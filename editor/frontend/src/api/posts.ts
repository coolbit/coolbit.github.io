import axios from 'axios'
import type { Post, PostPayload, PagedResponse } from '../types/post'

const http = axios.create({ baseURL: '/api' })

export interface PostsQuery {
  page?: number
  page_size?: number
  category?: string
  q?: string
}

export const listPosts = (query?: PostsQuery) =>
  http.get<PagedResponse<Post>>('/posts', { params: query })
export const getPost = (id: number) => http.get<Post>(`/posts/${id}`)
export const createPost = (data: PostPayload) => http.post<Post>('/posts', data)
export const updatePost = (id: number, data: PostPayload) => http.put<Post>(`/posts/${id}`, data)
export const deletePost = (id: number) => http.delete(`/posts/${id}`)
