provider "definednet" {
  token = "supersecret"
}

resource "definednet_host" "test" {
  name       = "host.defined.test"
  network_id = "network-id"
}
