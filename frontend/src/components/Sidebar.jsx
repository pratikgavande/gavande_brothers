import { NavLink } from 'react-router-dom'

const navigation = [
  { label: 'Overview', path: '/', icon: 'bi-grid-1x2' },
  { label: 'Login & authentication', path: '/authentication', icon: 'bi-person-lock' },
  { label: 'Customers', path: '/customers', icon: 'bi-people' },
  { label: 'Bills', path: '/bills', icon: 'bi-receipt' },
  { label: 'Bill items', path: '/bill-items', icon: 'bi-list-check' },
  { label: 'Accounting', path: '/accounting', icon: 'bi-journal-text' },
  { label: 'Users & roles', path: '/users-roles', icon: 'bi-shield-check' },
]

function Sidebar({ isOpen, onNavigate }) {
  return (
    <aside className={`sidebar d-flex flex-column ${isOpen ? 'is-open' : ''}`}>
      <div className="d-flex align-items-center gap-3 px-2">
        <span className="brand-mark" aria-hidden="true">
          <i className="bi bi-compass" />
        </span>
        <div>
          <div className="brand-name">Gavande Bros</div>
          <div className="brand-caption">Billing management</div>
        </div>
      </div>

      <div className="nav-section-label">Workspace</div>
      <nav className="nav flex-column" aria-label="Main navigation">
        {navigation.map((item) => (
          <NavLink
            key={item.path}
            to={item.path}
            end={item.path === '/'}
            className={({ isActive }) => `nav-link${isActive ? ' active' : ''}`}
            onClick={onNavigate}
          >
            <i className={`bi ${item.icon}`} aria-hidden="true" />
            <span>{item.label}</span>
          </NavLink>
        ))}
      </nav>

      <div className="sidebar-footer">
        <i className="bi bi-geo-alt me-2" aria-hidden="true" />
        Land surveying services
      </div>
    </aside>
  )
}

export default Sidebar
