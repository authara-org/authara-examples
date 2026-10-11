import { useCallback, useEffect, useState } from "react";
import { AppleButton, GoogleButton } from "@authara/react";
import {
  authara,
  errorMessage,
  isRecentAuthenticationRequired,
  isUnauthorized,
  loadProtectedMessage,
  login,
  recoverWithPassword,
  signup,
  verifySignup,
  withCSRF,
} from "./auth";

function Feedback({ error, notice }) {
  return (
    <div className="feedback" aria-live="polite">
      {error && <p className="error">{error}</p>}
      {notice && <p className="notice">{notice}</p>}
    </div>
  );
}

function ProviderButtons({ label = "sign-in", onAuthenticated, onError }) {
  return (
    <div className="provider-buttons">
      <GoogleButton
        action="login"
        audience="app"
        label={label}
        mode="system"
        onSuccess={onAuthenticated}
        onError={onError}
      />
      <AppleButton
        action="login"
        audience="app"
        label={label}
        mode="system"
        onSuccess={onAuthenticated}
        onError={onError}
      />
    </div>
  );
}

function Recovery({ recovery, onAuthenticated, onCancel }) {
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [busy, setBusy] = useState(false);
  const methods = new Set(recovery.proof_methods);

  async function usePassword(event) {
    event.preventDefault();
    setBusy(true);
    setError("");
    try {
      await recoverWithPassword(recovery.link_id, password);
      await onAuthenticated();
    } catch (requestError) {
      setError(errorMessage(requestError));
    } finally {
      setBusy(false);
    }
  }

  return (
    <section className="card auth-card">
      <p className="eyebrow">Account found</p>
      <h1>Confirm the existing account.</h1>
      <p>Use one of its existing sign-in methods. Authara will then safely link the new provider.</p>

      {methods.has("password") && (
        <form onSubmit={usePassword}>
          <label>
            Existing password
            <input
              type="password"
              autoComplete="current-password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              disabled={busy}
              required
            />
          </label>
          <button type="submit" disabled={busy}>Continue with password</button>
        </form>
      )}

      <div className="provider-buttons">
        {methods.has("google") && (
          <GoogleButton
            action="recover"
            audience="app"
            recoveryLinkId={recovery.link_id}
            mode="system"
            onSuccess={onAuthenticated}
            onError={(providerError) => setError(errorMessage(providerError))}
          />
        )}
        {methods.has("apple") && (
          <AppleButton
            action="recover"
            audience="app"
            recoveryLinkId={recovery.link_id}
            mode="system"
            onSuccess={onAuthenticated}
            onError={(providerError) => setError(errorMessage(providerError))}
          />
        )}
      </div>

      <Feedback error={error} />
      <button className="text-button" type="button" onClick={onCancel}>Start over</button>
    </section>
  );
}

function AuthScreen({ onAuthenticated }) {
  const [mode, setMode] = useState("login");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [challengeID, setChallengeID] = useState("");
  const [code, setCode] = useState("");
  const [recovery, setRecovery] = useState(null);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");

  if (recovery) {
    return (
      <main className="centered">
        <Recovery
          recovery={recovery}
          onAuthenticated={onAuthenticated}
          onCancel={() => setRecovery(null)}
        />
      </main>
    );
  }

  function providerError(providerError) {
    if (providerError.code === "account_link_required" && providerError.recovery) {
      setRecovery(providerError.recovery);
      setError("");
      return;
    }
    setError(errorMessage(providerError));
  }

  async function submit(event) {
    event.preventDefault();
    setBusy(true);
    setError("");
    setNotice("");
    try {
      if (mode === "verify") {
        await verifySignup(challengeID, code);
        await onAuthenticated();
        return;
      }
      if (mode === "login") {
        await login(email, password);
        await onAuthenticated();
        return;
      }

      const challenge = await signup(email, password);
      if (!challenge) {
        await onAuthenticated();
        return;
      }
      setChallengeID(challenge.challenge_id);
      setPassword("");
      setMode("verify");
      setNotice("Check your email for the verification code.");
    } catch (requestError) {
      setError(errorMessage(requestError));
    } finally {
      setBusy(false);
    }
  }

  function selectMode(nextMode) {
    setMode(nextMode);
    setPassword("");
    setCode("");
    setError("");
    setNotice("");
  }

  return (
    <main className="centered">
      <section className="card auth-card">
        <p className="eyebrow">Custom React UI</p>
        <h1>{mode === "verify" ? "Verify your email." : mode === "signup" ? "Create an account." : "Welcome back."}</h1>
        <p>All authentication calls go directly to Authara on this origin.</p>

        <form onSubmit={submit}>
          {mode === "verify" ? (
            <label>
              Verification code
              <input
                value={code}
                onChange={(event) => setCode(event.target.value)}
                inputMode="numeric"
                autoComplete="one-time-code"
                required
              />
            </label>
          ) : (
            <>
              <label>
                Email
                <input
                  type="email"
                  value={email}
                  onChange={(event) => setEmail(event.target.value)}
                  autoComplete="email"
                  required
                />
              </label>
              <label>
                Password
                <input
                  type="password"
                  value={password}
                  onChange={(event) => setPassword(event.target.value)}
                  autoComplete={mode === "signup" ? "new-password" : "current-password"}
                  required
                />
              </label>
            </>
          )}
          <button type="submit" disabled={busy}>
            {busy ? "Working…" : mode === "verify" ? "Verify" : mode === "signup" ? "Create account" : "Sign in"}
          </button>
        </form>

        {mode !== "verify" && (
          <>
            <div className="divider"><span>or</span></div>
            <ProviderButtons
              label={mode === "signup" ? "sign-up" : "sign-in"}
              onAuthenticated={onAuthenticated}
              onError={providerError}
            />
          </>
        )}

        <Feedback error={error} notice={notice} />
        {mode !== "verify" && (
          <button className="text-button" type="button" onClick={() => selectMode(mode === "login" ? "signup" : "login")}>
            {mode === "login" ? "Create an account" : "I already have an account"}
          </button>
        )}
      </section>
    </main>
  );
}

