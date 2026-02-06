variable "name" {
  description = "SMTP configuration name"
  type        = string
}

variable "api_key" {
  description = "Resend API key (used as SMTP password)"
  type        = string
  sensitive   = true
}
