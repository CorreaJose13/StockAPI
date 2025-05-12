variable "api_gateway_name" {
  description = "Name of the API Gateway"
  type        = string
  default     = "stock-api"
}

variable "lambda_bucket" {
  description = "S3 bucket for Lambda functions"
  type        = string
}

variable "lambda_role" {
  description = "IAM role for Lambda functions"
  type        = string
}

variable "lambda_logs_policy" {
  description = "CloudWatch log policy for Lambda functions"
  type        = string
}

variable "lambda_scheduler_role" {
  description = "IAM role for AWS Scheduler to invoke Lambda functions"
  type        = string
}

variable "lambda_scheduler_policy" {
  description = "IAM policy for AWS Scheduler to invoke Lambda functions"
  type        = string
}

variable "DB_URL" {
  description = "CockRoachDB connection string"
  type        = string
  sensitive   = true
}

variable "API_URL" {
  description = "External API URL"
  type        = string
  sensitive   = true
}

variable "API_KEY" {
  description = "API key for external API"
  type        = string
  sensitive   = true
}

variable "BEARER_TOKEN" {
  description = "Bearer token for external API"
  type        = string
  sensitive   = true
}

variable "stage" {
  description = "Stage of the API Gateway"
  type        = string
}
