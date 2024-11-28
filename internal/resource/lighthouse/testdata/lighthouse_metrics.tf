provider "definednet" {
  token = "supersecret"
}

variable "metrics_listen" {
  type = string
}

variable "metrics_path" {
  type = string
}

variable "metrics_namespace" {
  type = string
}

variable "metrics_subsystem" {
  type = string
}

variable "metrics_enable_extra" {
  type = bool
}

resource "definednet_lighthouse" "metrics_test" {
  name             = "metrics-test"
  network_id       = "network-id"
  listen_port      = 4242
  static_addresses = ["127.0.0.1"]

  metrics {
    enabled              = true
    listen               = var.metrics_listen
    path                 = var.metrics_path
    namespace            = var.metrics_namespace
    subsystem            = var.metrics_subsystem
    enable_extra_metrics = var.metrics_enable_extra
  }
}
