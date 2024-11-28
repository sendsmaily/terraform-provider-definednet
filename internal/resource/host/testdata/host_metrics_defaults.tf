provider "definednet" {
  token = "supersecret"
}

resource "definednet_host" "metrics_default_test" {
  name       = "metrics-test"
  network_id = "network-id"
  role_id    = "role-id"

  metrics {
    enabled = true
  }
}
