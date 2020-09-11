# terraform init && terraform providers schema -json
# with v0.12 generates relevant JSON file
terraform {
  required_providers {
    random = "3.0.0"
    null = "3.0.0"
  }
}

provider "null" {
  alias = "foobar"
}

resource "null_resource" "name" {

}

resource "random_string" "name" {

}

data "terraform_remote_state" "vpc" {

}
