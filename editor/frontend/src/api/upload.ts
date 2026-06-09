import axios from 'axios'

const http = axios.create({ baseURL: '/api' })

export const uploadImage = (postId: number, file: File) => {
  const form = new FormData()
  form.append('file', file)
  return http.post<{ url: string }>(`/upload/${postId}`, form)
}

export const listMedia = (postId: number) =>
  http.get<{ files: string[] }>(`/media/${postId}`)
