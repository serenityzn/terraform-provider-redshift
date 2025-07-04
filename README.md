# Terraform Provider for AWS Redshift (Fork with Serverless Support)

> [!NOTE]
> This is a forked version of the original [brainly/terraform-provider-redshift](https://github.com/brainly/terraform-provider-redshift) repository.
> 
> **Original Repository Status:** Deprecated - The original repository is no longer maintained.
>
> **This Fork:** Adds support for AWS Redshift Serverless features and continues maintenance.

This provider allows to manage with Terraform [AWS Redshift](https://aws.amazon.com/redshift/) objects like users, groups, schemas, etc., including support for Redshift Serverless.

It's published on the [Terraform registry](https://registry.terraform.io/providers/serenityzn/redshift/latest/docs).

## Features

- **Redshift Serverless Support**: Use the `type = "serverless"` parameter to work with Redshift Serverless deployments
- **Backward Compatibility**: Works with regular Redshift clusters when `type` is omitted or empty
- **All Original Features**: Users, groups, schemas, grants, and more

## Provider Configuration

### For Redshift Serverless
```hcl
provider "redshift" {
  host     = "your-serverless-endpoint.redshift-serverless.amazonaws.com"
  username = "your_username"
  password = "your_password"
  database = "your_database"
  port     = 5439
  sslmode  = "require"
  type     = "serverless"  # Required for Redshift Serverless
}
```

### For Regular Redshift Clusters
```hcl
provider "redshift" {
  host     = "your-cluster.redshift.amazonaws.com"
  username = "your_username"
  password = "your_password"
  database = "your_database"
  port     = 5439
  sslmode  = "require"
  # type parameter omitted for regular clusters
}
```

## Requirements

  - [Terraform](https://www.terraform.io/downloads.html) >= 1.0
  - [Go](https://golang.org/doc/install) 1.17 (to build the provider plugin)

## Building The Provider

```sh
$ git clone git@github.com:serenityzn/terraform-provider-redshift
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

## Final Solution: Use a Local Provider Plugin Directory

The most reliable way to test your local provider is to use a local plugin directory. Here's what you need to do:

### Step 1: Create a Local Plugin Directory
```bash
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/serenityzn/redshift/1.0.0/darwin_arm64/
```

### Step 2: Copy Your Provider
```bash
cp /Users/volodymyrl/.asdf/installs/golang/1.24.2/bin/terraform-provider-redshift ~/.terraform.d/plugins/registry.terraform.io/serenityzn/redshift/1.0.0/darwin_arm64/terraform-provider-redshift_v1.0.0
```

### Step 3: Remove .terraformrc
```bash
rm ~/.terraformrc
```

### Step 4: Test
```bash
cd /Users/volodymyrl/ODDITY/projects/infra-terragrunt-code/live/brand3/stg/us-east-1/datateam/glue-provider
terragrunt plan --non-interactive
```

This approach:
- ✅ **Completely bypasses registry queries**
- ✅ **Works with Terragrunt**
- ✅ **Uses your local provider**
- ✅ **Is the standard way to test providers locally**