variable "definednet_token" {
  description = "Defined.net HTTP API token"
  sensitive   = true
}

provider "definednet" {
  token = var.definednet_token
}

resource "definednet_host" "example" {
  name       = "example.defined.test"
  network_id = "network-7P81MCS2TVAY9XJWQTNJ3PWYPD"
  role_id    = "role-WSG78880Z655TQJVQFL5CZ405B"
  tags       = ["service:app"]
}

resource "definednet_host" "metrics_minimal" {
  name       = "example.defined.test"
  network_id = "network-7P81MCS2TVAY9XJWQTNJ3PWYPD"
  role_id    = "role-WSG78880Z655TQJVQFL5CZ405B"
  tags       = ["service:app"]

  metrics {
    enabled = true
  }
}

resource "definednet_host" "metrics" {
  name       = "example.defined.test"
  network_id = "network-7P81MCS2TVAY9XJWQTNJ3PWYPD"
  role_id    = "role-WSG78880Z655TQJVQFL5CZ405B"
  tags       = ["service:app"]

  metrics {
    enabled              = true
    listen               = "127.0.0.1:9100"
    path                 = "/-/metrics"
    namespace            = "infra"
    subsystem            = "nebula"
    enable_extra_metrics = true
  }
}
