output "image_uri" {
  description = "Full image URI in Artifact Registry"
  value       = local.image_tag
  depends_on  = [null_resource.docker_build]
}
