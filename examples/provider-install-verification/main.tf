terraform {
  required_providers {
    workos = {
      source = "vellanci/workos"
    }
  }
}

provider "workos" {
  host = "https://api.workos.com"
}

data "workos_organizations" "example" {}
