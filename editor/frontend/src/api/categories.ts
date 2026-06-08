import axios from 'axios'

const http = axios.create({ baseURL: '/api' })

export const listCategories = () => http.get<string[]>('/categories')
