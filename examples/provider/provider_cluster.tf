provider "redshift" {
  host     = "my-cluster.123456789.us-east-1.redshift.amazonaws.com"
  username = var.redshift_user
  password = var.redshift_password
  database = "dev"
  # type parameter omitted for regular Redshift clusters
}

