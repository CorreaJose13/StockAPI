locals {
  zip_file = "${var.function_name}.zip"
}

resource "terraform_data" "this" {
  triggers_replace = {
    always_run = "${timestamp()}"
  }
  provisioner "local-exec" {
    working_dir = dirname(var.lambda_source_path)
    command     = "make publish"
    environment = {
      BUCKET_NAME = var.s3_bucket
      BUILD_NAME  = local.zip_file
    }
  }

  depends_on = [var.s3_bucket]
}

module "lambda_function" {
  source = "../compute/lambda/"

  s3_bucket = var.s3_bucket
  s3_key    = local.zip_file

  function_name = "stock-api-${var.function_name}"
  description   = "Lambda function for ${var.function_name}"
  role          = var.lambda_role

  runtime            = var.runtime
  handler            = var.handler
  timeout            = var.timeout
  memory_size        = var.memory_size
  log_retention_days = var.log_retention_days
  env_vars           = var.env_vars

  depends_on = [terraform_data.this]
}

resource "aws_scheduler_schedule" "this" {
  name = var.scheduler_name

  flexible_time_window {
    mode = "OFF"
  }

  schedule_expression          = var.schedule_expression
  schedule_expression_timezone = "America/Bogota"

  target {
    arn      = "arn:aws:scheduler:::aws-sdk:lambda:invoke"
    role_arn = module.lambda_function.arn
  }
}
