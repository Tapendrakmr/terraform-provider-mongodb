
- This terraform providre allows to perform Create ,Read ,Update and Import MongoDB user
- To import a user use EmailID
- MongoDB doesn't provide an API to remove a user from organization

## Requirements

- [Go Lang](https://golang.org/doc/install)>=1.11 <br>
- [Terraform](https://www.terraform.io/downloads.html)>=v0.13.0 <br/>
- MongoDB Basic Account
- [MongoDB API Credentials](https://docs.atlas.mongodb.com/configure-api-access/)
- [MongoDB API documentation](https://docs.cloudmanager.mongodb.com/reference/api/users/)

## Application Account
This provider can be successfully tested on a basic mongoDB cloud manager account.

## Setup 
1.Create a mongoDB account(https://account.mongodb.com/account/)

## API Authentication
1. Login  to mongodb account
2. Click the Organization Settings icon next to the Organizations menu.
3. Toggle the Require IP Access List for Public API setting to On
4. Click Access Manager in the sidebar, or click Access Manager in the navigation bar, then click your organization.
5. Click Create API Key.
6. Enter the API Key Information
7. Click Next and save the Public and Private key.


## Building The Provider
1. Clone the repository, add all the dependencies and create a vendor directory that contains all dependencies. For this, run the following commands:
```bash
go mod init terraform-provider-mongodb
go mod tidy
go mod vendor
```
## Managing terraform plugins

- Clone or Download the repositor. <br>
- Add the Public key,Private key and Organization id generated in API Key
- Run the following command :

```bash
go mod init terraform-provider-mongodb
go mod tidy
```

- Run go mod vendor to create a vendor directory that contains all the provider's dependencies.

## Managing plugins for terraform
1. Run the following command to create a vendor subdirectory which will comprise of all provider dependencies. 
```
%APPDATA%/terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${target}
```
Command:
```
mkdir -p ~/.terraform.d/plugins/mongdb.com/edu/mongodb/0.2/[OS_ARCH]
```
For eg. mkdir -p %APPDATA%/.terraform.d/plugins/mongdb.com/edu/mongodb/0.2/windows_amd64

2. Run `go build -o terraform-provider-mongodb.exe`. This will save the binary (`.exe`) file in the main/root directory. <br>
3. Run this command to move this binary file to the appropriate location.

```
move terraform-provider-mongodb.exe %APPDATA%\terraform.d\plugins\mongdb.com\edu\mongodb\0.2\[OS_ARCH]
```
[OR]
1. Manually move the file from current directory to destination directory.<br>

## Working with terraform
### Application Credential Integration in terraform
1. Add terraform block and provider block as shown in example usage.  
2. Get credentials :Public Key, Private Key and Organization ID  
3.Assign the above credentials to the respective field in the provider block.` <br>

### Basic Terraform Commands

- `terraform init` -To initialize a working directory containing Terraform configuration files.
- `terraform plan` - To create an execution plan. Displays the changes to be done.
- `terraform apply` - To execute the actions proposed in a Terraform plan. Apply the changes.
- `terraform destroy` - Destroy previously-created
  infrastructure

### Invite User

1.  Add the user email and [Organization roles](https://docs.cloudmanager.mongodb.com/reference/api/user-update/) in the respective field as shown in  example usage
2. Run the basic terraform commands.
3.  You will see that a user has been successfully created and an invitation mail has been sent to the user.

```
resource "mongodb_user" "user1" {
  username="user@gmail.com"
  roles{
    org_id="[ORGANIZATION ID]"
    role_name="[ROLE NAME]"
  }
}
```
### Update the user

:heavy_exclamation_mark: [IMPORTANT] : If you own an organization or project, you can update the user roles for any user with membership in that organization or project. You cannot modify any other user profile information.

1. Update the data of the user resource block as show in example usage un the basic terraform commands to update user.User is not allowed to update email.
<br>

2. You can also allot a project with its projectId or groupId and [project roles](https://docs.cloudmanager.mongodb.com/reference/api/user-update/) to a user .

```
resource "mongodb_user" "user1" {
  username="user@gmail.com"
  roles{
    org_id="[ORGANIZATION ID]"
    role_name="[ROLE NAME]"
  }
  roles{
    group_id="[ORGANIZATION ID]"
    role_name="[ROLE NAME]"
  }
}
```
### Read the User Data

1. Add `data` and `output` blocks as shown in the [example usage](#example-usage) and run the basic terraform commands.

### Delete the user

:heavy_exclamation_mark: [IMPORTANT] : MongoDb doesn’t provide an API to delete users. To delete user from state file follow below instructions<br>

1. Delete the `resource` block of the user and run `terraform apply`.

### Import a User command

1.  Write manually a `resource` configuration block for the user as shown in [example usage](#example-usage). Imported user will be mapped to this block
2. Run the command `terraform import mongodb_user.resource_name [EMAIL_ID]`
3.Run `terraform plan`, if output shows `0 to add, 0 to change and 0 to destroy` user import is successful, otherwise recheck the user data in `resource` block with user data in mongoDB website.

### Testing the Provider

1. Navigate to the test file directory.
2. Run command `go test` for unit testing and for acceptance testing run command `TF_ACC=1 go test` . These commands will give combined test results for the execution or errors if any failure occurs.
3. If you want to see test result of each test function individually while running test in a single go, run command `go test -v`
4. To check test cover run `go test -cover`

## Example Usage <a id="example-usage"></a>

```terraform
terraform {
  required_providers {
    mongodb = {
      source = "Tapendrakmr/mongodb"
      version = "0.2.4"
    }
  }
}

provider "mongodb" {
  publickey  = "[PUBLIC KEY]"
  privatekey = "[PRIVATE KEY]"
  orgid="[ORGANIZATION ID]"
}

resource "mongodb_user" "user1" {
  username="user@gmail.com"
  roles{
    org_id="[ORGANIZATION ID]"
    role_name="[ROLE NAME]"
  }
  roles{
    group_id="[ORGANIZATION ID]"
    role_name="[ROLE NAME]"
  }
}

data "mongodb_user" "user" {
  id = "user@gmail.com"
}

output "datasource_user" {
  value = data.mongodb_user.user
}
```

## Argument Reference

- `email` (string) - The email id associated with the user account.
- `org_id` (string) - The id associated with the organization.
- `role_name` (string) - Type of [Roles](https://docs.cloudmanager.mongodb.com/reference/api/user-update/).
- `group_id` (string) - The id associated with the project.
- `role_name` (string) - Type of [Roles](https://docs.cloudmanager.mongodb.com/reference/api/user-update/).
