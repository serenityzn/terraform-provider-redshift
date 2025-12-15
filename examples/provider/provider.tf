provider "redshift" {
  host     = "your-workgroup.123456789.us-east-1.redshift-serverless.amazonaws.com"
  username = var.redshift_user
  password = var.redshift_password
  database = "dev"
  type     = "serverless"  # Required for Redshift Serverless
}
