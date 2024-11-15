#!/bin/bash

TF_LOG=INFO terraform plan -var-file="vars.tfvars"

TF_LOG=INFO terraform plan -out tfplan -var-file="vars.tfvars"