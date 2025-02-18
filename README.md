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
│   └── settings.json                 # VS Code settings (format on save for Go)
├── cmd/
│   └── scrapeycli/
│       └── main.go
├── configs/
│   └── default.json                  # Default/example configuration file
├── pkg/
│   ├── config/
│   │   └── config.go                 # Config loading logic
│   ├── crawler/
│   │   └── crawler.go                # Core web crawling logic
│   ├── parser/
│   │   └── parser.go                 # HTML parsing logic
│   ├── storage/
│   │   └── storage.go                # Storage logic
│   └── utils/
│       ├── printcolor.go             # Colorized terminal output utility
│       └── printstruct.go            # Utility for printing non-empty struct fields
├── scripts/
│   └── coverage_formatter.go         # Formats and colorizes Go test coverage output
├── test/                             # Optional integration tests
│   └── fail_test.go                  # Test case designed to always fail, used to debug test output
├── .gitignore
├── LICENSE                           # MIT License file
├── Makefile                          # Build & run script for CLI (includes targets for build, run, and test)
├── go.mod
├── go.sum
└── README.md
```

---

## 🔧 Configuration Options

Scrapey CLI is configured using a JSON file that defines how websites are crawled and scraped. Below is a detailed breakdown of the available configuration options.

### 🌍 URL Configuration

```json
"url": {
  "base": "https://example.com",
  "routes": [
    "/route1",
    "/route2",
    "*"
  ],
  "includeBase": false
}
```

- **base**: The primary domain to scrape.
- **routes**: List of specific paths to scrape. Supports `*` as a wildcard for full site crawling.
- **includeBase**: Whether to include the base URL in the scrape.

### 🔍 Parsing Rules

```json
"parseRules": {
  "title": "title",
  "metaDescription": "meta[name='description']",
  "articleContent": "article",
  "author": ".author-name",
  "datePublished": "meta[property='article:published_time']"
}
```

- **title**: Extracts the page title.
- **metaDescription**: Extracts the meta description.
- **articleContent**: Defines the main article section.
- **author**: Selector for extracting author names.
- **datePublished**: Extracts the publication date from meta properties.

### 💾 Storage Options

```json
"storage": {
  "outputFormats": ["json", "csv", "xml"],
  "savePath": "output/",
  "fileName": "scraped_data"
}
```

- **outputFormats**: List of formats in which data will be stored.
- **savePath**: Directory where scraped content is saved.
- **fileName**: Base name for output files.

### ⚡ Scraping Behavior

```json
"scrapingOptions": {
  "maxDepth": 2,
  "rateLimit": 1.5,
  "retryAttempts": 3,
  "userAgent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
}
```

- **maxDepth**: Defines how deep the scraper should follow links.
- **rateLimit**: Time delay (in seconds) between requests to avoid rate-limiting.
- **retryAttempts**: Number of retries for failed requests.
- **userAgent**: Custom user-agent string to mimic a browser.

### 🛠 Data Formatting

```json
"dataFormatting": {
  "cleanWhitespace": true,
  "removeHTML": true
}
```

- **cleanWhitespace**: Removes unnecessary whitespace in extracted content.
- **removeHTML**: Strips HTML tags from extracted content for cleaner output.

This configuration file allows fine-tuning of scraping behavior, data extraction, and storage formats for ultimate flexibility in web scraping.

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
