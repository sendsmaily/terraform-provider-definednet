provider "definednet" {
  token = "supersecret"
}

variable "name" {
  type = string
}

variable "network_id" {
  type = string
}

resource "definednet_host" "minimal_test" {
  name       = var.name
  network_id = var.network_id
}
