terraform {
  required_providers {
    jiraserverfatih = {
      source = "mrkresnofatih/jiraserverfatih"
      version = "1.0.2-dev"
    }
  }
}

provider "jiraserverfatih" {
  # Configuration options
  authorizationMethod = "Bearer"
  token = "MzI2OTExNDY0MjI1OgH5gVAN5+57IbtOHB4TbB5EDSkZ"
  host     = "support.app-dev.cloudvanti.com"
}

resource "jiraserverfatih_projectrole" "anewrole" {
  name = "Abrandnewrole"
  description = "a brand new role 0938"
}
