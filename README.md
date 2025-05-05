<h1 align="center" id="title">Stock Wise</h1>

<p id="description">StockWise is a web application for visualizing and analyzing stock market analyst ratings. The app fetches stock data from an external API, processes it, stores it in CockRoachDB, and presents it to users through an intuitive UI built with Vue + Tailwind, using the library PrimeVue. The backend is developed in Go, providing efficient data processing capabilities. All required API endpoints consumed by the frontend are deployed on AWS infrastructure, with the entire cloud resources managed through Infrastructure as Code (IaC) using Terraform.</p>

## Table of Contents

- [Demo](#demo)
- [Features](#features)
- [Installation steps](#installation-steps)
- [Built with](#built-with)
- [Project structure](#project-structure)
- [Useful resources](#useful-resources)
- [License](#license)
- [Author](#author)

## Demo

👉 **[Try StockWise Now](https://stock-wise-khaki.vercel.app)** 👈

## Features

Here're some of the project's best features:

- 📊 **Metrics Dashboard**: Displays total stock ratings tracked and their growth trends over time.
- 📰 **Recent Ratings Feed**: Stay updated with the latest market analyses and stock ratings.
- 🌟 **Top Stocks Carousel**: Highlights the top suggested stocks to invest in.
- 🧮 **Scoring Algorithm**: Identifies the best investment opportunities based on multiple calculations.
- 📋 **Interactive Ratings Table**: Filter by ticker, company, or analyst, and sort by ticker, company, analyst, or date.
- 🔍 **Detailed Stock View**: Dive deeper into stock information with an in-depth modal view.

## Installation steps

Each folder (`backend`, `frontend`, `infra`) contains a dedicated `README` file with detailed instructions to help you set up and run the respective part of the project locally. Check them out!

## Built with

Main technologies used in this project:

| **Technology**   | **Version** | **Description**                                                        |
| ---------------- | ----------- | ---------------------------------------------------------------------- |
| **Go**           | `v1.24.2`   | Backend programming language for efficient data processing.            |
| **TypeScript**   | `v5.8.0`    | Strongly typed programming language that builds on JavaScript.         |
| **Vue**          | `v3.5.13`   | Frontend framework for building a dynamic and reactive user interface. |
| **Tailwind CSS** | `v4.1.4`    | Utility-first CSS framework for fast UI development.                   |
| **PrimeVue**     | `v4.3.3`    | Vue component library for UI elements.                                 |
| **Node.js**      | `v18.19.1`  | JavaScript runtime for server-side and build processes.                |
| **npm**          | `v9.2.0`    | Package manager for JavaScript dependencies.                           |
| **AWS CLI**      | `v2.15.58`  | Command-line interface for managing AWS services.                      |
| **Terraform**    | `v1.11.4`   | Infrastructure as Code (IaC) tool for cloud resource management.       |
| **CockroachDB**  | `v25.1.6`   | Distributed SQL database for scalable and resilient data storage.      |

## Project structure

```bash
StockWise/
├── backend/ # Go server-side code
│   ├── config/ # Configuration management
│   ├── internal/ # Internal packages
│   │   ├── analysis/ # Stock analysis algorithms
│   │   ├── api/ # API consumer and response handlers
│   │   ├── db/ # Database interactions
│   │   ├── functions/ # Lambda function handlers
│   │   └── repository/ # Db repository pattern implementation
│   ├── models/ # Data models
│   └── utils/ # Utility functions
├── frontend/ # Vue.js client-side code
│   ├── public/ # Static assets
│   └── src/ # Source files
│       ├── components/ # Reusable Vue components
│       ├── composables/ # Composable functions
│       ├── types/ # TypeScript type definitions
│       ├── utils/ # Utility functions
│       └── views/ # Page components
└── infra/ # Terraform infrastructure code
    ├── environments/ # Environment-specific configurations
    └── modules/ # Reusable Terraform modules
```

## Useful resources

- [Quickstart with CockroachDB](https://www.cockroachlabs.com/docs/cockroachcloud/quickstart)
- [PrimeVue Docs](https://primevue.org/vite)
- [AWS with Terraform](https://registry.terraform.io/providers/hashicorp/aws/latest/docs)

## License

> This project is licensed under the MIT License - see the LICENSE file for details.

## Author

<a href="https://github.com/CorreaJose13/StockWise/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=CorreaJose13/StockWise" />
</a>
