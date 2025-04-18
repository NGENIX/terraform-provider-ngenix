terraform {
  required_providers {
    ngenix = {
      source = "ngenix.net/api/ngenix"
    }
  }
}

provider "ngenix" {
  host     = var.ngenix_api_url
  username = var.ngenix_email
  password = var.ngenix_token
}

data "ngenix_traffic_patterns" "example" {}

output "example_dnszone" {
  value = data.ngenix_traffic_patterns.example
}
