import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import AppLayout from './layouts/AppLayout.jsx'
import DashboardPage from './pages/DashboardPage.jsx'
import ModulePlaceholder from './pages/ModulePlaceholder.jsx'

const modules = [
  { path: 'authentication', title: 'Login & authentication' },
  { path: 'customers', title: 'Customers' },
  { path: 'bills', title: 'Bills' },
  { path: 'bill-items', title: 'Bill items' },
  { path: 'accounting', title: 'Accounting' },
  { path: 'users-roles', title: 'Users & roles' },
]

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<AppLayout />}>
          <Route index element={<DashboardPage />} />
          {modules.map((module) => (
            <Route
              key={module.path}
              path={module.path}
              element={<ModulePlaceholder title={module.title} />}
            />
          ))}
          <Route path="*" element={<Navigate to="/" replace />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export default App
