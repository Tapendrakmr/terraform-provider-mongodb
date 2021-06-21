terraform {
  required_providers {
    mongodb = {
      version = "0.2"
      source  = "mongdb.com/edu/mongodb"
    }
  }
}

provider "mongodb" {
  mongodb_public_key  = "REPLACE_MONGODBCLOUD_PUBLIC_KEY"
  mongodb_private_key = "REPLACE_MONGODBCLOUD_PRIVATE_KEY"
  mongodb_orgid="REPLACE_MONGODBCLOUD_ORGID"
}

resource "mongodb_user" "user" {
  username="user@gmail.com"
  roles{
    org_id="ORGID"
    role_name="ORG_READ_ONLY"
  }
}

data "mongodb_user" "user" {
  id = "user@gmail.com"
}

output "datasource_user" {
  value = data.mongodb_user.user
}