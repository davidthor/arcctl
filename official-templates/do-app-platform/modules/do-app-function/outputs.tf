output "function_id" {
  description = "App Platform application ID"
  value       = digitalocean_app.function.id
}

output "url" {
  description = "Function endpoint URL"
  value       = digitalocean_app.function.live_url
}
