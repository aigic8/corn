# Remote Job with Password Authentication

This example shows how to configure a remote job using password authentication.

```yaml
disableNotifications: true
remotes:
  my-remote-password:
    username: myuser
    address: 192.168.1.100
    port: 22
    auth:
      passwordAuth:
        password: "your_password"

jobs:
  remoteJob:
    schedules:
      - "0 0 * * * *" # Run at the beginning of every hour
    command: "ls -la /home/myuser"
    remoteName: my-remote-password
```
