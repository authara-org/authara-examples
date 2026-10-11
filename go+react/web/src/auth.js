import {
  AutharaApiError,
  AutharaBrowserClient,
  authFetch,
} from "@authara/browser";

export const authara = new AutharaBrowserClient({
  automaticRefresh: { audience: "app" },
});

export function errorMessage(error) {
  return error instanceof Error ? error.message : "Something went wrong.";
}

export function isUnauthorized(error) {
  return error instanceof AutharaApiError && error.status === 401;
}

export function isRecentAuthenticationRequired(error) {
  return (
    error?.code === "recent_authentication_required" &&
    Boolean(error?.authenticationChallenge?.id)
  );
}

export async function withCSRF(action) {
  await authara.getCsrfToken();
  return action();
}

export function login(identifier, password) {
  return withCSRF(() =>
    authara.loginWithPassword({
      audience: "app",
      body: { identifier, password },
    }),
  );
}

export async function signup(email, password) {
  const body = { email, password };
  try {
    await withCSRF(() => authara.signupDirect({ audience: "app", body }));
    return null;
  } catch (error) {
    if (!(error instanceof AutharaApiError) || error.status !== 404) throw error;
  }
  return withCSRF(() =>
    authara.startSignupChallenge({ audience: "app", body }),
  );
}

export function verifySignup(challengeID, code) {
  return withCSRF(() =>
    authara.verifySignupChallenge({
      audience: "app",
      body: { challenge_id: challengeID, code },
    }),
  );
}

export function recoverWithPassword(linkID, password) {
  return withCSRF(() =>
    authara.completeAccountRecoveryLinkWithPassword({
      linkID,
      audience: "app",
      body: { password },
    }),
  );
}

export async function loadProtectedMessage() {
  const response = await authFetch("/api/private");
  if (!response.ok) throw new Error(`Protected API failed (${response.status})`);
  return response.json();
}
