<h1 align="center">
  <br>
    corn
  <br>
  <br>
</h1>

<h4 align="center">A modern cronjob manager with logging, notification and remote worker support.</h4>
<p align="center">
  <a href="https://github.com/aigic8/corn/actions/workflows/docker-pusher.yaml"><img alt="Docker build & Push" src="https://github.com/aigic8/corn/actions/workflows/docker-pusher.yaml/badge.svg"></a>
  <a href="https://github.com/aigic8/corn/actions/workflows/go-releaser.yaml"><img alt="Go Releaser" src="https://github.com/aigic8/corn/actions/workflows/go-releaser.yaml/badge.svg"></a>
</p>
<br>

## Features

- **Structured Logging:** Execution logs of both tasks and agent itself are stored in files in json format.
- **Agentless Remote Execution:** Jobs can be be run on remote servers without installing corn on them. (through `ssh`)
- **Notification:** Notifications are support on `Telegram` and `Discord`
- **Fail Strategies:** Jobs can retried when failed

## Installation

### Single Line Script

Run the following command:

```sh
curl -s https://raw.githubusercontent.com/aigic8/corn/refs/heads/main/install/install.sh | sh
```

This script will do the following tasks:

- Download the latest release of corn
- Add `systemd` service for corn in your services

> [!NOTE]
> Single line script is only supported in **Linux** with **Systemd**.

### Docker Compose

Create a file named `docker-compose.yaml` and add the following content:

```yaml
services:
  corn:
    image: ghcr.io/aigic8/corn:0.2
    volumes:
      - ~/.corn:/root/.corn # TODO: REPLACE WITH YOUR INTENDED PATH ON YOUR SYSTEM
      - ~/.config/corn:/root/.config/corn # TODO: REPLACE WITH YOUR INTENDED PATH ON YOUR SYSTEM
```

Do not forget to modify the path to the volumes.

### Downloading Releases

Download the latest release of corn from [the releases page](https://github.com/aigic8/corn/releases/latest)

### Building Manually

Run the following command to clone the repo and build the project:

```sh
git clone https://github.com/aigic8/corn
cd corn
CGO_ENABLED=1 go build ./main.go
```

> [!NOTE]
> When downloading releases or building manually, you have to manually add
> corn to your systems daemon manager to run it in background.

## Configuration

Create a file named `corn.yaml` or `corn.yml` in `~/.config/corn/` directory and write the config file by modifying the following template:

```yaml
timezone: Europe/London # optional, system's timezone will be used if not provided
disableNotifications: false # disabled notifications on telegram and discord. DEFAULT: false
defaultNotifier: base # optional if disabledNotifications is true, default notifier if no notifier is provided for the job
defaultFailNotifier: base # default fail notifier if no fail notifier is provided for the job. If not provided, defaultNotifier will be used
notifyTimeoutMs: 5000 # optional, time out to send the notification in milliseconds
defaultTimeoutS: 5 # optional, default time out to run the jobs in seconds
# remotes are the worker servers to run the jobs on
remotes:
  my-remote-key:
    username: username # username of the remote
    address: 1.1.1.1
    port: 22
    auth:
      keyAuth:
        keyPath: /path/to/ssh/private.key # required, path to ssh private key
        passphrase: "" # optional, ssh key passphrase if set
  my-remote-password:
    username: username
    address: 1.1.1.1
    port: 22
    auth:
      passwordAuth:
        password: "12345" # required, password of the remote user on the server
jobs:
  localJob:
    schedules: # schedules, multiple schedules can be provided for the same job
      - "0 1 11 * * *" # similar to cron syntax, except the first item (0) is the second
      - "0 1 10 * * *"
    command: |
    onlyLogOnFail: false
    onlyNotifyOnFail: false # optional, wether only notify on failure
    notifier: base # optional, notifier of the job. If not provided defaultNotifier will be used.
    failNotifier: more # optional, fail notifier of the job.
    timeoutS: 10 # optional, timeout of running the job
    failStrategy:
      retry:
        maxRetries: 3
  remoteJob:
    schedules:
      - "* 1 11 * * *"
      - "* 1 10 * * *"
    command: whoami
    onlyLogOnFail: false
    remoteName: my-remote-key # name of the remote to run the job on
notifiers:
  base: # name of the notifier
    telegram:
      token: "TELEGRAM_TOKEN"
      receivers:
        - 1111111111 # chat id to send the message to, for private chats it is equivalent to userId
  more:
    discord:
      oAuth2Token: "OAUTH2_TOKEN" # you either have to pass oAuth2Token or botToken
      botToken: "botToken"
      channels: # required (at least 1 item), list of channel ids
        - CHANNEL_ID_1
        - CHANNEL_ID_2
```

## Usage

You can run the application with the following command:

```sh
corn run
```

A single job can be executed with the following command:

```sh
corn test -j jobName
```

You can also pass the `--dev` (or `-d`) option to printing the logs to stdout (shell) instead of pushing them to log files.

```sh
corn test -j jobName --dev
```
