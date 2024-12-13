## DServer

DServer is a lightweight, configurable, and dynamic mock server built in Go. It allows developers to simulate RESTful APIs for testing and development purposes, using a simple configuration file.

### Features
- **Dynamic Routing**: Easily define API routes with support for different methods like GET, POST, etc.
- **JSON Response**: Serve customizable JSON responses out of the box.
- **Delay Simulation**: Simulate delays globally or at the per-response level to test timeout handling in clients.
- **Hot Reloading**: Automatically reloads configuration when the file is updated, ensuring seamless development.
- Fallback Logic:
  - Use a global `response_body` if no `responses` array is provided.
  - Apply global defaults for `status_code` and `content_type` and `delay_ms` where not explicitly overridden in responses.

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
path = "/api/v1/resource"
method = "GET"
status_code = 200                   # Global status code
delay_ms = 500                      # Global delay in milliseconds
content_type = "application/json"   # Global content type
response_body = """{
  "message": "Fallback response"
}"""

[[routes]]
path = "/api/v1/resource"
method = "POST"
status_code = 201
content_type = "application/json"
responses = [
{ query = {}, response_body = """{ "message": "Default resources" }""" },
{ query = { "type" = "active" }, delay_ms = 1000, response_body = """{ "message": "Active resources" }""" },
{ query = { "type" = "inactive" }, response_body = """{ "message": "Inactive resources" }""" }
]
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

### Features in Detail

#### Global Defaults
DServer supports default values at the `[[routes]]` level:
- `status_code`: Used for all responses unless explicitly overridden.
- `content_type`: Global content type for the route.
- `delay_ms`: Default delay applied if not overridden in individual response.

#### Per-Response Overrides
Individual responses can override the global defaults:
- `status_code`
- `content_type`
- `delay_ms`

#### Fallback to `response_body`
If no `responses` array is provided, the server falls back to the `response_body` defined at the `[[routes]]` level.

#### Hot Reloading
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
