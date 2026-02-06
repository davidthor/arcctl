variable "name" {
  description = "Secret name"
  type        = string
}

variable "data" {
  description = "Secret data as key-value pairs"
  type        = map(string)
  sensitive   = true
}
