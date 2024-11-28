provider "definednet" {
  token = "supersecret"
}

variable "name" {
  type = string
}

variable "network_id" {
  type = string
}

variable "listen_port" {
  type = number
}

variable "static_addresses" {
  type = list(string)
}

resource "definednet_lighthouse" "minimal_test" {
  name             = var.name
  network_id       = var.network_id
  listen_port      = var.listen_port
  static_addresses = var.static_addresses
}
