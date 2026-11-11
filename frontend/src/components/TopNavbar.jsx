import { useLocation } from 'react-router-dom'
import { useAuth } from '../context/useAuth.js'

const pageTitles = {
  '/': 'Overview',
  '/customers': 'Customers',
  '/bills': 'Bills',
  '/bill-items': 'Bill items',
  '/accounting': 'Accounting',
  '/users-roles': 'Users & roles',
}

function TopNavbar({ onMenuClick }) {
  const { pathname } = useLocation()
  const { user, logout } = useAuth()
  const pageTitle = pageTitles[pathname] ?? 'Workspace'

  return (
    <header className="topbar d-flex align-items-center justify-content-between px-3 px-lg-4">
      <div className="d-flex align-items-center gap-3">
        <button
          type="button"
          className="btn btn-light d-md-none border-0"
          aria-label="Open navigation menu"
          onClick={onMenuClick}
        >
          <i className="bi bi-list fs-4" aria-hidden="true" />
        </button>
        <div>
          <div className="topbar-title">{pageTitle}</div>
          <div className="topbar-subtitle">Gavande Bros Billing System</div>
        </div>
      </div>
      <div className="d-flex align-items-center gap-2 gap-sm-3">
        <span className="d-none d-sm-inline text-secondary small">
          {user?.full_name ?? user?.email}
        </span>
        <span className="avatar" aria-label={user?.full_name ?? 'Signed-in user'}>
          {user?.full_name?.trim()?.charAt(0)?.toUpperCase() ?? 'U'}
        </span>
        <button
          type="button"
          className="btn btn-outline-secondary btn-sm"
          onClick={logout}
        >
          <i className="bi bi-box-arrow-right me-sm-1" aria-hidden="true" />
          <span className="d-none d-sm-inline">Sign out</span>
        </button>
      </div>
    </header>
  )
}

export default TopNavbar
