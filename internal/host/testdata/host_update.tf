provider "definednet" {
  token = "supersecret"
}

resource "definednet_host" "test" {
  name       = "updated-host.defined.test"
  network_id = "network-id"
  role_id    = "updated-role-id"
  tags       = ["tag:one", "tag:three"]
}
