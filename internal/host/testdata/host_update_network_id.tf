provider "definednet" {
  token = "supersecret"
}

resource "definednet_host" "test" {
  name       = "host.defined.test"
  network_id = "updated-network-id"
  role_id    = "role-id"
  tags       = ["tag:one", "tag:two"]
}
