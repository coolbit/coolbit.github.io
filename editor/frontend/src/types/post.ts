export interface Post {
  id: number
  title: string
  summary: string
  content: string
  category: string
  created_at: string
  updated_at: string
}

export type PostPayload = Pick<Post, 'title' | 'summary' | 'content' | 'category'>

export interface PagedResponse<T> {
  data: T[]
  total: number
  page: number
  page_size: number
}
