output "cluster_arn" {
  description = "EKS cluster ARN"
  value       = aws_eks_cluster.this.arn
}

output "cluster_name" {
  description = "EKS cluster name"
  value       = aws_eks_cluster.this.name
}

output "endpoint" {
  description = "EKS cluster API endpoint"
  value       = aws_eks_cluster.this.endpoint
}

output "kubeconfig" {
  description = "Kubeconfig for accessing the cluster"
  value       = local.kubeconfig
  sensitive   = true
}

output "certificate_authority" {
  description = "Cluster certificate authority data"
  value       = aws_eks_cluster.this.certificate_authority[0].data
}
