# Corn

Corn is a cron job tool with built in logging and notification system.

## Tasks

### Feature Sets

- **Basics**
  - [x] Test a single cron job
  - [x] Set timezone in config
- **Deployable**
  - [x] Add Github releases
  - [x] Docker image & CI/CD pipeline
- **Final Features Before Beta**
  - [x] Add timeout to the config and use the time out in the settings
  - [ ] Setting env vars and env files
  - [ ] Add remote servers SSH support for cron jobs
  - [ ] Add Installation script
  - [ ] Add Failure strategies (from [jobber](https://github.com/dshearer/jobber))
- **Later**
  - [ ] Setting command executer (zsh/fish/bash)
  - [ ] Logging level support from config
  - [ ] Systemd file for running the app in background and starting on power on
  - [ ] Notification in email
  - [ ] Report (Daily/Weekly/Monthly) on stats (how many jobs has failed/passed)
  - [ ] CPU/Memory info on each job
  - [ ] Dashboard
  - [ ] Different Job intervals from [gocron](https://github.com/go-co-op/gocron)
  - [ ] Add landing page
  - [ ] Auto-Log removal after some time (with config option to set the time interval)

### Enhancements

- [ ] Add Testing
- [ ] Check why [jobber](https://github.com/dshearer/jobber) uses ipc for long running jobs and what is it
- [ ] Add logging to terminal (instead of writing to file) for testing jobs (with a flag) and in config
  - Pretty logging should also be an option for this case
- [ ] Add support for multi-step docker image
- [ ] Add docs and clean the github page
- [ ] Test Discord Notification system
- [ ] Fix the bug with quoted args in function `SeperateArgsFromCommand`
- [ ] Notify the user on failure before command running
- [ ] Add Emojis to notifications
