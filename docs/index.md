---
# https://vitepress.dev/reference/default-theme-home-page
layout: home

hero:
  name: Corn
  tagline: A Modern Cronjob Manager
  actions:
    - theme: brand
      text: Getting Started
      link: /guide/getting-started
    - theme: alt
      text: Configuration Examples
      link: /examples

features:
  - title: Structured Logging
    details: Execution logs of both tasks and agent itself are stored in files in json format.
  - title: Agentless Remote Execution
    details: Jobs can be be run on remote servers without installing corn on them. (through ssh)
  - title: Notification
    details: Notifications are support on Telegram and Discord
---
