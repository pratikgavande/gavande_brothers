import { useState } from 'react'
import { Navigate, Outlet } from 'react-router-dom'
import Sidebar from '../components/Sidebar.jsx'
import TopNavbar from '../components/TopNavbar.jsx'
import { useAuth } from '../context/useAuth.js'
import '../App.css'

function AppLayout() {
  const [sidebarOpen, setSidebarOpen] = useState(false)
  const { isAuthenticated } = useAuth()

  if (!isAuthenticated) {
    return <Navigate to="/login" replace />
  }

  function closeSidebar() {
    setSidebarOpen(false)
  }

  return (
    <div className="app-shell d-flex">
      <Sidebar isOpen={sidebarOpen} onNavigate={closeSidebar} />
      {sidebarOpen && (
        <button
          type="button"
          className="sidebar-backdrop d-md-none"
          aria-label="Close navigation menu"
          onClick={closeSidebar}
        />
      )}
      <div className="workspace flex-grow-1">
        <TopNavbar onMenuClick={() => setSidebarOpen(true)} />
        <main className="main-content">
          <Outlet />
        </main>
      </div>
    </div>
  )
}

export default AppLayout