function Header({ onLogout }) {
  return (
    <header className="topbar shell">
      <a className="brand" href="/">Authara React</a>
      <nav>
        <a href="/">Home</a>
        <a href="/account">Account</a>
        <button className="text-button" type="button" onClick={onLogout}>Sign out</button>
      </nav>
    </header>
  );
}

function Dashboard({ user, onLogout }) {
  const [result, setResult] = useState(null);
  const [error, setError] = useState("");

  async function callBackend() {
    setError("");
    try {
      setResult(await loadProtectedMessage());
    } catch (requestError) {
      setError(errorMessage(requestError));
    }
  }

  return (
    <>
      <Header onLogout={onLogout} />
      <main className="shell stack">
        <section className="card">
          <p className="eyebrow">Signed in</p>
          <h1>Hello, {user.username}.</h1>
          <p>{user.email}</p>
        </section>
        <section className="card">
          <h2>Protected Go API</h2>
          <p>The browser helper refreshes once if necessary; the Go SDK verifies the app access token.</p>
          <button type="button" onClick={callBackend}>Call /api/private</button>
          {result && <pre>{JSON.stringify(result, null, 2)}</pre>}
          <Feedback error={error} />
        </section>
      </main>
    </>
  );
}

function Reauthentication({ account, challengeID, onComplete, onCancel, onError }) {
  const [password, setPassword] = useState("");
  const methods = new Set(account.auth_methods.map((method) => method.provider));

  async function usePassword(event) {
    event.preventDefault();
    try {
      await withCSRF(() =>
        authara.reauthenticateWithPassword({
          body: { authentication_challenge_id: challengeID, password },
        }),
      );
      await onComplete();
    } catch (requestError) {
      onError(requestError);
    }
  }

  return (
    <section className="card reauthentication">
      <h2>Confirm it is you</h2>
      <p>Authara requires a recent sign-in before this account change.</p>
      {methods.has("password") && (
        <form onSubmit={usePassword}>
          <label>
            Password
            <input type="password" value={password} onChange={(event) => setPassword(event.target.value)} required />
          </label>
          <button type="submit">Continue with password</button>
        </form>
      )}
      <div className="provider-buttons">
        {methods.has("google") && (
          <GoogleButton
            action="reauthenticate"
            authenticationChallengeId={challengeID}
            mode="system"
            onSuccess={onComplete}
            onError={onError}
          />
        )}
        {methods.has("apple") && (
          <AppleButton
            action="reauthenticate"
            authenticationChallengeId={challengeID}
            mode="system"
            onSuccess={onComplete}
            onError={onError}
          />
        )}
      </div>
      <button className="text-button" type="button" onClick={onCancel}>Cancel</button>
    </section>
  );
}

