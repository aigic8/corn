# Configuration

corn is configured using a YAML file located at `~/.config/corn/corn.yaml` or `~/.config/corn/corn.yml`.

The structure of configuration file is like following:

```yaml
timezone: Europe/London # optional, system's timezone will be used if not provided
disableNotifications: false # disabled notifications on telegram and discord. DEFAULT: false
defaultNotifier: base # optional if disabledNotifications is true, default notifier if no notifier is provided for the job
defaultFailNotifier: base # default fail notifier if no fail notifier is provided for the job. If not provided, defaultNotifier will be used
notifyTimeoutMs: 5000 # optional, time out to send the notification in milliseconds
defaultTimeoutS: 5 # optional, default time out to run the jobs in seconds
jobs: ...
# remotes are the worker servers to run the jobs on
remotes: ...
notifiers: ...
```

There are three main objects that can be defined:

- `jobs`: The jobs to be scheduled
- `remotes`: Servers to run the jobs on through ssh
- `notifiers`: Notification services used

For working examples check the [examples page](../examples/index).

## Generic Options

These settings are at the root of the configuration file.

- `timezone`: (Optional) Sets the timezone for scheduling. If not provided, the system's timezone is used.
- `disableNotifications`: (Optional) Set to `true` to disable all notifications. Defaults to `false`.
- `defaultNotifier`: The default notifier to use for jobs that don't specify one.
- `defaultFailNotifier`: The default notifier to use for failed jobs. If not set, `defaultNotifier` is used.
- `notifyTimeoutMs`: (Optional) Timeout in milliseconds for sending notifications.
- `defaultTimeoutS`: (Optional) Default timeout in seconds for jobs.

## Jobs

The `jobs` block contains the definitions for all the tasks you want to run.

- `schedules`: A list of cron-like schedules for the job. The format is `second minute hour day-of-month month day-of-week`.
- `command`: The command to be executed.
- `onlyLogOnFail`: (Optional) If `true`, logs are only written if the job fails.
- `onlyNotifyOnFail`: (Optional) If `true`, notifications are only sent if the job fails.
- `notifier`: (Optional) The notifier to use for this job. Overrides `defaultNotifier`.
- `failNotifier`: (Optional) The notifier to use for this job on failure. Overrides `defaultFailNotifier`.
- `timeoutS`: (Optional) Timeout in seconds for this specific job.
- `failStrategy`: (Optional) Defines what to do when a job fails. If not defined job will continue being scheduled normally.
  - `retry`:
    - `maxRetries`: Number of times to retry a failed job.
- `remoteName`: (Optional) The name of a remote defined in the `remotes` block to run the job on. If not defined job will be run locally

An example would be:

```yaml
jobs:
  localJob:
    schedules: # schedules, multiple schedules can be provided for the same job
      - "0 1 11 * * *" # similar to cron syntax, except the first item (0) is the second
      - "0 1 10 * * *"
    command: whoami
    onlyLogOnFail: false
    onlyNotifyOnFail: false # optional, wether only notify on failure
    notifier: base # optional, notifier of the job. If not provided defaultNotifier will be used.
    failNotifier: more # optional, fail notifier of the job.
    timeoutS: 10 # optional, timeout of running the job
    failStrategy:
      retry:
        maxRetries: 3
```

## Remotes

Remotes is for defining remote servers to execute jobs on. Authentication on the servers is done with `ssh`.

There are two main authentication methods:

- key authentication (using ssh keys)
- password authentication

an example of key authentication would be:

```yaml
remotes:
  my-remote-key:
    username: username
    address: 1.1.1.1
    port: 22
    auth:
      keyAuth:
        keyPath: /path/to/ssh/private.key
        passphrase: ""
```

`passphrase` is only required if the key has a passphrase.

and example of password authentication would be:

```yaml
remotes:
  my-remote-password:
    username: username
    address: 1.1.1.1
    port: 22
    auth:
      passwordAuth:
        password: "12345"
```

## Notifiers

The `notifiers` block is where you configure notifications for services like Telegram and Discord.

Each notifier can have multiple services and the notification would be send to all of them.

### Telegram

- `token`: Your Telegram bot token.
- `receivers`: A list of chat IDs to send notifications to.

an example of telegram notifier would be:

```yaml
notifiers:
  base: # name of the notifier should be used in the job
    telegram:
      token: "TELEGRAM_TOKEN"
      receivers:
        - 1111111111 # chat id to send the message to, for private chats it is equivalent to userId
```

### Discord

- `oAuth2Token` or `botToken`: Your Discord token.
- `channels`: A list of channel IDs to send notifications to.

an example would be:

```yaml
notifiers:
  dis: # name of the notifier should be used in the job
    discord:
      oAuth2Token: "OAUTH2_TOKEN" # you either have to pass oAuth2Token or botToken
      botToken: "botToken"
      channels: # required (at least 1 item), list of channel ids
        - CHANNEL_ID_1
        - CHANNEL_ID_2
```
