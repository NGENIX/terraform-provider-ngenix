#!/bin/bash

#TF_LOG=INFO terraform apply -var-file="vars.tfvars"

#terraform state show ngenix_dnszone.ngenixterraformex5

# Create / Update
terraform plan -var-file="vars.tfvars"
terraform apply -var-file="vars.tfvars"

# Import check
terraform show
terraform state rm ngenix_dnszone.ngenixterraformex10
terraform show
terraform import -var-file="vars.tfvars" ngenix_dnszone.ngenixterraformex10 5905015