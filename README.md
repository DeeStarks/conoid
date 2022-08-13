# CONOID

[![Release](https://github.com/DeeStarks/conoid/actions/workflows/release.yml/badge.svg)](https://github.com/DeeStarks/conoid/actions/workflows/release.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![GitHub tag](https://img.shields.io/github/tag/deestarks/conoid.svg)](https://github.com/deestarks/conoid/releases/latest)


<img src="./assets/welcome/img/icon.png" width="100">

---

Conoid is a simple HTTP server that can serve static files as well as a TCP tunneling tool that uses [localtunnel](http://localtunnel.me/) to bypass a firewall or NAT and expose your local development server to the internet.

While you can expose as many local servers as you want, conoid connects to the [localtunnel](http://localtunnel.me/) server only on a single port - **5000**, and directs traffic to your servers based on a mapping created while making the initail connections.

## Architecture

<img src="./assets/imgs/architecture.png">

## Installation
### Homebrew
*[In progress...]*

### Go
```
go install github.com/deestarks/conoid@latest && \
export PATH=$PATH:$(go env GOPATH)/bin
```

## Usage

### Start conoid server
To start conoid server in the foreground, execute:
```
conoid
```
Go to http://127.0.0.1:5000 to view the welcome page.

### List services
- Running services
```
conoid ps
```

- All services
```
conoid ps -a
```

### Expose a local server
```
conoid add \
    --name <your_app_name> --type server \
    --listener <your_server_address> --tunnel
```

E.g.
```
conoid add \
    --name my_app --type server \
    --listener <your_server_address> --tunnel
```

### Serving static files
```
conoid add \
    --name <your_app_name> --type static \
    --directory <document_directory>
```

or add the `--tunnel` flag to expose to the internet


Use the help flag `conoid [command] --help` or `-h` for more commands.
