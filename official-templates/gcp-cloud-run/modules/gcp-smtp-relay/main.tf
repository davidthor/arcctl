terraform {
  # Pass-through module for external SMTP relay configuration.
  # GCP does not provide a native SMTP service; this module simply
  # forwards the configured SMTP credentials as outputs.
}
