<!-- File: README.md -->

# ✨ Scrapey CLI

[![Build & Test](https://github.com/heinrichb/scrapey-cli/actions/workflows/ci.yml/badge.svg?branch=develop)](https://github.com/heinrichb/scrapey-cli/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/heinrichb/scrapey-cli/branch/develop/graph/badge.svg?token=P45PASDIKF)](https://codecov.io/gh/heinrichb/scrapey-cli)
[![Go Reference](https://pkg.go.dev/badge/github.com/heinrichb/scrapey-cli.svg)](https://pkg.go.dev/github.com/heinrichb/scrapey-cli)

Scrapey CLI is a lightweight, configurable web crawler and scraper. It collects data from websites based on rules defined in a config file. It can handle HTML parsing, data extraction, and plans to offer multiple storage options (JSON, XML, Excel, databases, etc.).

---

## 🚀 Features

- **Lightweight & Modular CLI:** Built with clean, DRY code principles.
- **Configurable Input:** Accepts configuration via a JSON file or command-line flags.
- **Extensible Parsing:** Customizable HTML parsing logic.
- **Planned Storage Options:** Future support for multiple output formats including JSON, XML, Excel, MongoDB, MySQL.

---

## 🌱 Getting Started

1.  **Clone the Repo**

        git clone https://github.com/heinrichb/scrapey-cli.git

2.  **Initialize Go Modules & Build the CLI**

    - **Option 1:** Using the Makefile (recommended)

          make build

      - This command runs `go mod tidy` and then builds the binary into the `build` folder.

    - **Option 2:** Directly via Go

          go mod tidy
          go build -o build/scrapeycli ./cmd/scrapeycli

3.  **Run the CLI**

    - **Direct Execution:**

          ./build/scrapeycli --config configs/default.json

    - **Using the Makefile:**

      The Makefile provides a `run` target which allows you to pass in optional variables:

    - **Default Run:**

          make run

      - This uses the default configuration file (`configs/default.json`).

    - **Override Config:**

          make run CONFIG=configs/other.json

    - **Pass a URL:**

          make run URL=https://example.org

    - **Combined:**

          make run CONFIG=configs/other.json URL=https://example.org

---

## ⚙️ Project Structure

```
scrapey-cli/
├── .github/
│   └── workflows/
│       └── ci.yml
├── .vscode/
│   └── settings.json        # VS Code settings (format on save for Go)
├── cmd/
│   └── scrapeycli/
│       └── main.go
├── configs/
│   └── default.json         # Default/example configuration file
├── docs/                    # Project documentation
├── pkg/
│   ├── config/
│   │   └── config.go        # Config loading logic
│   ├── crawler/
│   │   └── crawler.go       # Core web crawling logic
│   ├── parser/
│   │   └── parser.go        # HTML parsing logic
│   ├── storage/
│   │   └── storage.go       # Storage logic
│   └── utils/
│       └── utils.go         # Utility functions
├── test/                    # Optional integration tests
├── .gitignore
├── LICENSE                  # MIT License file
├── Makefile                 # Build & run script for CLI (includes targets for build, run, and test)
├── go.mod
├── go.sum
└── README.md

```

---

## 🛠 Usage

- **Basic Execution:**

      ./build/scrapeycli --url https://example.com

- **With a Config File:**

      ./build/scrapeycli --config configs/default.json

- **Using the Makefile:**

  - Run with defaults:

        make run

  - Override configuration and/or URL:

        make run CONFIG=configs/other.json URL=https://example.org

- **Future Enhancements:**
  - Save scraped data to JSON.
  - Support for scraping multiple URLs simultaneously.
  - Concurrency and rate-limiting.

---

## 🧪 Tests

- **Run Unit Tests Locally:**
  To run tests for all modules and the test folder (if it exists), use:

      make test

  This command first runs "go test ./..." to execute tests in all packages, and then, if the "test" folder exists and contains Go files, it will run tests in that folder as well.

- **Automated Tests on GitHub Actions:**
  - Tests are triggered on every push and pull request to the "main" or "develop" branches.
  - See Build & Test (https://github.com/heinrichb/scrapey-cli/actions) for logs and results.

---

## 🤝 Contributing

1.  Fork the project.
2.  Create your feature branch:

        git checkout -b feature/amazing-feature

3.  Commit your changes:

        git commit -m 'Add some amazing feature'

4.  Push to the branch:

        git push origin feature/amazing-feature

5.  Open a Pull Request.

---

## 📄 License

This project is licensed under the MIT License ([LICENSE](LICENSE)).

---

Happy Scraping!
