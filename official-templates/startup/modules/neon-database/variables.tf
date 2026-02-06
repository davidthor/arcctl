variable "name" {
  description = "Database and role name"
  type        = string
}

variable "api_key" {
  description = "Neon API key"
  type        = string
  sensitive   = true
}

variable "project_id" {
  description = "Neon project ID"
  type        = string
}

variable "parent_branch" {
  description = "Parent branch name to fork from (null for production/main branch)"
  type        = string
  default     = null
}

variable "branch_name" {
  description = "Branch name for this environment"
  type        = string
}
