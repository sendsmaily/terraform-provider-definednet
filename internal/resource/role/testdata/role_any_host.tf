provider "definednet" {
  token = "supersecret"
}

variable "name" {
  type = string
}

variable "description" {
  type = string
  default = ""
}

variable "rules" {
  type = list(object({
    port = number
    protocol = string
    description = string
  }))

  default = []
}

resource "definednet_role" "test" {
  name       = var.name
  description = var.description

  dynamic "rule" {
    for_each = var.rules
    content {
      port = rule.value.port
      protocol = rule.value.protocol
      description = rule.value.description
    }
  }
}
