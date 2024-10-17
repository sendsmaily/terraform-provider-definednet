variable "definednet_token" {
  description = "Defined.net HTTP API token"
  sensitive   = true
}

provider "definednet" {
  token = var.definednet_token
}
