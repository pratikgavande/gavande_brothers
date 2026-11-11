import { BrowserRouter, Navigate, Route, Routes, useLocation } from 'react-router-dom'
import { AuthProvider } from './context/AuthContext.jsx'
import { useAuth } from './context/useAuth.js'
import AppLayout from './layouts/AppLayout.jsx'
import DashboardPage from './pages/DashboardPage.jsx'
import LoginPage from './pages/LoginPage.jsx'
import ModulePlaceholder from './pages/ModulePlaceholder.jsx'

const modules = [
  { path: 'customers', title: 'Customers' },
  { path: 'bills', title: 'Bills' },
  { path: 'bill-items', title: 'Bill items' },
  { path: 'accounting', title: 'Accounting' },
  { path: 'users-roles', title: 'Users & roles' },
]

function ProtectedRoute({ children }) {
  const { isAuthenticated } = useAuth()
  const location = useLocation()

  if (!isAuthenticated) {
    return <Navigate to="/login" replace state={{ from: location }} />
  }

  return children
}

function ApplicationRoutes() {
  const { isAuthenticated } = useAuth()

  return (
    <Routes>
      <Route
        path="/login"
        element={isAuthenticated ? <Navigate to="/" replace /> : <LoginPage />}
      />
      <Route
        element={
          <ProtectedRoute>
            <AppLayout />
          </ProtectedRoute>
        }
      >
        <Route index element={<DashboardPage />} />
        {modules.map((module) => (
          <Route
            key={module.path}
            path={module.path}
            element={<ModulePlaceholder title={module.title} />}
          />
        ))}
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <ApplicationRoutes />
      </AuthProvider>
    </BrowserRouter>
  )
}

export default App
