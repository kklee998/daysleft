# Daysleft
A bot for printing how many days are left in the year

# Fly.io deployment

```bash
fly auth login
fly apps create
fly secrets set -a $YOUR_APP_NAME BSKY_USERNAME=
fly secrets set -a $YOUR_APP_NAME BSKY_APP_PASSWORD=
fly secrets set -a $YOUR_APP_NAME SENTRY_DSN=

fly machines run . --rm --app $YOUR_APP_NAME --region $YOUR_REGION --schedule daily
```