output "otlp_endpoint" {
  description = "OTLP HTTP endpoint for sending telemetry"
  value       = digitalocean_app.monitoring.live_url
}

output "loki_endpoint" {
  description = "Log query endpoint (App Platform built-in logging)"
  value       = "https://api.digitalocean.com/v2/apps/${digitalocean_app.monitoring.id}/logs"
}

output "dashboard_url" {
  description = "DigitalOcean monitoring dashboard URL"
  value       = "https://cloud.digitalocean.com/apps/${digitalocean_app.monitoring.id}/overview"
}
