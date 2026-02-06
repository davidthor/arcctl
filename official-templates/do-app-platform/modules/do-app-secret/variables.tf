variable "name" {
  description = "Secret name"
  type        = string
}

variable "token" {
  description = "DigitalOcean API token"
  type        = string
  sensitive   = true
}

variable "key_type" {
  description = "Encryption key type: rsa, ecdsa, or symmetric (null for plain secrets)"
  type        = string
  default     = null
}

variable "key_size" {
  description = "Key size in bits"
  type        = number
  default     = null
}

variable "data" {
  description = "Secret data as key-value pairs"
  type        = map(string)
  default     = null
  sensitive   = true
}
