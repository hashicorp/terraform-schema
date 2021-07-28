terraform {
  required_providers {
    hashicup = {
      source  = "hashicorp/hashicup"
      version = "0.0.0"
    }
  }
}

module "test" {
  source = "source"
}
