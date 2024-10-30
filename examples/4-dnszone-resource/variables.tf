variable "ngenix_api_url" {
  description = "Ngenix API URL"
  type        = string
}

variable "ngenix_email" {
  description = "Ngenix User email address"
  type        = string
  sensitive   = true
}

variable "ngenix_token" {
  description = "Ngenix User access token"
  type        = string
  sensitive   = true
}