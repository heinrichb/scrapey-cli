# ✨ Scrapey CLI

[![Build & Test](https://github.com/heinrichb/scrapey-cli/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/heinrichb/scrapey-cli/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/heinrichb/scrapey-cli.svg)](https://pkg.go.dev/github.com/heinrichb/scrapey-cli)
[![Coverage Status](https://img.shields.io/badge/coverage-0%25-red)](https://example.com/coverage)

Scrapey CLI is a lightweight, configurable web crawler and scraper. It collects data from websites based on rules defined in a config file. It can handle HTML parsing, data extraction, and plans to offer multiple storage options (JSON, XML, Excel, databases, etc.).

---

## 🚀 Features

- Lightweight and modular CLI interface
- Configurable input (`.json` config file or command-line flags)
- Extensible parsing logic for targeted HTML elements
- Future support for multiple storage options (JSON, XML, Excel, MongoDB, MySQL)
- DRY and clean code principles

---

## 🌱 Getting Started

1. **Clone the repo**:
   git clone https://github.com/heinrichb/scrapey-cli.git

2. **Initialize Go modules**:
   cd scrapey-cli
   go mod tidy

3. **Build the CLI**:  
   Use the provided Makefile to build the binary into the `build` folder:
   make build
   Alternatively, you can run:
   go build -o build/scrapeycli ./cmd/scrapeycli

4. **Run**:
   ./build/scrapeycli --config configs/default.json

---

## ⚙️ Project Structure

```
scrapey-cli/
├── .github/                 # GitHub-specific configurations
│   └── workflows/
│       └── ci.yml           # GitHub Actions CI/CD pipeline configuration
├── .vscode/                 # VS Code settings
│   └── settings.json        # Editor settings (format on save for Go)
├── build/                   # Build scripts, Dockerfiles, etc.
├── cmd/
│   └── scrapeycli/          # CLI application entry point
│       └── main.go          # Main Go file for Scrapey CLI
├── configs/
│   └── default.json         # Default/example configuration file
├── docs/                    # Project documentation
├── pkg/                     # Public packages for external use
│   ├── config/
│   │   └── config.go        # Config loading logic (stubbed)
│   ├── crawler/
│   │   └── crawler.go       # Core web crawling logic (stubbed)
│   ├── parser/
│   │   └── parser.go        # HTML parsing logic (stubbed)
│   └── storage/
│       └── storage.go       # Storage logic (stubbed for JSON and others)
├── test/                    # Optional integration tests
├── .gitignore               # Git ignore file
├── LICENSE                  # MIT License file
├── Makefile                 # Build script to output binary into build folder
├── go.mod                   # Go module file
├── go.sum                   # Go module checksum file
└── README.md                # Project README
```

---

## 🛠 Usage

- **Basic**:
  ./build/scrapeycli --url https://example.com

- **With config file**:
  ./build/scrapeycli --config configs/default.json

- **Future**:
  - Save data to JSON
  - Multiple URLs at once
  - Concurrency and rate-limiting

---

## 🧪 Tests

- Run unit tests locally:
  go test ./...

- Automated tests on GitHub Actions:
  - Triggered on every push and pull request to main or develop branches.
  - See Build & Test (https://github.com/heinrichb/scrapey-cli/actions) for logs and results.

---

## 🤝 Contributing

1. Fork the project
2. Create your feature branch (git checkout -b feature/amazing-feature)
3. Commit your changes (git commit -m 'Add some amazing feature')
4. Push to the branch (git push origin feature/amazing-feature)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License (LICENSE).

---

Happy Scraping!
