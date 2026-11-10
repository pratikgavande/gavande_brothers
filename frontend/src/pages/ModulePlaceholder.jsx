function ModulePlaceholder({ title }) {
  return (
    <section>
      <div className="mb-4">
        <div className="page-eyebrow mb-2">Workspace module</div>
        <h1 className="page-title mb-0">{title}</h1>
      </div>
      <div className="module-placeholder">
        <h2 className="h5 mb-2">{title}</h2>
        <p className="mb-0">
          This section is reserved for future implementation. No business data
          or backend API is connected.
        </p>
      </div>
    </section>
  )
}

export default ModulePlaceholder
