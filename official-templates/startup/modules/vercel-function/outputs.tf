output "function_id" {
  description = "Vercel deployment/function ID"
  value       = vercel_deployment.this.id
}

output "endpoint" {
  description = "Function endpoint URL"
  value       = "https://${vercel_deployment.this.url}"
}
