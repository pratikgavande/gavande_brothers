import { useLocation } from 'react-router-dom'

const pageTitles = {
  '/': 'Overview',
  '/authentication': 'Login & authentication',
  '/customers': 'Customers',
  '/bills': 'Bills',
  '/bill-items': 'Bill items',
  '/accounting': 'Accounting',
  '/users-roles': 'Users & roles',
}

function TopNavbar({ onMenuClick }) {
  const { pathname } = useLocation()
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
      <div className="d-flex align-items-center gap-3">
        <span className="d-none d-sm-inline text-secondary small">Welcome</span>
        <span className="avatar" aria-label="User profile placeholder">
          GB
        </span>
      </div>
    </header>
  )
}

export default TopNavbar
