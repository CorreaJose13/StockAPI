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

variable "API_URL" {
  description = "External API URL"
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
