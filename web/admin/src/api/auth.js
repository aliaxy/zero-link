import http from './http'

export const loginApi = (username, password) =>
  http.post('/admin/login', { username, password })

export const getProfileApi = () =>
  http.get('/admin/profile')
