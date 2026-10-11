import React from "react";
import ReactDOM from "react-dom/client";
import { AutharaProvider } from "@authara/react";
import App from "./App";
import { authara } from "./auth";
import "./styles.css";

ReactDOM.createRoot(document.getElementById("root")).render(
	<React.StrictMode>
		<AutharaProvider client={authara}>
			<App />
		</AutharaProvider>
	</React.StrictMode>
);
