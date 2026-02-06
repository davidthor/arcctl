output "lb_id" {
  description = "Load balancer ID"
  value       = digitalocean_loadbalancer.lb.id
}

output "ip" {
  description = "Load balancer public IP"
  value       = digitalocean_loadbalancer.lb.ip
}

output "urn" {
  description = "Load balancer URN"
  value       = digitalocean_loadbalancer.lb.urn
}
