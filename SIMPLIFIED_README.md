# migrator

A lightweight file conversion tool that converts between various data formats including Excel, CSV, JSON, and SQL.

## ✨ Features

- Convert Excel files (.xlsx) to CSV, JSON, or SQL
- Convert CSV files to Excel, JSON, or SQL
- Convert JSON files to Excel, CSV, or SQL
- Convert SQL files to Excel, CSV, or JSON
- Simple web interface for uploading and converting files
- No external dependencies required (except Go)

## 🚀 Quick Start

### Prerequisites

- Go 1.24 or higher

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/lukaszbudnik/migrator.git
   cd migrator
   ```

2. Run the application:
   ```bash
   go run migrator.go
   ```

3. Open your browser and navigate to `http://localhost:8080`

### Usage

1. Select a file to convert (Excel, CSV, JSON, or SQL)
2. Choose the target format from the dropdown
3. Click "Start Migration" to process the file
4. Download the converted file using the "Download Migrated File" button

## 🛠️ Supported Conversions

| From → To | CSV | JSON | SQL | Excel |
|-----------|-----|------|-----|-------|
| CSV       | -   | ✓    | ✓   | ✓     |
| JSON      | ✓   | -    | ✓   | ✓     |
| SQL       | ✓   | ✓    | -   | ✓     |
| Excel     | ✓   | ✓    | ✓   | -     |

## 📁 Project Structure

```
migrator/
├── converter/          # File conversion logic
├── server/             # Web server implementation
├── static/             # Frontend files (HTML, CSS, JS)
├── uploads/            # Directory for uploaded files
└── migrator.go         # Main application entry point
```

## 🔧 Configuration

The application runs on port 8080 by default. To change the port, modify the configuration in `migrator.go`.

## 📄 License

This project is licensed under the Apache 2.0 License - see the [LICENSE](LICENSE) file for details.