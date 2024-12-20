---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "definednet_host Resource - definednet"
subcategory: ""
description: |-
  definednet_host enables managing Nebula overlay network hosts on Defined.net.
  The Defined.net API token must be configured with the following scope:
  hosts:createhosts:deletehosts:enrollhosts:listhosts:readhosts:update
---

# definednet_host (Resource)

`definednet_host` enables managing Nebula overlay network hosts on Defined.net.

The Defined.net API token must be configured with the following scope:

- `hosts:create`
- `hosts:delete`
- `hosts:enroll`
- `hosts:list`
- `hosts:read`
- `hosts:update`

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Host's name
- `network_id` (String) Enrolled Network ID

### Optional

- `metrics` (Block, Optional) Host's metrics exporter configuration (see [below for nested schema](#nestedblock--metrics))
- `role_id` (String) Host's role ID on Defined.net
- `tags` (List of String) Host's tags on Defined.net

### Read-Only

- `enrollment_code` (String, Sensitive) Host's enrollment code
- `id` (String) Host's ID
- `ip_address` (String) Host's IP address on Defined.net overlay network

<a id="nestedblock--metrics"></a>
### Nested Schema for `metrics`

Optional:

- `enable_extra_metrics` (Boolean) Enable extra metrics
- `enabled` (Boolean) Enable metrics exporter
- `listen` (String) Host-port for Prometheus metrics exporter listener
- `namespace` (String) Prometheus metrics' namespace
- `path` (String) Prometheus metrics exporter's HTTP path
- `subsystem` (String) Prometheus metrics' subsystem
