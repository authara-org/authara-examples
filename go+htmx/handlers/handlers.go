package handlers

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/authara-org/authara-go/authara"
)

type Handler struct {
	client *authara.Client
}

type privateData struct {
	Email    string
	Username string
	Logout   authara.LogoutFormData
}

func New(autharaBaseURL string) *Handler {
	return &Handler{client: authara.NewClient(autharaBaseURL)}
}

func (h *Handler) Home(w http.ResponseWriter, _ *http.Request) {
	render(w, homeTemplate, nil)
}

func (h *Handler) Private(w http.ResponseWriter, r *http.Request) {
	user, err := h.client.CallGetCurrentUser(r.Context(), r)
	if err != nil {
		http.Error(w, "could not load the Authara user", http.StatusBadGateway)
		return
	}

	logout, ok := authara.LogoutFormDataFromRequest(r, "/")
	if !ok {
		http.Redirect(w, r, "/auth/login?return_to=/private", http.StatusSeeOther)
		return
	}

	render(w, privateTemplate, privateData{
		Email:    user.Email,
		Username: user.Username,
		Logout:   logout,
	})
}

func (h *Handler) CurrentUser(w http.ResponseWriter, r *http.Request) {
	user, err := h.client.CallGetCurrentUser(r.Context(), r)
	if err != nil {
		http.Error(w, "could not load the Authara user", http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = userTemplate.Execute(w, user)
}

func render(w http.ResponseWriter, tmpl *template.Template, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "template error", http.StatusInternalServerError)
	}
}

const pageStyle = `
  :root { color-scheme: light dark; font-family: system-ui, sans-serif; }
  body { max-width: 48rem; margin: 4rem auto; padding: 0 1rem; line-height: 1.5; }
  nav, .actions { display: flex; gap: .75rem; flex-wrap: wrap; align-items: center; }
  nav { justify-content: space-between; margin-bottom: 3rem; }
  a, button { color: inherit; }
  .button { border: 1px solid currentColor; border-radius: .6rem; padding: .6rem .9rem; background: transparent; text-decoration: none; cursor: pointer; }
  .card { border: 1px solid color-mix(in srgb, currentColor 20%, transparent); border-radius: 1rem; padding: 1.5rem; }
  form { margin: 0; }
`

var homeTemplate = template.Must(template.New("home").Parse(fmt.Sprintf(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Authara Go + HTMX</title>
  <style>%s</style>
</head>
<body>
  <nav><strong>Go + HTMX</strong><a href="/private">Private page</a></nav>
  <main>
    <h1>Use Authara's native UI.</h1>
    <p>This application only owns its product pages. Authara owns login, signup, account management, and session security.</p>
    <div class="actions">
      <a class="button" href="/auth/login?return_to=/private">Sign in</a>
      <a class="button" href="/auth/signup?return_to=/private">Create account</a>
    </div>
  </main>
</body>
</html>`, pageStyle)))

var privateTemplate = template.Must(template.New("private").Parse(fmt.Sprintf(`<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Private page</title>
  <script src="https://cdn.jsdelivr.net/npm/htmx.org@2.0.8/dist/htmx.min.js" integrity="sha384-/TgkGk7p307TH7EXJDuUlgG3Ce1UVolAOFopFekQkkXihi5u/6OCvVKyz1W+idaz" crossorigin="anonymous"></script>
  <style>%s</style>
</head>
<body>
  <nav>
    <a href="/">Home</a>
    <div class="actions">
      <a class="button" href="/auth/account">Manage account</a>
      <form method="{{.Logout.Method}}" action="{{.Logout.Action}}">
        <input type="hidden" name="{{.Logout.CSRFName}}" value="{{.Logout.CSRFValue}}">
        <button class="button" type="submit">Sign out</button>
      </form>
    </div>
  </nav>
  <main class="card">
    <h1>Private page</h1>
    <div id="current-user">
      <p><strong>{{.Username}}</strong><br>{{.Email}}</p>
    </div>
    <button class="button" hx-get="/private/current-user" hx-target="#current-user" type="button">Reload with HTMX</button>
  </main>
</body>
</html>`, pageStyle)))

var userTemplate = template.Must(template.New("user").Parse(`<p><strong>{{.Username}}</strong><br>{{.Email}}</p>`))
