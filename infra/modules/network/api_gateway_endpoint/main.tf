data "aws_api_gateway_rest_api" "this" {
  name = var.api_gateway_name
}

data "aws_lambda_function" "this" {
  function_name = var.lambda_name
}

resource "aws_api_gateway_resource" "this" {
  rest_api_id = data.aws_api_gateway_rest_api.this.id
  parent_id   = data.aws_api_gateway_rest_api.this.root_resource_id
  path_part   = var.path
}

resource "aws_api_gateway_method" "this" {
  rest_api_id   = data.aws_api_gateway_rest_api.this.id
  resource_id   = aws_api_gateway_resource.this.id
  http_method   = var.method
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "this" {
  rest_api_id = data.aws_api_gateway_rest_api.this.id
  resource_id = aws_api_gateway_resource.this.id

  http_method             = aws_api_gateway_method.this.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = data.aws_lambda_function.this.invoke_arn
}

resource "aws_lambda_permission" "this" {
  action        = "lambda:InvokeFunction"
  function_name = data.aws_lambda_function.this.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${data.aws_api_gateway_rest_api.this.execution_arn}/${var.stage}/${var.method}/${var.path}"
}

resource "aws_api_gateway_deployment" "this" {
  rest_api_id = data.aws_api_gateway_rest_api.this.id

  triggers = {
    redeployment = sha1(jsonencode([
      aws_api_gateway_resource.this.id,
      aws_api_gateway_method.this.id,
      aws_api_gateway_integration.this.id,
    ]))
  }

  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "this" {
  deployment_id = aws_api_gateway_deployment.this.id
  rest_api_id   = data.aws_api_gateway_rest_api.this.id
  stage_name    = var.stage
}

// TO DO: Add CORS configuration
# resource "aws_api_gateway_method" "options" {
#   rest_api_id   = data.aws_api_gateway_rest_api.this.id
#   resource_id   = aws_api_gateway_resource.this.id
#   http_method   = "OPTIONS"
#   authorization = "NONE"
# }
