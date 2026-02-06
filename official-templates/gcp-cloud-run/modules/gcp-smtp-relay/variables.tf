variable "host" {
  description = "SMTP relay host"
  type        = string
}

variable "port" {
  description = "SMTP relay port"
  type        = number
}

variable "username" {
  description = "SMTP relay username"
  type        = string
}

variable "password" {
  description = "SMTP relay password"
  type        = string
  sensitive   = true
}
