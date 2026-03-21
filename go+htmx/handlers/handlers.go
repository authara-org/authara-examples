package handlers

import (
	"fmt"
	"html"
	"html/template"
	"net/http"

	"github.com/authara-org/authara-go/authara"
)

type Handler struct {
	autharaBaseURL string
	client         *authara.Client
}

type HomeData struct {
	LoggedIn bool
	Username string
	LoginURL string
}

type PrivateData struct {
	Email      string
	Username   string
	Logout     authara.LogoutFormData
	AccountURL string
}

func New(autharaBaseURL string) *Handler {
	return &Handler{
		autharaBaseURL: autharaBaseURL,
		client:         authara.NewClient(autharaBaseURL),
	}
}

func (h *Handler) Home(w http.ResponseWriter, r *http.Request) {
	user, err := h.client.GetCurrentUser(r.Context(), r)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	data := HomeData{
		LoggedIn: user != nil,
		LoginURL: "/auth/login?return_to=/private",
	}
	if user != nil {
		data.Username = user.Username
	}

	if err := homeTemplate.Execute(w, data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (h *Handler) Private(w http.ResponseWriter, r *http.Request) {
	user, err := h.client.GetCurrentUser(r.Context(), r)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Redirect(w, r, "/auth/login?return_to=/private", http.StatusSeeOther)
		return
	}

	logout, ok := authara.LogoutFormDataFromRequest(
		r,
		"/auth/login?return_to=/private",
	)
	if !ok {
		http.Redirect(w, r, "/auth/login?return_to=/private", http.StatusSeeOther)
		return
	}

	data := PrivateData{
		Email:      user.Email,
		Username:   user.Username,
		Logout:     logout,
		AccountURL: "/auth/account",
	}

	if err := privateTemplate.Execute(w, data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

func (h *Handler) PrivatePulse(w http.ResponseWriter, r *http.Request) {
	user, err := h.client.GetCurrentUser(r.Context(), r)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Redirect(w, r, "/auth/login?return_to=/private", http.StatusSeeOther)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(
		w,
		`<div class="glass">Authenticated as <strong>%s</strong></div>`,
		html.EscapeString(user.Username),
	)
}

var homeTemplate = template.Must(template.New("home").Parse(`
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <title>Authara Presents</title>
  <script src="https://cdn.jsdelivr.net/npm/htmx.org@2.0.8/dist/htmx.min.js" integrity="sha384-/TgkGk7p307TH7EXJDuUlgG3Ce1UVolAOFopFekQkkXihi5u/6OCvVKyz1W+idaz" crossorigin="anonymous"></script>
  <style>
    :root {
      --bg: #0b0d10;
      --panel: rgba(255,255,255,0.06);
      --line: rgba(255,255,255,0.12);
      --text: #f7f8fa;
      --muted: #a7adb7;
      --accent: #8b5cf6;
      --accent2: #22c55e;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      font-family: Inter, ui-sans-serif, system-ui, sans-serif;
      background:
        radial-gradient(circle at top left, rgba(139,92,246,0.25), transparent 30%),
        radial-gradient(circle at bottom right, rgba(34,197,94,0.18), transparent 30%),
        var(--bg);
      color: var(--text);
      min-height: 100vh;
    }
    .nav {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 22px 28px;
    }
    .brand {
      font-weight: 700;
      letter-spacing: 0.08em;
      text-transform: uppercase;
      color: var(--muted);
      font-size: 13px;
    }
    .btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      gap: 8px;
      border: 1px solid var(--line);
      background: rgba(255,255,255,0.04);
      color: var(--text);
      text-decoration: none;
      padding: 10px 16px;
      border-radius: 999px;
      transition: 160ms ease;
      cursor: pointer;
    }
    .btn:hover { background: rgba(255,255,255,0.08); }
    .btn-primary {
      background: linear-gradient(135deg, var(--accent), #6d28d9);
      border-color: transparent;
    }
    .btn-primary:hover { filter: brightness(1.08); }
    .hero {
      max-width: 1000px;
      margin: 0 auto;
      padding: 80px 28px 120px;
      text-align: center;
    }
    h1 {
      font-size: clamp(48px, 10vw, 110px);
      line-height: 0.95;
      margin: 0 0 18px;
      letter-spacing: -0.04em;
    }
    .sub {
      font-size: 18px;
      color: var(--muted);
      max-width: 720px;
      margin: 0 auto 30px;
      line-height: 1.6;
    }
    .actions {
      display: flex;
      gap: 14px;
      justify-content: center;
      flex-wrap: wrap;
      margin-bottom: 26px;
    }
    .glass {
      display: inline-block;
      padding: 16px 20px;
      border-radius: 18px;
      background: var(--panel);
      border: 1px solid var(--line);
      color: var(--muted);
      backdrop-filter: blur(10px);
    }
    .meta {
      margin-top: 26px;
      color: var(--muted);
      font-size: 14px;
    }
    code {
      background: rgba(255,255,255,0.06);
      padding: 2px 8px;
      border-radius: 8px;
      color: #ddd6fe;
    }
  </style>
</head>
<body>
  <nav class="nav">
    <div class="brand">Authara Presents</div>
    <div>
      {{if .LoggedIn}}
        <a class="btn" href="/private">Private</a>
      {{else}}
        <a class="btn btn-primary" href="{{.LoginURL}}">Login</a>
      {{end}}
    </div>
  </nav>

  <main class="hero">
    <h1>Authara<br>Presents</h1>
    <p class="sub">
      A tiny Go + HTMX starter with startup energy, protected routes, and Authara-based authentication.
    </p>

    <div class="actions">
      <button
        class="btn"
        hx-get="/private/pulse"
        hx-target="#present-slot"
        hx-swap="innerHTML"
      >
        Feel the vibe
      </button>

      {{if .LoggedIn}}
        <a class="btn btn-primary" href="/private">Open private page</a>
      {{else}}
        <a class="btn btn-primary" href="{{.LoginURL}}">Login with Authara</a>
      {{end}}
    </div>

    <div id="present-slot">
      {{if .LoggedIn}}
        <div class="glass">Signed in as <strong>{{.Username}}</strong></div>
      {{else}}
        <div class="glass">Quiet infrastructure. Loud ambition.</div>
      {{end}}
    </div>

    <div class="meta">
      {{if .LoggedIn}}
        You are already authenticated.
      {{else}}
        Public landing page
      {{end}}
    </div>
  </main>
</body>
</html>
`))

var privateTemplate = template.Must(template.New("private").Parse(`
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width,initial-scale=1">
  <title>Private</title>
  <script src="https://cdn.jsdelivr.net/npm/htmx.org@2.0.8/dist/htmx.min.js" integrity="sha384-/TgkGk7p307TH7EXJDuUlgG3Ce1UVolAOFopFekQkkXihi5u/6OCvVKyz1W+idaz" crossorigin="anonymous"></script>

  <style>
    :root {
      --bg: #0b0d10;
      --panel: rgba(255,255,255,0.06);
      --line: rgba(255,255,255,0.12);
      --text: #f7f8fa;
      --muted: #a7adb7;
      --accent: #8b5cf6;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      font-family: Inter, ui-sans-serif, system-ui, sans-serif;
      background:
        radial-gradient(circle at top left, rgba(139,92,246,0.20), transparent 30%),
        var(--bg);
      color: var(--text);
      min-height: 100vh;
    }
    .wrap {
      max-width: 900px;
      margin: 0 auto;
      padding: 28px;
    }
    .top {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 28px;
    }
    .top-right {
      display: flex;
      gap: 10px;
      align-items: center;
    }
    .btn {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      border: 1px solid var(--line);
      background: rgba(255,255,255,0.04);
      color: var(--text);
      text-decoration: none;
      padding: 10px 16px;
      border-radius: 999px;
      cursor: pointer;
    }
    .card {
      padding: 28px;
      border-radius: 24px;
      background: var(--panel);
      border: 1px solid var(--line);
      backdrop-filter: blur(10px);
    }
    h1 {
      margin: 0 0 12px;
      font-size: 40px;
      letter-spacing: -0.03em;
    }
    p {
      margin: 0 0 20px;
      color: var(--muted);
      line-height: 1.6;
    }
    code {
      background: rgba(255,255,255,0.06);
      padding: 2px 8px;
      border-radius: 8px;
      color: #ddd6fe;
    }
    form { margin: 0; }
    .glass {
      display: inline-block;
      padding: 16px 20px;
      border-radius: 18px;
      background: rgba(255,255,255,0.06);
      border: 1px solid rgba(255,255,255,0.12);
      color: #a7adb7;
      backdrop-filter: blur(10px);
    }
  </style>
</head>
<body>
  <div class="wrap">
    <div class="top">
      <a class="btn" href="/">Home</a>

      <div class="top-right">
        <a class="btn" href="{{.AccountURL}}">Account</a>
        <form method="{{.Logout.Method}}" action="{{.Logout.Action}}">
          <input type="hidden" name="{{.Logout.CSRFName}}" value="{{.Logout.CSRFValue}}">
          <button class="btn" type="submit">Logout</button>
        </form>
      </div>
    </div>

    <div class="card">
      <h1>Private page</h1>
      <p>This page is protected by Authara.</p>
      <p><strong>Email:</strong> <code>{{.Email}}</code></p>
      <p><strong>Username:</strong> <code>{{.Username}}</code></p>

      <button
        class="btn"
        hx-get="/private/pulse"
        hx-target="#pulse"
        hx-swap="innerHTML"
      >
        Ping protected fragment
      </button>

      <div id="pulse" style="margin-top:16px;"></div>
    </div>
  </div>

  <script>
    window.addEventListener("pageshow", (event) => {
      if (event.persisted) {
        window.location.reload();
      }
    });
  </script>
</body>
</html>
`))
