# Remote Job with Notification

This example combines a remote job with notifications.

```yaml
disableNotifications: false
defaultNotifier: my-discord-notifier
defaultFailNotifier: my-discord-notifier
notifyTimeoutMs: 5000
defaultTimeoutS: 5
remotes:
  my-remote-key:
    username: remoteuser
    address: remote.server.com
    port: 22
    auth:
      keyAuth:
        keyPath: /home/user/.ssh/id_rsa

notifiers:
  my-discord-notifier:
    discord:
      botToken: "YOUR_DISCORD_BOT_TOKEN"
      channels:
        - "YOUR_DISCORD_CHANNEL_ID"

jobs:
  remoteNotifiedJob:
    schedules:
      - "0 15 * * * *" # Run every hour at 15 minutes past the hour
    command: "df -h"
    remoteName: my-remote-key
    notifier: my-discord-notifier
    onlyNotifyOnFail: true
```
