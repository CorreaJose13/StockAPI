locals {
  src_path = "${path.module}/../../../backend/internal/functions/metrics/main.go"
}

variable "metrics_key" { default = "lambda_metrics.zip" }

module "lambda_bucket" {
  source = "../../modules/storage/s3"

  name          = var.lambda_bucket
  force_destroy = true
}

resource "terraform_data" "function_binary" {

  provisioner "local-exec" {
    working_dir = dirname(local.src_path)
    command     = "make publish"
    environment = {
      BUCKET_NAME = module.lambda_bucket.bucket
      BUILD_NAME  = var.metrics_key
    }
  }

  depends_on = [module.lambda_bucket]
}

module "lambda_function" {
  source = "../../modules/compute/lambda/"

  s3_bucket = module.lambda_bucket.bucket
  s3_key    = var.metrics_key

  function_name = "stock-api-get-metrics"
  description   = "Lambda function to get metrics from the database"
  role          = module.lambda_role.arn

  env_vars = {
    "DB_URL" = var.DB_URL
  }

  depends_on = [terraform_data.function_binary]
}

module "api_gateway_endpoint" {
  source           = "../../modules/network/api_gateway_endpoint/"
  api_gateway_name = var.api_gateway_name
  path             = "get-metrics"
  method           = "GET"
  stage            = "production"
  lambda_name      = module.lambda_function.name

  depends_on = [module.api_gateway, module.lambda_function]
}

output "endpoint_url" {
  value       = module.api_gateway_endpoint.endpoint_url
  description = "The invoke URL for the API Gateway endpoint"
}
