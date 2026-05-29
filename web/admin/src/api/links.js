import http from './http'

export const listLinks = (params) =>
  http.get('/admin/links', { params })

export const getLink = (id) =>
  http.get(`/admin/links/${id}`)

export const createLink = (body) =>
  http.post('/admin/links', body)

export const updateLink = (id, body) =>
  http.patch(`/admin/links/${id}`, body)

export const deleteLink = (id) =>
  http.delete(`/admin/links/${id}`)

export const getLinkStats = (id, params) =>
  http.get(`/admin/links/${id}/stats`, { params })
