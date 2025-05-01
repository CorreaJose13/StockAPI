terraform {
  backend "s3" {}
}

module "api_gateway" {
  source             = "../../modules/network/api_gateway/"
  name               = var.api_gateway_name
  description        = "Stock API for managing stock data"
  log_retention_days = 7
}
