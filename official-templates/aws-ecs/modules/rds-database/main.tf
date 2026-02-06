terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

data "aws_subnets" "database" {
  filter {
    name   = "vpc-id"
    values = [var.vpc_id]
  }
}

locals {
  is_redis    = var.type == "redis"
  is_postgres = var.type == "postgres"
  is_mysql    = var.type == "mysql"
  db_name     = replace(var.name, "-", "_")

  engine         = local.is_postgres ? "postgres" : local.is_mysql ? "mysql" : "redis"
  default_port   = local.is_postgres ? 5432 : local.is_mysql ? 3306 : 6379
  engine_version = var.engine_version != null ? var.engine_version : local.is_postgres ? "16" : local.is_mysql ? "8.0" : "7.0"
}

resource "random_password" "this" {
  count   = local.is_redis ? 0 : 1
  length  = 32
  special = false
}

resource "aws_db_subnet_group" "this" {
  count      = local.is_redis ? 0 : 1
  name       = var.name
  subnet_ids = data.aws_subnets.database.ids

  tags = {
    Name      = var.name
    ManagedBy = "arcctl"
  }
}

resource "aws_db_instance" "this" {
  count = local.is_redis ? 0 : 1

  identifier     = var.name
  engine         = local.engine
  engine_version = local.engine_version

  instance_class    = var.instance_class
  allocated_storage = var.allocated_storage

  db_name  = local.is_postgres || local.is_mysql ? local.db_name : null
  username = local.is_postgres ? "arcctl_admin" : "admin"
  password = random_password.this[0].result

  db_subnet_group_name   = aws_db_subnet_group.this[0].name
  vpc_security_group_ids = [var.security_group_id]

  skip_final_snapshot = true
  publicly_accessible = false

  tags = {
    Name        = var.name
    Environment = var.name
    ManagedBy   = "arcctl"
  }
}

# ElastiCache for Redis
resource "aws_elasticache_subnet_group" "this" {
  count      = local.is_redis ? 1 : 0
  name       = var.name
  subnet_ids = data.aws_subnets.database.ids

  tags = {
    Name      = var.name
    ManagedBy = "arcctl"
  }
}

resource "random_password" "redis" {
  count   = local.is_redis ? 1 : 0
  length  = 32
  special = false
}

resource "aws_elasticache_replication_group" "this" {
  count = local.is_redis ? 1 : 0

  replication_group_id = var.name
  description          = "Redis cluster for ${var.name}"

  engine               = "redis"
  engine_version       = local.engine_version
  node_type            = replace(var.instance_class, "db.", "cache.")
  num_cache_clusters   = 1
  port                 = local.default_port
  subnet_group_name    = aws_elasticache_subnet_group.this[0].name
  security_group_ids   = [var.security_group_id]

  transit_encryption_enabled = true
  auth_token                 = random_password.redis[0].result

  tags = {
    Name        = var.name
    Environment = var.name
    ManagedBy   = "arcctl"
  }
}

locals {
  # RDS outputs
  rds_endpoint = local.is_redis ? "" : aws_db_instance.this[0].address
  rds_port     = local.is_redis ? 0 : aws_db_instance.this[0].port
  rds_username = local.is_redis ? "" : aws_db_instance.this[0].username
  rds_password = local.is_redis ? "" : random_password.this[0].result
  rds_database = local.is_redis ? "" : aws_db_instance.this[0].db_name

  # Redis outputs
  redis_endpoint = local.is_redis ? aws_elasticache_replication_group.this[0].primary_endpoint_address : ""
  redis_port     = local.is_redis ? aws_elasticache_replication_group.this[0].port : 0
  redis_password = local.is_redis ? random_password.redis[0].result : ""

  # Unified outputs
  endpoint      = local.is_redis ? local.redis_endpoint : local.rds_endpoint
  port          = local.is_redis ? local.redis_port : local.rds_port
  username      = local.is_redis ? "" : local.rds_username
  password      = local.is_redis ? local.redis_password : local.rds_password
  database_name = local.is_redis ? "0" : local.rds_database

  connection_url = local.is_redis ? "redis://:${local.redis_password}@${local.redis_endpoint}:${local.redis_port}/0" : local.is_postgres ? "postgresql://${local.rds_username}:${local.rds_password}@${local.rds_endpoint}:${local.rds_port}/${local.rds_database}" : "mysql://${local.rds_username}:${local.rds_password}@${local.rds_endpoint}:${local.rds_port}/${local.rds_database}"
}
