module "lambda_bucket" {
  source = "../../modules/storage/s3"

  name          = var.lambda_bucket
  force_destroy = true
}
