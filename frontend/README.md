# Gavande Bros frontend

React application built with Vite, JavaScript, Bootstrap 5, React Router, and
Axios. The current version provides the responsive application shell and
placeholder navigation only; it does not connect to the backend or contain
business data.

## Run locally

```powershell
cd D:\gavande_brothers\frontend
npm install
npm run dev
```

Open the local URL printed by Vite (normally `http://localhost:5173`).

## Structure

- `src/components/` — reusable sidebar and top navigation.
- `src/layouts/` — shared application layout.
- `src/pages/` — overview shell and module placeholders.
- `src/services/` — reserved for API clients and services.
- `src/context/` — reserved for application state and authentication context.
- `src/assets/` — local assets.
- `public/` — static public assets.

Build for production with `npm run build`.
