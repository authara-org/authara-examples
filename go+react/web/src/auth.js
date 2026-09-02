import { AutharaBrowserClient, authFetch } from "@authara/browser";

const authara = new AutharaBrowserClient();

export async function getCurrentUser() {
	const res = await authFetch("/api/me");

	if (res.status === 401) {
		return null;
	}

	if (!res.ok) {
		throw new Error("failed to load current user");
	}

	return await res.json();
}

export async function logoutAndRedirect() {
	await authara.logout();
	window.location.assign("/");
}
