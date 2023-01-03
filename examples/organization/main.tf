terraform {
  required_providers {
    workos = {
      source = "hashicorp.com/aleshchynskyi/workos"
    }
  }
  required_version = ">= 1.1.0"
}

provider "workos" {
  host = "https://api.workos.com"
}

resource "workos_organization" "example" {
  name = "Provided Org"
  domains = [{
    domain = "provided-org.co"
  }]
}

output "example_organizations" {
  value = workos_organization.example
}
