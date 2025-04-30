variable "api_gateway_name" {
  description = "Name of the API Gateway"
  type        = string
  default     = "stock-api"
}

variable "lambda_bucket" {
  description = "S3 bucket for Lambda functions"
  type        = string
}

variable "DB_URL" {
  description = "CockRoachDB connection string"
  type        = string
  sensitive   = true
}
