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
    port_from = number
    port_to = number
    protocol = string
    description = string
    allowed_role_id = string
    allowed_tags = list(string)
  }))

  default = []
}

resource "definednet_role" "test" {
  name       = var.name
  description = var.description

  dynamic "rule" {
    for_each = var.rules
    content {
      port_range = {
        from = rule.value.port_from
        to = rule.value.port_to
      }
      protocol = rule.value.protocol
      description = rule.value.description
      allowed_role_id = rule.value.allowed_role_id
      allowed_tags = rule.value.allowed_tags
    }
  }
}
