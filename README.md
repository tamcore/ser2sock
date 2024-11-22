# ser2sock

**ser2sock** is a lightweight application written in Go to expose a serial device over TCP/IP. It enables bidirectional communication between a serial device and a TCP client, with options for verbose logging and decoding traffic.

## Features
- Forward data between a serial port and a TCP connection.
- Configurable serial device, baud rate, and TCP listening address.
- Filter connections by allowed IP addresses.
- Optional verbose logging of incoming and outgoing data.
- Decode data into human-readable text (UTF-8) or display raw hexadecimal.

## Installation

### Prerequisites
- [Go](https://golang.org/) version 1.16 or later installed on your system.

### Build from Source
1. Clone
   ```bash
   git clone https://github.com/tamcore/ser2sock.git
   cd ser2sock
   ```
2. Build
   ```bash
   go build -o ser2sock
   ```

## Usage
```bash
./ser2sock -device <serial_device> -listen <address:port> -baud <baud_rate> [options]
```

## Options
| Option | Description | Example |
| ------ | ----------- | -------
| `-device` | Path to the serial device. Required. | `/dev/ttyUSB0`, `COM3`, `/dev/zigbee1` |
| `-listen` | TCP address and port to listen on. | `0.0.0.0:5000` |
| `-baud` | Baud rate for the serial device. Default: `9600`. | `115200` |
| `-allowed` | Comma-separated list of allowed client IPs. Leave empty to allow all IPs. | `192.168.1.100,192.168.1.101` |
| `-verbose` | Enable verbose logging for incoming and outgoing data.  (no value, just add the flag) |
| `-decode` | Attempt to decode data into human-readable UTF-8 text. Defaults to raw hexadecimal format. (no value, just add the flag) |

## Example Commands
### Basic Usage
Expose the serial device `/dev/ttyUSB0` on TCP port `5000` with a baud rate of `9600`:

```
./ser2sock -device /dev/ttyUSB0 -listen 0.0.0.0:5000 -baud 9600
```

### Restricting Access
Allow only clients from specific IPs:

```
./ser2sock -device /dev/ttyUSB0 -listen 0.0.0.0:5000 -baud 9600 -allowed "192.168.1.100,192.168.1.101"
```

### Verbose Logging
Enable detailed logging of IN/OUT traffic:

```
./ser2sock -device /dev/ttyUSB0 -listen 0.0.0.0:5000 -baud 9600 -verbose
```

### Decoding Data
Attempt to decode data into human-readable text when possible:

```
./ser2sock -device /dev/ttyUSB0 -listen 0.0.0.0:5000 -baud 9600 -verbose -decode
```

## Example Output
### Without Decoding (`-decode not used):
```
Accepted connection from 192.168.1.100:54832
IN  (TCP->Serial): 48656c6c6f
OUT (Serial->TCP): fe00
```

### With Decoding (`-decode enabled):
```

Accepted connection from 192.168.1.100:54832
IN  (TCP->Serial): "Hello"
OUT (Serial->TCP): (binary: fe00)
```

## Development
### Dependencies
* [go-serial](https://github.com/bugst/go-serial): For interacting with the serial device.

### Install dependencies:

```bash
go get ./...
```

### Testing
Run the application locally and connect using a TCP client like telnet or netcat:

```bash
telnet <server_ip> <port>
```
