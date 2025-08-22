# Basic Local Job

This the configuration for basic local job with multiple schedules:

```yaml
disableNotifications: true
jobs:
  localJob:
    schedules:
      - "0 1 11 * * *" # Run at 11:01:00 every day
      - "0 1 10 * * *" # Run at 10:01:00 every day
    command: whoami
```
