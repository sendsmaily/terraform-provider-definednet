variable "definednet_token" {
  description = "Defined.net HTTP API token"
  sensitive   = true
}

provider "definednet" {
  token = var.definednet_token
}

resource "definednet_lighthouse" "example" {
  name             = "example.defined.test"
  network_id       = "network-7P81MCS2TVAY9XJWQTNJ3PWYPD"
  role_id          = "role-WSG78880Z655TQJVQFL5CZ405B"
  listen_port      = 4242
  static_addresses = ["84.123.10.1"]
  tags             = ["service:app"]
}
