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
  - [x] Add remote servers SSH support for cron jobs
  - [x] Systemd file for running the app in background and starting on power on
  - [x] Add Installation script
  - [x] Add Failure strategies (from [jobber](https://github.com/dshearer/jobber))
    - `retry:` retry either instantly or with timing between them
    - `halt:` remove the job from schedules until it is fixed and notify the user
    - `continue:` DEFAULT option, continue writing commands like before
- **Later**
  - [ ] Setting command executer (zsh/fish/bash)
  - [ ] Logging level support from config
  - [ ] Notification in email
  - [ ] Report (Daily/Weekly/Monthly) on stats (how many jobs has failed/passed)
  - [ ] CPU/Memory info on each job
  - [ ] Dashboard
  - [ ] Different Job intervals from [gocron](https://github.com/go-co-op/gocron)
  - [ ] Add landing page
  - [ ] Auto-Log removal after some time (with config option to set the time interval)
  - [ ] Add an option to reload the config file
  - [ ] Add support for remote sqlite files

### Enhancements

- [ ] Choose and add License
- [ ] Add Docker compose and volume for sqlite db
- [ ] Test multiline commands using SSH
- [ ] Add Testing
- [ ] Add logging to terminal (instead of writing to file) for testing jobs (with a flag) and in config
  - Pretty logging should also be an option for this case
- [ ] Check why [jobber](https://github.com/dshearer/jobber) uses ipc for long running jobs and what is it
- [ ] Add support for multi-step docker image
- [ ] Add docs and clean the github page
- [ ] Test Discord Notification system
- [x] Fix the bug with quoted args in function `SeperateArgsFromCommand`
- [x] Notify the user on failure before command running
- [ ] Add Emojis to notifications
- [ ] In config get all the paths relative to the config file
- [ ] Add an CLI option to validate the config file
- [ ] Add comments to config options
- [x] Test retry count
- [x] Remove Windows builds from Github releases
