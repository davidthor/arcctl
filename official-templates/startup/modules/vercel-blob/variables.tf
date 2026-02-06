variable "name" {
  description = "Blob store name"
  type        = string
}

variable "token" {
  description = "Vercel API token"
  type        = string
  sensitive   = true
}

variable "team_id" {
  description = "Vercel team ID (optional for personal accounts)"
  type        = string
  default     = ""
}

variable "public" {
  description = "Whether the blob store allows public access"
  type        = bool
  default     = false
}
