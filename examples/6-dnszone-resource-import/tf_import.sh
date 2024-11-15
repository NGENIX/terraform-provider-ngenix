#!/bin/bash

TF_LOG=INFO terraform apply -var-file="vars.tfvars" -auto-approve

#terraform state show ngenix_dnszone.terraformimporttestex1

# Show plan
terraform show

# Remove TF state
terraform state rm ngenix_dnszone.terraformimporttestex1

# Import state by DNS zone ID (Read methos)
terraform import -var-file="vars.tfvars" ngenix_dnszone.terraformimporttestex1 <DNZ_ZONE_ID>
# TF_LOG=DEBUG terraform import -var-file="vars.tfvars" ngenix_dnszone.terraformimporttestex1 5905015

# View TF state 
terraform show