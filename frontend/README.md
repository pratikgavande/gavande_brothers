# Gavande Bros frontend

React application built with Vite, JavaScript, Bootstrap 5, React Router, and
Axios. The current version includes backend login and a responsive application
shell with placeholder module navigation.

## Run locally

```powershell
cd D:\gavande_brothers\frontend
npm install
npm run dev
```

Open the local URL printed by Vite (normally `http://localhost:5173`).
The Vite development server proxies `/api` requests to `http://localhost:8080`.
Set `VITE_API_BASE_URL` in `.env.local` when the backend API is hosted elsewhere.

The access token is kept in memory and sent as a Bearer token by the shared
Axios client. Reloading the page clears the session and requires signing in
again. Do not store access tokens in local storage.

## Structure

- `src/components/` — reusable sidebar and top navigation.
- `src/layouts/` — shared application layout.
- `src/pages/` — overview shell and module placeholders.
- `src/services/` — reserved for API clients and services.
- `src/context/` — reserved for application state and authentication context.
- `src/assets/` — local assets.
- `public/` — static public assets.

Build for production with `npm run build`.
