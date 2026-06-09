export interface Post {
  id: number
  title: string
  summary: string
  cover_image: string
  content: string
  created_at: string
  updated_at: string
}

export type PostPayload = Pick<Post, 'title' | 'summary' | 'content' | 'cover_image'>

export interface PagedResponse<T> {
  data: T[]
  total: number
  page: number
  page_size: number
}
