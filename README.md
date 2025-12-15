# Terraform Provider for AWS Redshift (Fork with Serverless Support)

> [!NOTE]
> This is a forked version of the original [brainly/terraform-provider-redshift](https://github.com/brainly/terraform-provider-redshift) repository.
> 
> **Original Repository Status:** Deprecated - The original repository is no longer maintained.
>
> **This Fork:** Adds support for AWS Redshift Serverless features and continues maintenance.

This provider allows you to manage [AWS Redshift](https://aws.amazon.com/redshift/) objects with Terraform, including users, groups, schemas, databases, and permissions. 

**✨ Key Addition:** Full support for **AWS Redshift Serverless** alongside traditional Redshift clusters.

📖 Full documentation available on the [Terraform Registry](https://registry.terraform.io/providers/serenityzn/redshift/latest/docs).

## Features

- 🚀 **Redshift Serverless Support**: Use `type = "serverless"` to work with Redshift Serverless workgroups
- ⚙️ **Redshift Cluster Support**: Works with traditional provisioned Redshift clusters (default behavior)
- 🔐 **Multiple Authentication Methods**: Fixed passwords, temporary credentials, and cross-account roles
- 📦 **Complete Resource Management**: Users, groups, schemas, databases, grants, default privileges, and datashares
- 🔄 **Backward Compatible**: Drop-in replacement for the original provider

## Quick Start

### Step 1: Add Provider to Configuration

```hcl
terraform {
  required_providers {
    redshift = {
      source  = "serenityzn/redshift"
      version = "~> 1.0"
    }
  }
}
```

### Step 2: Configure Provider

#### For Redshift Serverless
```hcl
provider "redshift" {
  host     = "your-workgroup.123456789.us-east-1.redshift-serverless.amazonaws.com"
  username = "admin"
  password = var.redshift_password
  database = "dev"
  type     = "serverless"  # Required for Redshift Serverless
}
```

#### For Regular Redshift Clusters
```hcl
provider "redshift" {
  host     = "my-cluster.123456789.us-east-1.redshift.amazonaws.com"
  username = "admin"
  password = var.redshift_password
  database = "dev"
  # type parameter omitted for regular clusters
}
```

### Step 3: Create Resources

```hcl
resource "redshift_schema" "example" {
  name  = "my_schema"
  owner = "admin"
}

resource "redshift_user" "example" {
  name     = "my_user"
  password = "MySecurePassword123!"
}
```

## Requirements

  - [Terraform](https://www.terraform.io/downloads.html) >= 1.0
  - [Go](https://golang.org/doc/install) 1.17 (to build the provider plugin)

## Building The Provider

```sh
$ git clone git@github.com:YOUR_USERNAME/terraform-provider-redshift
```

Enter the provider directory and build the provider

```sh
$ cd terraform-provider-redshift
$ make build
```
## Development

If you're new to provider development, a good place to start is the [Extending
Terraform](https://www.terraform.io/docs/extend/index.html) docs.

### Running Tests

Acceptance tests require a running real AWS Redshift cluster. 

```sh
REDSHIFT_HOST=<cluster ip or DNS>
REDSHIFT_USER=root
REDSHIFT_DATABASE=redshift
REDSHIFT_PASSWORD=<password>
make testacc
```

If your cluster is only accessible from within the VPC, you can connect via a socks proxy:
```sh
ALL_PROXY=socks5[h]://[<socks-user>:<socks-password>@]<socks-host>[:<socks-port>]
NO_PROXY=127.0.0.1,192.168.0.0/24,*.example.com,localhost
```

## Documentation

Documentation is generated with
[tfplugindocs](https://github.com/hashicorp/terraform-plugin-docs). Generated
files are in `docs/` and should not be updated manually. They are derived from:

* Schema `Description` fields in the provider Go code.
* [examples/](./examples)
* [templates/](./templates)

Use `go generate` to update generated docs.

## Releasing

Builds and releases are automated with GitHub Actions and [GoReleaser](https://github.com/goreleaser/goreleaser/). 
The changelog is managed with [github-changelog-generator](https://github.com/github-changelog-generator/github-changelog-generator).

Currently there are a few manual steps to this:

1. Update the changelog:

   ```sh
   RELEASE_VERSION=v... \
   CHANGELOG_GITHUB_TOKEN=... \
   make changelog
   ```

   This will commit the changelog locally.

2. Review generated changelog and push:

   View the committed changelog with `git show`. If all is well `git push origin
   master`.

3. Kick off the release:

   ```sh
   RELEASE_VERSION=v... \
   make release
   ```

   Once the command exits, you can monitor the rest of the process on the
   [Actions UI](https://github.com/serenityzn/terraform-provider-redshift/actions?query=workflow%3Arelease).

4. Publish release:

   The Action creates the release, but leaves it in "draft" state. Open it up in
   a [browser](https://github.com/serenityzn/terraform-provider-redshift/releases)
   and if all looks well, click the publish button.

## Local Development and Testing

To test your local provider build, you can use a local plugin directory:

### Step 1: Build the Provider
```bash
make build
```

### Step 2: Create a Local Plugin Directory
```bash
# For macOS ARM64 (adjust OS/architecture as needed)
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/serenityzn/redshift/99.0.0/darwin_arm64/
```

### Step 3: Copy Your Provider
```bash
cp $GOPATH/bin/terraform-provider-redshift ~/.terraform.d/plugins/registry.terraform.io/serenityzn/redshift/99.0.0/darwin_arm64/terraform-provider-redshift_v99.0.0
```

### Step 4: Use in Terraform
```hcl
terraform {
  required_providers {
    redshift = {
      source  = "serenityzn/redshift"
      version = "99.0.0"
    }
  }
}
```

This approach:
- ✅ **Uses your local provider build**
- ✅ **Works with Terraform and Terragrunt**
- ✅ **Is the standard way to test providers locally**