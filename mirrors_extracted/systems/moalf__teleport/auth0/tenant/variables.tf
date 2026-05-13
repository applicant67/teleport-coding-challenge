variable "auth0_domain" {
  type        = string
  description = "Auth0 tenant domain"
  sensitive   = true
}

variable "auth0_client_id" {
  type      = string
  sensitive = true
}

variable "auth0_client_secret" {
  type      = string
  sensitive = true
}
