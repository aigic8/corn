# Job with Notification

This example of a job with notifications in telegram:

```yaml
disableNotifications: false
defaultNotifier: my-telegram-notifier
defaultFailNotifier: my-telegram-notifier
notifyTimeoutMs: 5000
defaultTimeoutS: 5
notifiers:
  my-telegram-notifier:
    telegram:
      token: "YOUR_TELEGRAM_BOT_TOKEN"
      receivers:
        - 123456789 # Your Telegram Chat ID

jobs:
  notifiedJob:
    schedules:
      - "0 30 9 * * *" # Run at 9:30 AM every day
    command: "echo 'This job will send a notification.'"
    notifier: my-telegram-notifier
```
