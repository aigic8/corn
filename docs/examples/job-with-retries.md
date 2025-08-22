# Job with Retries

This example shows how to configure a job to retry on failure.

```yaml
disableNotifications: true
jobs:
  retryJob:
    schedules:
      - "0 0 2 * * *" # Run at 2:00 AM every day
    command: "this_command_will_fail"
    failStrategy:
      retry:
        maxRetries: 3
```
