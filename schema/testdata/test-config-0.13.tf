# terraform init && terraform providers schema -json
# with v0.13 generates relevant JSON file
terraform {
  required_providers {
    rand = {
      source  = "hashicorp/random"
      version = "3.0.0"
    }
    null = {
      source  = "hashicorp/null"
      version = "3.0.0"
    }
    grafana = {
      source  = "grafana/grafana"
      version = "1.6.0"
    }
  }
}

provider "null" {
  alias = "foobar"
}

resource "random_string" "name" {
  provider = rand
}

resource "null_resource" "name" {

}

data "terraform_remote_state" "vpc" {

}

resource "grafana_alert_notification" "slack" {

}
