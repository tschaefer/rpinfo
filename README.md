# rpinfo

**rpinfo** is a lightweight RESTful API server written in Go that exposes
detailed system information for Raspberry Pi devices. It utilizes the
`vcgencmd` utility to provide real-time hardware data such as CPU temperature,
voltages, firmware configuration, throttling status and clock frequencies.

## Features

- Exposes Raspberry Pi system metrics via a clean RESTful API
- Supports optional bearer token authentication
- Configurable host and port via command-line flags
- Fast and efficient Go implementation
- Ideal for integration with dashboards, monitoring tools, or automation scripts

## Getting Started

### Prerequisites

- Raspberry Pi running a Linux-based OS (e.g., Raspberry Pi OS)
- `vcgencmd` utility (preinstalled on Raspberry Pi OS)

### Installation and Usage

Download the latest release from the [releases page](https://github.com/tschaefer/rpinfo/releases).

Start the server on `localhost:8080` by default.

```bash
./rpinfo server
```
For further configuration, see the command-line options below.

| Flag            | Description                          | Default     |
|-----------------|--------------------------------------|-------------|
| `-H`, `--host`  | Host to bind the server to           | `localhost` |
| `-p`, `--port`  | Port to run the server on            | `8080`      |
| `-a`, `--auth`  | Enable bearer token authentication   | `false`     |
| `-t`, `--token` | Bearer token used for authentication |             |
| `-h`, `--help`  | Show help for the server command     |             |

Additional a systemd service file and environment file are provided in the
[contrib directory](https://github.com/tschaefer/rpinfo/tree/main/contrib) for automatic startup on boot and management of the
server.

## API Endpoints

| Endpoint                  | Description                    |
|---------------------------|--------------------------------|
| `/configuration`          | Returns firmware configuration |
| `/temperature`            | Returns CPU temperature        |
| `/throttled(?human=true)` | Returns throttling status      |
| `/voltages`               | Returns voltages               |
| `/clock`                  | Returns clock frequencies      |

All endpoints return JSON-formatted data.

The complete API specification is available at `/redoc`.

## Security Notes

- If authentication is enabled, all API calls must include the `Authorization`
header with the valid bearer token.
- Use strong and random tokens.
- Consider running the server behind HTTPS if exposed publicly.

## Contributing

Contributions are welcome! Please fork the repository and submit a pull request.
For major changes, open an issue first to discuss what you would like to change.

Ensure that your code adheres to the existing style and includes appropriate tests.

## License

This project is licensed under the [MIT License](LICENSE).
