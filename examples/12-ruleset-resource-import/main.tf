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

resource "ngenix_ruleset" "ngenixrulesetimport" {
  name    = "ngenixrulesetimport"
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
                id = 1003851
              }
            }
          ]
        }
      ],
      actions = [
        {
          action = "jsChallenge"
          params = ["3600"]
        }
      ]
    }
  ]
}
