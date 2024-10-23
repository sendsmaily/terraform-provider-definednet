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

variable "listen_port" {
  type = number
}

variable "static_addresses" {
  type = list(string)
}

variable "tags" {
  type = list(string)
}

resource "definednet_lighthouse" "test" {
  name             = var.name
  network_id       = var.network_id
  role_id          = var.role_id
  listen_port      = var.listen_port
  static_addresses = var.static_addresses
  tags             = var.tags
}
