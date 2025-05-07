# StockWise Infrastructure Installation Guide

### Prerequisites

- Terraform 1.11.4 or higher
- AWS CLI configured with appropiate permissions
- CockroachDB cluster instance
- S3 bucket for Terraform state management (this infra uses s3 backend)

### Environment Setup

1. Clone the repository:

```sh
git clone https://github.com/CorreaJose13/StockWise.git
cd StockWise/infra/environments/production
```

2. Create a `state.config` file in the directory with the following variables:

```sh
bucket=your-bucket-for-state
region=your-region
key=your-key
encrypt=true
use_lockfile=true
```

3. Create a `variables.tfvars` file in the directory with the following variables:

```sh
DB_URL           = your-cockroachdb-connection-string
api_gateway_name = your-api-gateway-name
lambda_bucket    = your-bucket-name-to-store-lambda-functions
stage            = your-stage
```

4.  Configure remote state:

```sh
terraform init -backend-config=state.config
```

5. Plan and apply:

```sh
terraform plan -var-file=variables.tfvars
terraform apply -var-file=variables.tfvars
```

> [!NOTE]
> The infrastructure contains a trigger that always updates the lambda functions defined in `backend` folder, so after code changes to the lambda functions just run `terraform apply` to deploy the updated code.
