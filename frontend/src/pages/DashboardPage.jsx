import { useCallback, useEffect, useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '../context/useAuth.js'
import { getDashboardRecords } from '../services/dashboardService.js'
import './DashboardPage.css'

const currencyFormatter = new Intl.NumberFormat('en-IN', {
  style: 'currency',
  currency: 'INR',
  maximumFractionDigits: 0,
})

const dateFormatter = new Intl.DateTimeFormat('en-IN', {
  day: '2-digit',
  month: 'short',
  year: 'numeric',
})

const summaryCards = [
  {
    key: 'customers',
    label: 'Total customers',
    icon: 'bi-people',
    color: 'blue',
  },
  { key: 'bills', label: 'Total bills', icon: 'bi-receipt', color: 'violet' },
  {
    key: 'billed',
    label: 'Total billed amount',
    icon: 'bi-wallet2',
    color: 'gold',
    format: 'currency',
  },
  {
    key: 'received',
    label: 'Total amount received',
    icon: 'bi-cash-stack',
    color: 'green',
    format: 'currency',
  },
  {
    key: 'pending',
    label: 'Total pending amount',
    icon: 'bi-hourglass-split',
    color: 'rose',
    format: 'currency',
  },
]

const quickActions = [
  {
    title: 'Add New Customer',
    description: 'Customer module placeholder',
    path: '/customers',
    icon: 'bi-person-plus',
    color: 'blue',
  },
  {
    title: 'Create New Bill',
    description: 'Bill module placeholder',
    path: '/bills',
    icon: 'bi-file-earmark-plus',
    color: 'gold',
  },
  {
    title: 'View All Bills',
    description: 'Open the bills section',
    path: '/bills',
    icon: 'bi-receipt',
    color: 'violet',
  },
  {
    title: 'View All Customers',
    description: 'Open the customers section',
    path: '/customers',
    icon: 'bi-people',
    color: 'green',
  },
]

function formatCurrency(value) {
  return currencyFormatter.format(value)
}

function formatBillDate(value) {
  const date = new Date(value)
  return Number.isNaN(date.getTime()) ? '—' : dateFormatter.format(date)
}

function formatBillStatus(status) {
  return `Status ${status}`
}

function getDashboardError(error) {
  if (error.response?.status === 401) {
    return 'Your session has expired. Sign in again to continue.'
  }
  if (!error.response) {
    return 'Could not connect to the billing service. Check that the backend is running and try again.'
  }
  return (
    error.response.data?.message ||
    error.message ||
    'Could not load dashboard information. Please try again.'
  )
}

function calculateTotals(customers, bills) {
  const totals = bills.reduce(
    (result, bill) => {
      result.billed += bill.amount
      result.received += bill.amount_received
      result.pending += bill.amount_pending
      return result
    },
    { billed: 0, received: 0, pending: 0 },
  )

  if (
    !Number.isSafeInteger(totals.billed) ||
    !Number.isSafeInteger(totals.received) ||
    !Number.isSafeInteger(totals.pending)
  ) {
    throw new Error('The bill totals exceed the safe range for display.')
  }

  return {
    customers: customers.length,
    bills: bills.length,
    ...totals,
  }
}

function SummaryCard({ card, value, loading }) {
  let displayValue = '—'
  if (!loading) {
    displayValue =
      card.format === 'currency' ? formatCurrency(value) : value.toLocaleString('en-IN')
  }

  return (
    <article className="dashboard-stat-card">
      <div className={`dashboard-stat-icon ${card.color}`}>
        <i className={`bi ${card.icon}`} aria-hidden="true" />
      </div>
      <div className="dashboard-stat-label">{card.label}</div>
      {loading ? (
        <div className="placeholder-glow mt-2" aria-label="Loading statistic">
          <span className="placeholder col-7 dashboard-stat-placeholder" />
        </div>
      ) : (
        <div className="dashboard-stat-value">{displayValue}</div>
      )}
    </article>
  )
}

function DashboardPage() {
  const { user } = useAuth()
  const [dashboard, setDashboard] = useState(null)
  const [isLoading, setIsLoading] = useState(true)
  const [errorMessage, setErrorMessage] = useState('')
  const [reloadKey, setReloadKey] = useState(0)

  const refreshDashboard = useCallback(() => {
    setIsLoading(true)
    setErrorMessage('')
    setReloadKey((key) => key + 1)
  }, [])

  useEffect(() => {
    const controller = new AbortController()
    const { signal } = controller

    getDashboardRecords(signal)
      .then((records) => {
        if (!signal.aborted) {
          const totals = calculateTotals(records.customers, records.bills)
          setDashboard({ ...records, totals })
        }
      })
      .catch((error) => {
        if (!signal.aborted) {
          setErrorMessage(getDashboardError(error))
          setDashboard(null)
        }
      })
      .finally(() => {
        if (!signal.aborted) {
          setIsLoading(false)
        }
      })

    return () => controller.abort()
  }, [reloadKey])

  const customerNames = useMemo(
    () => new Map(dashboard?.customers.map((customer) => [customer.customer_id, customer.name])),
    [dashboard?.customers],
  )

  const recentBills = useMemo(
    () =>
      dashboard?.bills
        .slice()
        .sort((left, right) => right.bill_id - left.bill_id)
        .slice(0, 5) ?? [],
    [dashboard?.bills],
  )

  const greetingName = user?.full_name?.trim().split(/\s+/)[0] || 'Gavande Bros'

  return (
    <section className="dashboard-page">
      <div className="dashboard-heading d-flex flex-wrap justify-content-between align-items-end gap-3 mb-4">
        <div>
          <div className="page-eyebrow mb-2">Gavande Bros</div>
          <h1 className="page-title mb-2">Welcome, {greetingName}</h1>
          <p className="dashboard-subtitle mb-0">
            Here is an overview of your billing business.
          </p>
        </div>
        <button
          className="btn btn-outline-secondary dashboard-refresh"
          disabled={isLoading}
          onClick={refreshDashboard}
          type="button"
        >
          <i className={`bi bi-arrow-clockwise me-2 ${isLoading ? 'dashboard-spinning' : ''}`} aria-hidden="true" />
          Refresh
        </button>
      </div>

      {errorMessage && (
        <div className="alert alert-danger d-flex flex-wrap align-items-center justify-content-between gap-3" role="alert">
          <span>
            <i className="bi bi-exclamation-circle me-2" aria-hidden="true" />
            {errorMessage}
          </span>
          <button
            className="btn btn-sm btn-outline-danger"
            onClick={refreshDashboard}
            type="button"
          >
            Try again
          </button>
        </div>
      )}

      {!errorMessage && (
        <>
          <div className="row g-3 g-xl-4 mb-4">
            {summaryCards.map((card) => (
              <div className="col-12 col-sm-6 col-xl" key={card.key}>
                <SummaryCard
                  card={card}
                  loading={isLoading}
                  value={dashboard?.totals[card.key] ?? 0}
                />
              </div>
            ))}
          </div>

          <div className="dashboard-section-card mb-4">
            <div className="dashboard-section-header d-flex flex-wrap align-items-center justify-content-between gap-3">
              <div>
                <h2 className="dashboard-section-title mb-1">Recent bills</h2>
                <p className="dashboard-section-caption mb-0">
                  Latest bills from your billing records
                </p>
              </div>
              <Link className="dashboard-text-link" to="/bills">
                View all bills <i className="bi bi-arrow-up-right ms-1" aria-hidden="true" />
              </Link>
            </div>

            {isLoading ? (
              <div className="dashboard-loading-state" role="status">
                <span className="spinner-border spinner-border-sm me-2" aria-hidden="true" />
                Loading billing records...
              </div>
            ) : recentBills.length === 0 ? (
              <div className="dashboard-empty-state">
                <span className="dashboard-empty-icon">
                  <i className="bi bi-receipt" aria-hidden="true" />
                </span>
                <h3>No bills found</h3>
                <p>
                  {dashboard?.customers.length === 0
                    ? 'There are no customers or bills in the available records yet.'
                    : 'There are no bill records to show yet.'}
                </p>
                <Link className="btn btn-outline-primary btn-sm" to="/bills">
                  Open bills placeholder
                </Link>
              </div>
            ) : (
              <div className="table-responsive">
                <table className="table dashboard-bills-table align-middle mb-0">
                  <thead>
                    <tr>
                      <th scope="col">Bill number</th>
                      <th scope="col">Customer</th>
                      <th scope="col">Bill date</th>
                      <th className="text-end" scope="col">Amount</th>
                      <th className="text-end" scope="col">Received</th>
                      <th className="text-end" scope="col">Pending</th>
                      <th scope="col">Status</th>
                      <th scope="col"><span className="visually-hidden">Action</span></th>
                    </tr>
                  </thead>
                  <tbody>
                    {recentBills.map((bill) => (
                      <tr key={bill.bill_id}>
                        <td className="dashboard-bill-number">#{bill.bill_no}</td>
                        <td className="dashboard-customer-cell">
                          {customerNames.get(bill.customer_id) ?? 'Customer unavailable'}
                        </td>
                        <td className="text-nowrap">{formatBillDate(bill.date)}</td>
                        <td className="text-end text-nowrap">{formatCurrency(bill.amount)}</td>
                        <td className="text-end text-nowrap">{formatCurrency(bill.amount_received)}</td>
                        <td className="text-end text-nowrap">{formatCurrency(bill.amount_pending)}</td>
                        <td>
                          <span className="dashboard-status">{formatBillStatus(bill.status)}</span>
                        </td>
                        <td className="text-end">
                          <Link
                            aria-label={`Open bills placeholder for bill ${bill.bill_no}`}
                            className="dashboard-view-link"
                            to="/bills"
                            title="Bill detail page is not implemented yet"
                          >
                            <i className="bi bi-arrow-up-right" aria-hidden="true" />
                          </Link>
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}
          </div>

          <div className="dashboard-section-header mb-3">
            <h2 className="dashboard-section-title mb-1">Quick actions</h2>
            <p className="dashboard-section-caption mb-0">
              Shortcuts to the customer and bill sections
            </p>
          </div>
          <div className="row g-3">
            {quickActions.map((action) => (
              <div className="col-12 col-sm-6 col-xl-3" key={action.title}>
                <Link className="dashboard-action-card" to={action.path}>
                  <span className={`dashboard-action-icon ${action.color}`}>
                    <i className={`bi ${action.icon}`} aria-hidden="true" />
                  </span>
                  <span className="dashboard-action-copy">
                    <span className="dashboard-action-title">{action.title}</span>
                    <span className="dashboard-action-description">{action.description}</span>
                  </span>
                  <i className="bi bi-arrow-right dashboard-action-arrow" aria-hidden="true" />
                </Link>
              </div>
            ))}
          </div>
        </>
      )}

      {!isLoading && !errorMessage && dashboard && (
        <p className="dashboard-data-note mt-3 mb-0">
          Summary totals use all records across the API pages. Customer and bill data are read-only here.
        </p>
      )}
    </section>
  )
}

export default DashboardPage
