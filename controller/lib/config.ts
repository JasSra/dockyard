export const CONFIG = {
  ADMIN_SECRET: process.env.ADMIN_SECRET || 'admin',
  AGENT_URL: process.env.AGENT_URL || 'http://localhost:8080',
  DOMAIN: process.env.DOMAIN || '{{ DOMAIN }}', // { SPECULATION }
  APPS_ROOT: process.env.APPS_ROOT || '{{ APPS_ROOT }}', // { SPECULATION }
}
