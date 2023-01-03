terraform {
  required_providers {
    workos = {
      source = "hashicorp.com/aleshchynskyi/workos"
    }
  }
}

provider "workos" {}

data "workos_organizations" "example" {}
