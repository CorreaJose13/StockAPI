module "lambda_role" {
  source = "../../modules/iam/iam_role/"

  name     = "stock-api-lambda-role"
  services = ["lambda.amazonaws.com"]
}

module "lambda_logs_policy" {
  source = "../../modules/iam/iam_policy_attachment/"

  name        = "stock-api-lambda-policies"
  description = "IAM policy for Lambda function to write logs to CloudWatch"
  action      = ["logs:CreateLogStream", "logs:PutLogEvents"]
  resource    = "arn:aws:logs:*:*:*"
  role_name   = module.lambda_role.name
}


