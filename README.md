# CONOID

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![forthebadge](https://forthebadge.com/images/badges/made-with-go.svg)](https://forthebadge.com)

![conoid icon](./assets/welcome/img/icon.png)

---

Conoid is a simple HTTP server and TCP tunnelling tool that uses [localtunnel](http://localtunnel.me/) to bypass a firewall or NAT, to enable your local development server be exposed to the internet.

While you can expose as many local servers as you want, conoid connects to the [localtunnel](http://localtunnel.me/) server only on a single port - `5000`, and directs traffic to your servers based on a mapping created while making the inital connections.

## Architecture

![conoid architecture](https://drive.google.com/file/d/1Kl-I6rShVYIN8RxFXvMijDxWiPG3k9Xy/view?usp=sharing)

## Installation

*Working on it...*

## Usage

### List services
- Running services
```
conoid ps
```

- All services
```
conoid ps -a
```

### Expose a server
```
conoid add \
    --name your_app_name --type server \
    --listener your_server_addr --tunnel
```

### Serving static files
```
conoid add \
    --name your_app_name --type static \
    --directory document_directory
```

or add the `--tunnel` flag to expose to the internet


Use the help flag `conoid [command] --help` or `-h` to view more available commands.
