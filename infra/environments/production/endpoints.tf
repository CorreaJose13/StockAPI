variable "stage" { default = "production" }

module "get_metrics_endpoint" {
  source             = "../../modules/lambda_api_integration/"
  lambda_source_path = "${path.module}/../../../backend/internal/functions/metrics/main.go"
  s3_bucket          = module.lambda_bucket.bucket
  lambda_role        = module.lambda_role.arn
  env_vars           = { DB_URL = var.DB_URL }

  endpoint_name    = "get-metrics"
  api_gateway_name = var.api_gateway_name
  endpoint_path    = "get-metrics"
  stage            = var.stage
}

module "analyze_endpoint" {
  source             = "../../modules/lambda_api_integration/"
  lambda_source_path = "${path.module}/../../../backend/internal/functions/analysis/main.go"
  s3_bucket          = module.lambda_bucket.bucket
  lambda_role        = module.lambda_role.arn
  env_vars           = { DB_URL = var.DB_URL }

  endpoint_name    = "analyze"
  api_gateway_name = var.api_gateway_name
  endpoint_path    = "analyze"
  stage            = var.stage
}

// TO DO: Improve redeployment strategy
resource "aws_api_gateway_deployment" "deployment" {
  rest_api_id = module.api_gateway.id

  triggers = {
    redeployment = sha1(jsonencode([
      module.get_metrics_endpoint,
      module.analyze_endpoint,
    ]))
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "stage" {
  deployment_id = aws_api_gateway_deployment.deployment.id
  rest_api_id   = module.api_gateway.id
  stage_name    = var.stage
}

output "invoke_url" {
  value = aws_api_gateway_stage.stage.invoke_url
}
