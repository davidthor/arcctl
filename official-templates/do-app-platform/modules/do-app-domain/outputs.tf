output "url" {
  description = "Application URL"
  value       = local.effective_url
}

output "host" {
  description = "Application hostname"
  value       = local.effective_host
}
