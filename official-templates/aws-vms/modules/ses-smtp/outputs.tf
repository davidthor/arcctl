output "host" {
  description = "SES SMTP host"
  value       = "email-smtp.${var.region}.amazonaws.com"
}

output "port" {
  description = "SES SMTP port"
  value       = 587
}
