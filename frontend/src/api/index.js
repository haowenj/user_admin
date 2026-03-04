import axios from 'axios'

const api = axios.create({
  baseURL: '/api',
  timeout: 10000
})

export const userAPI = {
  register: (data) => api.post('/user/register', data),
  login: (data) => api.post('/user/login', data),
  changePassword: (data) => api.post('/user/change-password', data)
}

export const employeeAPI = {
  getList: () => api.get('/employee'),
  get: (id) => api.get(`/employee/${id}`),
  create: (data) => api.post('/employee', data),
  update: (id, data) => api.put(`/employee/${id}`, data),
  delete: (id) => api.delete(`/employee/${id}`)
}

export default api
