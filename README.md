# Terraform Provider JiraServerFatih

A Terraform Provider for JIRA Server (Self-Hosted Option). At the moment, resources available for provisioning in Terraform JiraServerFatih are as follows:

- Project Role
- Group
- Permission Scheme
- Permission Scheme Grant
- Issue Type

```terraform
terraform {
  required_providers {
    jiraserverfatih = {
      source = "mrkresnofatih/jiraserverfatih"
      version = "1.0.20230316224912-dev"
    }
  }
}

provider "jiraserverfatih" {
    authorization_method = "Bearer"     # Valid Values: Bearer | Basic
    token = "XXXXXXXXXXXXXXXXXXXXXXXX"  # Required
    domain     = "myjiraselfhosted.com" # Required
}

resource "jiraserverfatih_projectrole" "partneradminrole" {
  name = "Partner Admin"              # Required
  description = "a fatih new role23"  # Required
}

resource "jiraserverfatih_group" "myhostadmingroup" {
  # Groups cannot be updated, can only be created or destroyed
  name = "myhostadmingrupp" # Required
}

resource "jiraserverfatih_permissionscheme" "mypmsch" {
  name = "mypmsch"                                    # Required
  description = "fatih's test permission schemeee"    # Required
}

resource "jiraserverfatih_grant" "partneradminroleaddcomment" {
  permission_scheme_id = jiraserverfatih_permissionscheme.mypmsch.permission_scheme_id  # Required
  permission_name = "ADD_COMMENTS"                                                      # Required, Valid Values Coming Soon
  security_type = "projectrole"                                                         # Required, Valid Values: projectroles
  security_param = jiraserverfatih_projectrole.partneradminrole.project_role_id         # Required
}

resource "jiraserverfatih_issuetype" "mysuperissuetype" {
  name = "mysuperissuetyp"                    # Required
  description = "my super issue type desc2"   # Required
  avatar_id = 10304                           # Required
}
```