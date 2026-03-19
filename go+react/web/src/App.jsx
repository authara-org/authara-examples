import { useEffect, useState } from "react";
import { getCurrentUser, logoutAndRedirect } from "./auth";

export default function App() {
	const [user, setUser] = useState(undefined);

	useEffect(() => {
		getCurrentUser()
			.then(setUser)
			.catch(() => setUser(null));
	}, []);

	const isPrivate = window.location.pathname === "/private";

	const loginURL = "/auth/login?return_to=/private";

	if (user === undefined) {
		return (
			<div className="page">
				<div className="shell">
					<div className="glass">Loading…</div>
				</div>
			</div>
		);
	}

	if (isPrivate && !user) {
		window.location.href = loginURL;
		return null;
	}

	if (isPrivate) {
		return (
			<div className="page">
				<div className="shell">
					<div className="topbar">
						<a className="btn" href="/">Home</a>
						<div className="topbar-right">
							<a className="btn" href="/auth/account">Account</a>
							<button className="btn" onClick={logoutAndRedirect}>Logout</button>
						</div>
					</div>

					<div className="card">
						<h1>Private Page</h1>
						<p className="muted">This page is protected by Authara.</p>
						<p><strong>Email:</strong> <code>{user.email}</code></p>
						<p><strong>Username:</strong> <code>{user.username}</code></p>
					</div>
				</div>
			</div>
		);
	}

	return (
		<div className="page">
			<div className="topbar shell">
				<div className="brand">Authara Presents</div>
				{user ? (
					<a className="btn btn-primary" href="/private">Private</a>
				) : (
					<a className="btn btn-primary" href={loginURL}>Login</a>
				)}
			</div>

			<main className="hero shell">
				<h1>Authara<br />Presents</h1>
				<p className="muted hero-copy">
					A tiny Go + React example with Authara authentication.
				</p>

				<div className="actions">
					{user ? (
						<a className="btn btn-primary" href="/private">Open private page</a>
					) : (
						<a className="btn btn-primary" href={loginURL}>Login with Authara</a>
					)}
				</div>

				<div className="glass">
					{user ? `Signed in as ${user.username}` : "Quiet infrastructure. Loud ambition."}
				</div>
			</main>
		</div>
	);
}
