terraform {
  required_providers {
    hashicups = {
      source  = "hashicorp.com/edu/hashicups-pf"
    }
  }
  required_version = ">= 1.1.0"
}

provider "hashicups" {
  authorizationMethod = "Bearer"
  token = "dsdadadsadsdasdadsadasadasdasdas"
  host     = "sample.app-dev.fatihcompany.com"
}

resource "hashicups_projectrole" "anewrole" {
  name = "Abrandnewrole"
  description = "a brand new role 0938"
}
