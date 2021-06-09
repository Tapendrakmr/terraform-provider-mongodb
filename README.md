# MongoDB Cloud Manger Terraform Provider

- This terraform providre allows to perform Create ,Read ,Update and Import MongoDB user
- To import a user use EmailID
- MongoDB doesn't provide an API to remove a user from organization
- [MongoDB API documentation](https://docs.cloudmanager.mongodb.com/reference/api/users/)

## Requirements

- [Go Lang](https://golang.org/doc/install)>=1.11 <br>
- [Terraform](https://www.terraform.io/downloads.html)>=v0.13.0 <br/>
- MongoDB Basic Account
- [MongoDB API Credentials](https://docs.atlas.mongodb.com/configure-api-access/)

## Application Setup

- Create an organization and save Public,Private key of organization

## Building The Provider

- Clone or Download the repositor. <br>
- Add the Public key,Private key and Organization id generated in API Key
- Run the following command :

```bash
go mod init terraform-provider-mongodb
go mod tidy
```

- Run go mod vendor to create a vendor directory that contains all the provider's dependencies.

## Managing plugins for terraform

- Run `go build -o terraform-provider-mongodb.exe`. This will save the binary (`.exe`) file in the main/root directory. <br>
- Run this command to move this binary file to the appropriate location.

```
~/.terraform.d/plugins/${host_name}/${namespace}/${type}/${version}/${target}
```

```
mkdir -p ~/.terraform.d/plugins/mongdb.com/edu/mongodb/0.2/[OS_ARCH]

move terraform-provider-mongodb.exe %APPDATA%\terraform.d\plugins\mongdb.com\edu\mongodb\0.2\[OS_ARCH]
```

Otherwise you can manually move the file from current directory to destination directory.<br>

[OR]

1. Download required binaries <br>
2. move binary `~/.terraform.d/plugins/[architecture name]/`

## Working with terraform

1. Add the Public Key, Private Key and Organization ID to respective fields in `main.tf` <br>

### Basic Terraform Commands

- `terraform init` - Prepare your working directory for other commands
- `terraform plan` - Show changes required by the current configuration
- `terraform apply` - Create or update infrastructure
- `terraform destroy` - Destroy previously-created
  infrastructure

### Invite User

1.  Add the user email and [Organization roles](https://docs.cloudmanager.mongodb.com/reference/api/user-update/) in the respective field in `main.tf`
2.  Initialize the terraform provider `terraform init`
3.  Check the changes applicable using `terraform plan` and apply using `terraform apply`
4.  You will see that a user has been successfully created and an invitation mail has been sent to the user.

```
resource "mongodb_user" "user1" {
  username="user@gmail.com"
  roles{
    org_id="[ORGANIZATION ID]"
    role_name="[ROLE NAME]"
  }
}
```

### Read the User Data

1. Add data and output blocks in the `main.tf` file after that add email field as id and run `terraform plan` to read user data

### Update the user

:heavy_exclamation_mark: [IMPORTANT] : If you own an organization or project, you can update the user roles for any user with membership in that organization or project. You cannot modify any other user profile information.

1. Update the data of the user in the `main.tf` file and apply using `terraform apply`<br>
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

### Delete the user

:heavy_exclamation_mark: [IMPORTANT] : MongoDb doesn’t provide an API to delete users. To delete user from state file follow below instructions<br>

1. Delete the resource block of the particular user from the `main.tf` file and run `terraform apply`.

### Import a User command

1. Write manually a resource configuration block for the resource in `main.tf`, to which the imported object will be mapped.
2. RUN terraform import mongodb_user.resource_name <user_id>

### Testing the Provider

1. Navigate to the test file directory.
2. Run command `go test` for unit testing and for acceptance testing run command `TF_ACC=1 go test` . These commands will give combined test results for the execution or errors if any failure occurs.
3. If you want to see test result of each test function individually while running test in a single go, run command `go test -v`
4. To check test cover run `go test -cover`

## Example Usage

```terraform
terraform {
  required_providers {
    mongodb = {
      version = "0.2"
      source  = "mongodb.com/edu/mongodb"
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

output "user_result" {
  value = mongodb_user.user1
}
```

## Argument Reference

- `email` (string) - The email id associated with the user account.
- `org_id` (string) - The id associated with the organization.
- `role_name` (string) - Type of [Roles](https://docs.cloudmanager.mongodb.com/reference/api/user-update/).
- `group_id` (string) - The id associated with the project.
- `role` (string) - Type of [Roles](https://docs.cloudmanager.mongodb.com/reference/api/user-update/).
