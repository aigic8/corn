import { defineConfig } from "vitepress";

// https://vitepress.dev/reference/site-config
export default defineConfig({
  base: "/corn/",
  title: "Corn",
  description: "A Modern Cron Manager",
  themeConfig: {
    // https://vitepress.dev/reference/default-theme-config
    nav: [
      { text: "Home", link: "/" },
      { text: "Getting Started", link: "/guide/getting-started" },
    ],

    sidebar: [
      {
        text: "Guide",
        items: [
          { text: "Getting Started", link: "/guide/getting-started" },
          { text: "Installation", link: "/guide/installation" },
          { text: "Configuration", link: "/guide/configuration" },
        ],
      },
      {
        text: "Examples",
        items: [
          { text: "Examples", link: "/examples" },
          { text: "Basic Local Job", link: "/examples/basic-local-job" },
          {
            text: "Basic Remote Job with Password Authentication",
            link: "/examples/basic-remote-job-password-auth",
          },
          {
            text: "Job with Notification",
            link: "/examples/job-with-notification",
          },
          { text: "Job with Retries", link: "/examples/job-with-retries" },
          {
            text: "Remote Job with Notification",
            link: "/examples/remote-job-with-notification",
          },
        ],
      },
    ],

    socialLinks: [{ icon: "github", link: "https://github.com/aigic8/corn" }],
  },
});
