output "vpc_id" {
  description = "VPC ID"
  value       = digitalocean_vpc.vpc.id
}

output "urn" {
  description = "VPC URN"
  value       = digitalocean_vpc.vpc.urn
}
