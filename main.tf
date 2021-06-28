# terraform {
#   required_providers {
#     mongodb = {
#       version = "0.2"
#       source  = "mongdb.com/edu/mongodb"
#     }
#   }
# }

# provider "mongodb" {
#   mongodb_public_key  = var.mongodb_public_key
#   mongodb_private_key = var.mongodb_private_key
#   mongodb_orgid=var.mongodb_private_key
# }

# resource "mongodb_user" "user" {
#   username="user@gmail.com"
#   roles{
#     org_id="ORGID"
#     role_name="ORG_READ_ONLY"
#   }
# }

# data "mongodb_user" "user" {
#   id = "user@gmail.com"
# }

# output "datasource_user" {
#   value = data.mongodb_user.user
# }
terraform {
  required_providers {
    mongodb = {
      version = "0.2"
      source  = "hashicorp.com/edu/mongodb"
    }
  }
}

