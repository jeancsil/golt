# Golt: Go Load Tester

A simple, lightweight, and concurrent HTTP load testing tool written in Go.

## 🚀 Features

- **Concurrent Requests**: Control the number of concurrent workers.
- **Configurable**: Set target URL, total requests, and concurrency via command-line flags.
- **Simple Output**: View status codes and timing durations for your tests.

## 🛠️ Usage

### Prerequisites

- [Go](https://go.dev/dl/) installed on your machine.

### Running the Tool

You can run the tool directly using `go run`:

```bash
go run main.go [flags]
```

### Configuration Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-u` | The target URL to test | `http://google.com` |
| `-n` | Total number of requests to make | `10` |
| `-c` | Number of concurrent workers | `1` |

### Examples

**Basic run with defaults:**
```bash
go run main.go
```

**Custom target with high concurrency:**
Run 100 requests to `https://api.example.com` using 10 concurrent workers:
```bash
go run main.go -u https://api.example.com -n 100 -c 10
```

## 🤝 Contributing

Contributions, issues, and feature requests are welcome!

## 📄 License

This project is open source and available under the [MIT License](LICENSE).
