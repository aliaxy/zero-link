import axios from 'axios'
import { ElMessage } from 'element-plus'

const http = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? '/api',
  timeout: 10000,
})

http.interceptors.request.use((config) => {
  const token = localStorage.getItem('zl_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

http.interceptors.response.use(
  (response) => response.data.data ?? response.data,
  (error) => {
    const status = error.response?.status
    const code = error.response?.data?.code
    const message = error.response?.data?.message ?? error.message

    if (status === 401 || code === 'UNAUTHENTICATED') {
      localStorage.removeItem('zl_token')
      localStorage.removeItem('zl_admin')
      window.location.href = '/login'
      return Promise.reject(error)
    }

    ElMessage.error(message)
    return Promise.reject(Object.assign(error, { code, message }))
  },
)

export default http
