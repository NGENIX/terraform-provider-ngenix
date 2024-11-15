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

// Addr
resource "ngenix_traffic_pattern" "ngenixtpex1" {
  name         = "test-import"
  type         = "commonlist"
  content_type = "addr"
  patterns = [
    {
      addr    = "98.164.15.2/32"
      expires = 1734182958
    }
  ]
}
