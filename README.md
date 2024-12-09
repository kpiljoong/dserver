## DServer

DServer is a lightweight, configurable, and dynamic mock server built in Go. It allows developers to simulate RESTful APIs for testing and development purposes, using a simple configuration file.

### Features
- **Dynamic Routing**: Easily define API routes with support for different methods like GET, POST, etc.
- **JSON Response**: Serve customizable JSON responses out of the box.
- **Delay Simulation**: Simulate delays to test timeout handling in clients.
- **Hot Reloading**: Automa

### Prerequisites
- **Go**: Ensure you have Go (version 1.21 or higher) installed. [Install Go](https://golang.org/doc/install).

### Installation
1. Clone the repository:
```bash
git clone https://github.com/kpiljoong/dserver.git
cd dserver
```
2. Build the binary:
```bash
go build -o dserver
```
3. Alternatively, run directly:
```bash
go run main.go
```

### Usage
#### Configuration File (config.toml)
Define your API endpoints in a `config.toml` file. Example:
```toml
[[routes]]
path = "/api/v1/health"
method = "GET"
status_code = 200
response_body = """{
  "status": "healthy"
}"""
content_type = "application/json"
```

#### Running the Server
Start the server with:
```bash
./dserver -config=config.toml -port=8080
```

Flags:
- `-config`: Path to the configuration file (default: `config.toml`).
- `-port`: Port to run the server on (default: `8080`).
- `-verbose`: Enable verbose logging (default: `false`).

Example:
```bash
./dserver -config=custom-config.toml -port=9090 -verbose
```

### Hot Reloading
DServer automatically reloads the configuration file when it is updated. Simply edit the `config.toml` file, and the changes will be applied without restarting the server.

### Development
#### Running Tests
Run the test suite:
```bash
go test ./...
```

### Adding New Features
Contributions are welcome! Feel free to open a pull request or create an issue for feature requests.

### License
This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

### Acknowledgements
- Built with [Go](https://golang.org/).
