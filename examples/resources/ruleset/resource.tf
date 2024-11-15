# Manage ruleset example
resource "ngenix_ruleset" "ngenixrulesetexampleone" {
  name    = "ngenixrulesetexampleone"
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
          action = "deny"
          params = ["403", "wontbe"]
        }
      ]
    }
  ]
}