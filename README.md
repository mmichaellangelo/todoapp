# ToDo
> [!WARNING]
> This is a work in progress. It most certainly is not fully functional.

This is a simple but completely overengineered todo app.

The purpose of this app is to help me learn how to build a full-stack application, so I decided to stick to something simple that I can finish in a reasonable amount of time.

The frontend web application is built with SvelteKit and the backend api is written in Go.

Postgres serves as the database with Adminer for database admin things.

The whole stack is dockerized and can be run with docker compose.
It is recommended to run with:
> docker compose up --no-attach proxy

to avoid a massive amount of logs from nginx.

After starting, the webapp will be live at dev.localhost.

Adminer can be accessed at adminer.localhost and the api is proxied at api.localhost.