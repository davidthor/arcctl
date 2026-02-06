output "app_id" {
  description = "App Platform application ID"
  value       = digitalocean_app.service.id
}

output "url" {
  description = "App Platform live URL"
  value       = digitalocean_app.service.live_url
}

output "default_ingress" {
  description = "Default ingress URL"
  value       = digitalocean_app.service.default_ingress
}
