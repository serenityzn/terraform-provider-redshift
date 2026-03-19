provider "redshift" {
  host     = var.redshift_host
  username = var.redshift_user
  type     = "serverless"
  temporary_credentials_serverless {
    workgroup_name = "my-workgroup"
  }
}
