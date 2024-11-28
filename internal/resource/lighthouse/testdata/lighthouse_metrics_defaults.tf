provider "definednet" {
  token = "supersecret"
}

resource "definednet_lighthouse" "metrics_default_test" {
  name             = "metrics-test"
  network_id       = "network-id"
  listen_port      = 4242
  static_addresses = ["127.0.0.1"]

  metrics {
    enabled = true
  }
}
