# Authara examples

Two small Go applications show the two supported browser-integration styles.

- **Go + HTMX:** Authara owns the hosted login, signup, and account UI. The app only protects its own pages with `authara-go` and renders a CSRF-safe logout form.
- **Go + React:** React owns the authentication and account UI. It uses the generated `@authara/browser` client, `@authara/react` Google, Apple, and Microsoft buttons, and a Go API protected by `authara-go`.

The stack pins Authara Core `v0.24.0`, migrations `v0.3.0`, Gateway `v1.2.2`, Go SDK `v0.20.0`, Browser SDK `v0.11.0`, and React components `v0.2.0`.

## Run an example

```bash
git clone https://github.com/authara-org/authara-examples.git
cd authara-examples
cp .env.example .env
docker compose up --build
```

Open <http://localhost:3000>.

`go+htmx` is selected by default. To run the custom React UI, set this in `.env` and rebuild:

```dotenv
AUTHARA_EXAMPLE=go+react
```

```bash
docker compose up --build
```

The examples share one Postgres volume, so accounts remain available when you switch apps. Use `docker compose down --volumes` only when you intentionally want to delete that local data.

## Go + HTMX

This is the smallest server-rendered integration:

1. Links navigate to `/auth/login`, `/auth/signup`, and `/auth/account`.
2. Authara Gateway sends `/auth/*` to Core and everything else to the app.
3. `authara-go` protects `/private` and refreshes an expired app session once.
4. The page uses the SDK's logout form data and one protected HTMX fragment.

There is no custom authentication form or account-management logic in this app.

## Go + React

The React app never navigates to Authara's hosted login or account pages. It includes:

- password login, direct signup, and email-code signup;
- Google, Apple, and Microsoft login through `@authara/react`;
- symmetric provider-collision recovery with password, Google, Apple, or Microsoft proof;
- a custom account screen for username, password, provider links, and sessions;
- password, Google, Apple, and Microsoft reauthentication for sensitive changes;
- API logout and automatic refresh-once behavior;
- a call to `/api/private`, where the Go SDK verifies the Authara access token.

Authara still owns the cookies, CSRF enforcement, session issuance, provider proof validation, and authorization decisions.

For frontend development, start the Compose stack with `AUTHARA_EXAMPLE=go+react`, then run:

```bash
cd go+react/web
npm ci
npm run dev
```

Vite proxies `/auth` and `/api` to the gateway on port `3000`.

## Google, Apple, and Microsoft

Provider configuration belongs to Authara Core, not the React app. Add the desired providers to `.env` using the documented variables in `.env.example`.

For Google, register every origin that displays the button, such as `http://localhost:3000` and, when using Vite, `http://localhost:5173`.

Apple requires a real HTTPS domain. Register this exact return URL for the Services ID:

```text
https://your-domain/auth/oauth/apple/callback
```

Set `PUBLIC_URL=https://your-domain` and route that domain to the gateway.

For Microsoft, register this SPA redirect URI in Microsoft Entra:

```text
https://your-domain/auth/oauth/microsoft/redirect
```

For local development, also register `http://localhost:3000/auth/oauth/microsoft/redirect` and, when using Vite, `http://localhost:5173/auth/oauth/microsoft/redirect`.

## Validate locally

```bash
(cd go+htmx && go test ./...)
(cd go+react && go test ./...)
(cd go+react/web && npm ci && npm run build)
docker compose config --quiet
```

## License

MIT
