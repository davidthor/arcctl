output "ready" {
  description = "Whether the ALB controller is ready"
  value       = true
  depends_on  = [helm_release.alb_controller]
}

output "alb_dns_name" {
  description = "ALB DNS name (populated after ingress creation)"
  value       = ""
}
