# Corn

Corn is a cron job tool with built in logging and notification system.

## Tasks

### Feature Sets

- **Basics**
  - [ ] Setting command executer (zsh/fish/bash)
  - [ ] Setting env vars and env files
  - [ ] Test a single cron job
  - [ ] Set timezone in config
- **Deployable**
  - [ ] Add Github releases
  - [ ] Docker image
- **Later**
  - [ ] Logging level support from config
  - [ ] Systemd file for running the app in background and starting on power on
  - [ ] Notification in email
  - [ ] Report (Daily/Weekly/Monthly) on stats (how many jobs has failed/passed)
  - [ ] CPU/Memory info on each job
  - [ ] Dashboard
  - [ ] Different Job intervals from [gocron](https://github.com/go-co-op/gocron)
