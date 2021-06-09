terraform {
  required_providers {
    mongodb = {
      version = "0.2"
      source  = "mongdb.com/edu/mongodb"
    }
  }
}

provider "mongodb" {
  publickey  = "PUBLIC KEY"
  privatekey = "PRIVATE KEY"
  orgid="ORGID"
}
# For new user
resource "mongodb_user" "user" {
  username="user@gmail.com"
  roles{
    org_id="ORGID"
    role_name="ORG_READ_ONLY"
  }
}

# import user
resource "mongodb_user" "newuser"{
  
}
# Read existing user
data "mongodb_user" "user" {
  id = "user@gmail.com"
}

output "user" {
  value = data.mongodb_user.user
}