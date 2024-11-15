#!/bin/bash

#TF_LOG=INFO terraform apply -var-file="vars.tfvars"

# Create / Update
terraform plan -var-file="vars.tfvars"
terraform apply -var-file="vars.tfvars"

# Import check
terraform show
terraform state rm ngenix_ruleset.ngenixrulesetimport
terraform show
terraform import -var-file="vars.tfvars" ngenix_ruleset.ngenixrulesetimport 1352871