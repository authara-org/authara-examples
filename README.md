# Authara Starter

Minimal starter project demonstrating how to integrate an application with **Authara**.

Includes a small example application built with **Go + HTMX** running behind the **Authara Gateway**.

## Stack

- Go
- HTMX
- Chi router
- Docker
- PostgreSQL
- Authara

## Quickstart

Clone the repository:

```bash
git clone https://github.com/authara-org/authara-starter
cd authara-starter
```

Create the environment file:

```bash
cp .env.example .env
```

Start the stack:

```bash
docker compose up --build
```

Open the app:

```
http://localhost:3000
```

## What this example shows

- Public homepage
- Login with Authara
- Protected `/private` page
- Logout flow
- Link to `/auth/account`

The example app uses the **Authara Go SDK** to fetch the current user.

## Project structure

```
authara-starter
│
├─ docker-compose.yml
├─ .env.example
│
├─ go+htmx
│   ├─ main.go
│   ├─ handlers/
│   └─ Dockerfile
```

## Learn more

Authara documentation  
https://docs.authara.org

Authara repository  
https://github.com/authara-org/authara

Authara Go SDK  
https://github.com/authara-org/authara-go

## License

MIT
