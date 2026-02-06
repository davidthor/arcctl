output "host" {
  description = "SMTP host"
  value       = "smtp.resend.com"
}

output "port" {
  description = "SMTP port (TLS)"
  value       = 465
}

output "username" {
  description = "SMTP username"
  value       = "resend"
}

output "password" {
  description = "SMTP password (Resend API key)"
  value       = var.api_key
  sensitive   = true
}
