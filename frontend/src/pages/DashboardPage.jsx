function DashboardPage() {
  return (
    <section>
      <div className="mb-4">
        <div className="page-eyebrow mb-2">Gavande Bros</div>
        <h1 className="page-title mb-0">Welcome to your workspace</h1>
      </div>
      <div className="welcome-panel">
        <div className="page-eyebrow mb-3">Billing management</div>
        <h2 className="mb-3">Your workspace is ready</h2>
        <p className="mb-0">
          Use the navigation to access the application modules. Each section
          will be built incrementally; no billing data or API connections are
          configured yet.
        </p>
      </div>
    </section>
  )
}

export default DashboardPage
