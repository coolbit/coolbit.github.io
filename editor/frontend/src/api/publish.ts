import axios from 'axios'

const http = axios.create({ baseURL: '/api' })

export interface PublishResult {
  dir: string
  count: number
}

export const publish = () => http.post<PublishResult>('/publish')
