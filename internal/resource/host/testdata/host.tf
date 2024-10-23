provider "definednet" {
  token = "supersecret"
}

variable "name" {
  type = string
}

variable "network_id" {
  type = string
}

variable "role_id" {
  type = string
}

variable "tags" {
  type = list(string)
}

resource "definednet_host" "test" {
  name       = var.name
  network_id = var.network_id
  role_id    = var.role_id
  tags       = var.tags
}
