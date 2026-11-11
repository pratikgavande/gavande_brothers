import apiClient from './apiClient.js'

export async function login(email, password) {
  const { data } = await apiClient.post('/auth/login', { email, password })

  if (
    typeof data?.access_token !== 'string' ||
    !data.access_token ||
    typeof data?.token_type !== 'string' ||
    data.token_type.toLowerCase() !== 'bearer' ||
    typeof data?.expires_at !== 'string' ||
    !Number.isFinite(Date.parse(data.expires_at)) ||
    !data?.user
  ) {
    throw new Error('The server returned an invalid login response.')
  }

  return data
}
