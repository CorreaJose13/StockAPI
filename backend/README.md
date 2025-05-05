# StockWise Backend Installation Guide

### Prerequisites

- Go 1.24.2
- AWS CLI configured with appropiate permissions
- CockroachDB cluster instance

### Environment Setup

1. Clone the repository:

```sh
git clone https://github.com/CorreaJose13/StockWise.git
cd StockWise/backend
```

2. Create a `.env` file in the root directory with the following variables:

```sh
DB_URL=your-cockroachdb-connection-string
API_URL=external-api-given-url
BEARER_TOKEN=external-api-given-token
```

3. Install dependencies:

```sh
go mod download
```

### Local Development

To run the application locally to fetch the stocks from the API and store them in CockroachDB:

```sh
go run cmd/stockapi/main.go
```

### Testing

Run the test availables:

```sh
go test ./...
```
