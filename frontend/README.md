# StockWise Frontend Installation Guide

### Prerequisites

- Node.js 18.19.1 or higher
- npm 9.2.0 or higher
- AWS API Gateway invoke url

## Enviroment Setup

1. Clone the repository:

```sh
git clone https://github.com/CorreaJose13/StockWise.git
cd StockWise/frontend
```

2. Install dependencies:

```sh
npm install
```

2. Create a `.env` file in the root directory with the following variables:

```sh
VITE_API_BASE_URL=apigateway-base-url(example:https://apigateway-id.region.amazonaws.com/stage)
```

3. Compile and Hot-Reload for Development

```sh
npm run dev
```

4. Type-Check, Compile and Minify for Production

```sh
npm run build
```

### Lint with [ESLint](https://eslint.org/)

```sh
npm run lint
```
