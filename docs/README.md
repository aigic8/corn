# Corn

Corn is a cron job tool with built in logging and notification system.

## Tasks

### Feature Sets

- **Basics**
  - [x] Test a single cron job
  - [x] Set timezone in config
- **Deployable**
  - [ ] Add Github releases
  - [ ] Docker image & CI/CD pipeline
- **Later**
  - [ ] Setting command executer (zsh/fish/bash)
  - [ ] Setting env vars and env files
  - [ ] Logging level support from config
  - [ ] Systemd file for running the app in background and starting on power on
  - [ ] Notification in email
  - [ ] Report (Daily/Weekly/Monthly) on stats (how many jobs has failed/passed)
  - [ ] CPU/Memory info on each job
  - [ ] Dashboard
  - [ ] Different Job intervals from [gocron](https://github.com/go-co-op/gocron)
  - [ ] Add installation landing page & script
  - [ ] Add remote servers SSH support for cron jobs

### Enhancements

- [ ] Add Testing
- [ ] Add logging to terminal (instead of writing to file) for testing jobs (with a flag) and in config
  - Pretty logging should also be an option for this case
