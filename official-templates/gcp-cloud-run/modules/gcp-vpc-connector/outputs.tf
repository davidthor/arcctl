output "id" {
  description = "The ID of the VPC connector"
  value       = google_vpc_access_connector.main.id
}

output "name" {
  description = "The name of the VPC connector"
  value       = google_vpc_access_connector.main.name
}
