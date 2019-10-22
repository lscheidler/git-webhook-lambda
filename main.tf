terraform {
  required_version = ">= 0.8, < 0.13"
}

# Specify the provider and access details
provider "aws" {
  region = "eu-central-1"
}

