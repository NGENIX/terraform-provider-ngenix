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

// Traffic pattern httpMethod
resource "ngenix_traffic_pattern" "testhttpruleset" {
  name         = "test-http-for-ruleset"
  type         = "commonlist"
  content_type = "httpMethod"
  patterns = [
    {
      http_method = "PATCH"
    },
    {
      http_method = "POST"
    },
    {
      http_method = "DELETE"
    }
  ]
}

// Ruleset example with TP id link 
resource "ngenix_ruleset" "ngenixrulesetex1" {
  name    = "ngenixrulesetexample"
  enabled = true
  rules = [
    {
      name    = "rule_one"
      enabled = true
      conditions = [
        {
          function = "in"
          negation = true
          params = [
            {
              variable = "http.request.method"
            },
            {
              trafficpattern_ref = {
                id = ngenix_traffic_pattern.testhttpruleset.id
              }
            }
          ]
        }
      ],
      actions = [
        {
          action = "deny"
          params = ["403", "access_denied"]
        }
      ]
    }
  ]
  depends_on = [ngenix_traffic_pattern.testhttpruleset]
}

