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

resource "definednet_role" "test" {
  name       = var.name
  description = var.description
}
