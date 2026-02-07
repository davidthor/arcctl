terraform {
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

resource "google_service_account" "gke_nodes" {
  account_id   = "${var.name}-gke-nodes"
  display_name = "GKE Node Service Account for ${var.name}"
  project      = var.project
}

resource "google_project_iam_member" "gke_log_writer" {
  project = var.project
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.gke_nodes.email}"
}

resource "google_project_iam_member" "gke_metric_writer" {
  project = var.project
  role    = "roles/monitoring.metricWriter"
  member  = "serviceAccount:${google_service_account.gke_nodes.email}"
}

resource "google_project_iam_member" "gke_artifact_reader" {
  project = var.project
  role    = "roles/artifactregistry.reader"
  member  = "serviceAccount:${google_service_account.gke_nodes.email}"
}

resource "google_container_cluster" "main" {
  name     = var.name
  project  = var.project
  location = var.region

  network    = var.network
  subnetwork = var.subnet

  # Use separately managed node pool
  remove_default_node_pool = true
  initial_node_count       = 1

  # Enable Gateway API
  gateway_api_config {
    channel = "CHANNEL_STANDARD"
  }

  # Enable workload identity
  workload_identity_config {
    workload_pool = "${var.project}.svc.id.goog"
  }

  # VPC-native cluster
  ip_allocation_policy {
    cluster_ipv4_cidr_block  = "/16"
    services_ipv4_cidr_block = "/22"
  }

  release_channel {
    channel = "REGULAR"
  }

  resource_labels = {
    managed-by = "cldctl"
  }

  deletion_protection = false
}

resource "google_container_node_pool" "primary" {
  name     = "${var.name}-primary"
  project  = var.project
  location = var.region
  cluster  = google_container_cluster.main.name

  node_config {
    machine_type    = var.node_pool.machine_type
    service_account = google_service_account.gke_nodes.email
    oauth_scopes    = ["https://www.googleapis.com/auth/cloud-platform"]

    workload_metadata_config {
      mode = "GKE_METADATA"
    }

    labels = {
      managed-by = "cldctl"
    }
  }

  autoscaling {
    min_node_count = var.node_pool.min_nodes
    max_node_count = var.node_pool.max_nodes
  }

  management {
    auto_repair  = true
    auto_upgrade = true
  }
}
