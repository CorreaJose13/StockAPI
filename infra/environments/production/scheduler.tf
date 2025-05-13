module "update_scheduler" {
  source             = "../../modules/lambda_scheduler_integration/"
  lambda_source_path = "${path.module}/../../../backend/internal/functions/schedule/main.go"
  function_name      = "update-db-lambda"
  s3_bucket          = module.lambda_bucket.bucket
  lambda_role        = module.lambda_role.arn
  timeout            = 100
  memory_size        = 128
  log_retention_days = 7

  env_vars = {
    DB_URL       = var.DB_URL
    API_URL      = var.API_URL
    BEARER_TOKEN = var.BEARER_TOKEN
  }

  schedule_expression     = "cron(35 19 ? * mon-fri *)" # Every hour from 8 AM to 3 PM, Monday to Friday
  scheduler_name          = "update_db_scheduler"
  lambda_scheduler_role   = var.lambda_scheduler_role
  lambda_scheduler_policy = var.lambda_scheduler_policy
}
