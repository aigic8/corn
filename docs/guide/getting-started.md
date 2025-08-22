# Getting started

## Installation

For a quick start, you can use the single-line installation script for Linux systems with Systemd.

```sh
curl -s https://raw.githubusercontent.com/aigic8/corn/refs/heads/main/install/install.sh | sh
```

This script will download the latest release and set up a `systemd` service for corn.

For other installation methods, such as Docker or manual installation, please see the [full installation guide](./installation.md).

## Basic Configuration

Create a configuration file at `~/.config/corn/corn.yaml` (or `~/.config/corn/corn.yml`) with a simple job to get started:

```yaml
jobs:
  localJob:
    schedules:
      - "0 1 11 * * *" # Runs at 11:01:00 every day
    command: echo "Hello from corn!"
```

This configuration defines a single job named `localJob` that runs at a specified schedule.

The scheduling syntax is like [cron syntax](https://cron.help/) only the first number is for seconds (`0` in this case)

For more detailed configuration options, see the [configuration guide](./configuration.md).

## Usage

To run corn with your configuration, use the following command:

```sh
corn run
```

You can also test a specific job:

```sh
corn test -j localJob
```
