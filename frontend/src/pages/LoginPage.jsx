import { useState } from 'react'
import { Navigate, useLocation, useNavigate } from 'react-router-dom'
import { useAuth } from '../context/useAuth.js'
import './LoginPage.css'

function getLoginError(error) {
  if (error.response?.status === 401) {
    return 'The email or password is incorrect. Please try again.'
  }
  if (error.response?.status === 400) {
    return 'Enter a valid email address and password.'
  }
  if (!error.response) {
    return 'Could not reach the server. Check that the backend is running and try again.'
  }
  if (error.message?.includes('invalid login response') || error.message?.includes('expired access token')) {
    return 'The server response could not be used to start a session. Please contact support.'
  }
  return error.response?.data?.error ?? 'Sign in failed. Please try again.'
}

function LoginPage() {
  const { isAuthenticated, login } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [showPassword, setShowPassword] = useState(false)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [errorMessage, setErrorMessage] = useState('')

  if (isAuthenticated) {
    return <Navigate to="/" replace />
  }

  async function handleSubmit(event) {
    event.preventDefault()
    setErrorMessage('')

    const form = event.currentTarget
    if (!form.checkValidity()) {
      form.reportValidity()
      return
    }

    setIsSubmitting(true)
    try {
      await login(email.trim(), password)
      const destination = location.state?.from?.pathname ?? '/'
      navigate(destination, { replace: true })
    } catch (error) {
      setErrorMessage(getLoginError(error))
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <main className="login-page">
      <section className="login-brand-panel">
        <div className="login-brand">
          <span className="login-brand-mark" aria-hidden="true">
            <i className="bi bi-compass" />
          </span>
          <div>
            <div className="login-brand-name">Gavande Bros</div>
            <div className="login-brand-caption">Land surveying services</div>
          </div>
        </div>
        <div className="login-brand-message">
          <div className="login-overline">Billing management</div>
          <h1>Welcome back.</h1>
          <p>Sign in to continue to your business workspace.</p>
        </div>
        <div className="login-brand-footer">
          <i className="bi bi-shield-check me-2" aria-hidden="true" />
          Secure access for your team
        </div>
      </section>

      <section className="login-form-panel">
        <div className="login-form-wrap">
          <div className="login-mobile-brand">
            <span className="login-brand-mark" aria-hidden="true">
              <i className="bi bi-compass" />
            </span>
            <span className="login-brand-name">Gavande Bros</span>
          </div>
          <div className="login-form-heading">
            <div className="page-eyebrow mb-2">Your workspace</div>
            <h2>Sign in to your account</h2>
            <p>Enter the email and password associated with your account.</p>
          </div>

          <form onSubmit={handleSubmit} noValidate>
            {errorMessage && (
              <div className="alert alert-danger d-flex gap-2" role="alert">
                <i className="bi bi-exclamation-circle mt-1" aria-hidden="true" />
                <span>{errorMessage}</span>
              </div>
            )}

            <div className="mb-3">
              <label className="form-label" htmlFor="login-email">
                Email address
              </label>
              <div className="input-group login-input-group">
                <span className="input-group-text">
                  <i className="bi bi-envelope" aria-hidden="true" />
                </span>
                <input
                  autoComplete="username"
                  className="form-control"
                  id="login-email"
                  name="email"
                  onChange={(event) => setEmail(event.target.value)}
                  placeholder="name@company.com"
                  required
                  type="email"
                  value={email}
                />
              </div>
            </div>

            <div className="mb-4">
              <label className="form-label" htmlFor="login-password">
                Password
              </label>
              <div className="input-group login-input-group">
                <span className="input-group-text">
                  <i className="bi bi-lock" aria-hidden="true" />
                </span>
                <input
                  autoComplete="current-password"
                  className="form-control"
                  id="login-password"
                  name="password"
                  onChange={(event) => setPassword(event.target.value)}
                  placeholder="Enter your password"
                  required
                  type={showPassword ? 'text' : 'password'}
                  value={password}
                />
                <button
                  aria-label={showPassword ? 'Hide password' : 'Show password'}
                  aria-pressed={showPassword}
                  className="btn btn-outline-secondary login-password-toggle"
                  onClick={() => setShowPassword((visible) => !visible)}
                  type="button"
                >
                  <i className={`bi ${showPassword ? 'bi-eye-slash' : 'bi-eye'}`} aria-hidden="true" />
                </button>
              </div>
            </div>

            <button className="btn login-submit w-100" disabled={isSubmitting} type="submit">
              {isSubmitting ? (
                <>
                  <span className="spinner-border spinner-border-sm me-2" aria-hidden="true" />
                  Signing in...
                </>
              ) : (
                <>
                  Sign in
                  <i className="bi bi-arrow-right ms-2" aria-hidden="true" />
                </>
              )}
            </button>
          </form>

          <p className="login-form-footnote">
            <i className="bi bi-lock-fill me-1" aria-hidden="true" />
            Your session is protected with a short-lived access token.
          </p>
        </div>
      </section>
    </main>
  )
}

export default LoginPage