function AccountPage({ user, onUserChanged, onLogout }) {
  const [account, setAccount] = useState(null);
  const [username, setUsername] = useState(user.username);
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [error, setError] = useState("");
  const [notice, setNotice] = useState("");
  const [reauthentication, setReauthentication] = useState(null);

  const loadAccount = useCallback(async () => {
    try {
      setAccount(await authara.getCurrentAccount());
    } catch (requestError) {
      setError(errorMessage(requestError));
    }
  }, []);

  useEffect(() => {
    void loadAccount();
  }, [loadAccount]);

  function handleError(requestError, retry) {
    if (isRecentAuthenticationRequired(requestError)) {
      setReauthentication({
        challengeID: requestError.authenticationChallenge.id,
        retry,
      });
      return;
    }
    setError(errorMessage(requestError));
  }

  async function run(action, successMessage) {
    setError("");
    setNotice("");
    try {
      await action();
      await loadAccount();
      setNotice(successMessage);
    } catch (requestError) {
      handleError(requestError, () => run(action, successMessage));
    }
  }

  async function finishReauthentication() {
    const retry = reauthentication?.retry;
    setReauthentication(null);
    if (retry) await retry();
    else setNotice("Identity confirmed. Retry the account action.");
  }

  if (!account) {
    return (
      <>
        <Header onLogout={onLogout} />
        <main className="shell"><section className="card">Loading account…<Feedback error={error} /></section></main>
      </>
    );
  }

  const providers = new Set(account.auth_methods.map((method) => method.provider));

  return (
    <>
      <Header onLogout={onLogout} />
      <main className="shell stack">
        <section className="card">
          <p className="eyebrow">React-owned account page</p>
          <h1>Account</h1>
          <p>{account.user.email}</p>
        </section>

        {reauthentication && (
          <Reauthentication
            account={account}
            challengeID={reauthentication.challengeID}
            onComplete={finishReauthentication}
            onCancel={() => setReauthentication(null)}
            onError={(requestError) => setError(errorMessage(requestError))}
          />
        )}

        <section className="card">
          <h2>Profile</h2>
          <form
            onSubmit={(event) => {
              event.preventDefault();
              void run(
                () => withCSRF(() => authara.changeCurrentUsername({ body: { username } })),
                "Username updated.",
              ).then(onUserChanged);
            }}
          >
            <label>
              Username
              <input value={username} onChange={(event) => setUsername(event.target.value)} required />
            </label>
            <button type="submit">Save username</button>
          </form>
        </section>

        <section className="card">
          <h2>Sign-in methods</h2>
          <ul className="method-list">
            {account.auth_methods.map((method) => (
              <li key={method.provider}>
                <span>{method.provider}</span>
                <button
                  className="text-button"
                  type="button"
                  onClick={() => run(
                    () => withCSRF(() => authara.unlinkCurrentUserAuthMethod({ provider: method.provider })),
                    `${method.provider} removed.`,
                  )}
                >
                  Remove
                </button>
              </li>
            ))}
          </ul>

          {(!providers.has("google") || !providers.has("apple")) && (
            <div className="provider-buttons">
              {!providers.has("google") && (
                <GoogleButton
                  action="link"
                  mode="system"
                  onSuccess={() => loadAccount().then(() => setNotice("Google linked."))}
                  onError={(requestError) => handleError(requestError)}
                />
              )}
              {!providers.has("apple") && (
                <AppleButton
                  action="link"
                  mode="system"
                  onSuccess={() => loadAccount().then(() => setNotice("Apple linked."))}
                  onError={(requestError) => handleError(requestError)}
                />
              )}
            </div>
          )}

          <form
            onSubmit={(event) => {
              event.preventDefault();
              const action = providers.has("password")
                ? () => authara.changeCurrentUserPassword({ body: { current_password: currentPassword, new_password: newPassword } })
                : () => authara.addCurrentUserPassword({ body: { password: newPassword } });
              void run(() => withCSRF(action), providers.has("password") ? "Password changed." : "Password added.");
            }}
          >
            {providers.has("password") && (
              <label>
                Current password
                <input type="password" value={currentPassword} onChange={(event) => setCurrentPassword(event.target.value)} autoComplete="current-password" required />
              </label>
            )}
            <label>
              New password
              <input type="password" value={newPassword} onChange={(event) => setNewPassword(event.target.value)} autoComplete="new-password" required />
            </label>
            <button type="submit">{providers.has("password") ? "Change password" : "Add password"}</button>
          </form>
        </section>

        <section className="card">
          <h2>Sessions</h2>
          <ul className="method-list">
            {account.sessions.map((session) => (
              <li key={session.id}>
                <span>{session.current ? "This session" : session.user_agent || "Session"}</span>
                {!session.current && (
                  <button
                    className="text-button"
                    type="button"
                    onClick={() => run(
                      () => withCSRF(() => authara.revokeCurrentUserSession({ sessionID: session.id })),
                      "Session revoked.",
                    )}
                  >
                    Revoke
                  </button>
                )}
              </li>
            ))}
          </ul>
        </section>

        <Feedback error={error} notice={notice} />
      </main>
    </>
  );
}

export default function App() {
  const [user, setUser] = useState(undefined);
  const [error, setError] = useState("");

  const loadUser = useCallback(async () => {
    try {
      setUser(await authara.getCurrentUser());
      setError("");
    } catch (requestError) {
      if (isUnauthorized(requestError)) {
        setUser(null);
        return;
      }
      setError(errorMessage(requestError));
      setUser(null);
    }
  }, []);

  useEffect(() => {
    void loadUser();
  }, [loadUser]);

  async function logout() {
    try {
      await withCSRF(() => authara.logout());
      setUser(null);
      window.location.assign("/");
    } catch (requestError) {
      setError(errorMessage(requestError));
    }
  }

  if (user === undefined) return <main className="centered"><section className="card">Loading…</section></main>;
  if (!user) return <><AuthScreen onAuthenticated={loadUser} /><Feedback error={error} /></>;
  if (window.location.pathname === "/account") {
    return <AccountPage user={user} onUserChanged={loadUser} onLogout={logout} />;
  }
  return <Dashboard user={user} onLogout={logout} />;
}
