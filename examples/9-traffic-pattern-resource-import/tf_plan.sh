#!/bin/bash

#TF_LOG=INFO terraform apply -var-file="vars.tfvars"

# Create / Update
terraform plan -var-file="vars.tfvars"
terraform apply -var-file="vars.tfvars"

# Import check
terraform show
terraform state rm ngenix_traffic_pattern.ngenixtpex1
terraform show
terraform import -var-file="vars.tfvars" ngenix_traffic_pattern.ngenixtpex1 1003850