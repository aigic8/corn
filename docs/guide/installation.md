# Installation

There are several ways to install corn. Choose the one that best fits your environment.

## Single Line Script

For **Linux** systems with **Systemd**, you can use the following command to install the latest release and set up a `systemd` service:

```sh
curl -s https://raw.githubusercontent.com/aigic8/corn/refs/heads/main/install/install.sh | sh
```

## Docker Compose

You can run corn using Docker. Create a `docker-compose.yaml` file with the following content:

```yaml
services:
  corn:
    image: ghcr.io/aigic8/corn:0.2
    volumes:
      - ~/.corn:/root/.corn # Path for corn's data (e.g., logs)
      - ~/.config/corn:/root/.config/corn # Path for corn's configuration file
```

Make sure to replace the volume paths with the intended paths on your host system.

The two volumes in docker compose files are used for following purposes:

- `/root/.corn`: used for log files and sqlite database (`corn.sqlite`)
- `/root/.config/corn`: used for configuration file (`corn.yaml`)

## Downloading Releases

You can download the latest binary for your operating system from the [releases page](https://github.com/aigic8/corn/releases/latest).

## Building Manually

To build corn from source, you'll need Go installed.

```sh
git clone https://github.com/aigic8/corn
cd corn
CGO_ENABLED=1 go build ./main.go
```

> [!NOTE]
> When downloading a release or building manually, you will need to set up a daemon or service manager (like `systemd` or `supervisord`) to run corn in the background.
