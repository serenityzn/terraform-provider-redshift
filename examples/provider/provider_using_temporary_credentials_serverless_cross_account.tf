provider "redshift" {
  host     = var.redshift_host
  username = var.redshift_user
  type     = "serverless"
  temporary_credentials_serverless {
    workgroup_name = "my-workgroup"
    assume_role {
      arn = "arn:aws:iam::012345678901:role/role-name-with-path"
    }
  }
}
