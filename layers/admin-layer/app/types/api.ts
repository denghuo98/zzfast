/** 分页相关通用类型 */
export interface PaginatedResponse<T> {
  list: T[]
  total: number
  page: number
  itemsPerPage: number
}

export interface PaginationParams {
  current?: number
  page?: number
  size?: number
  pageSize?: number
  [key: string]: unknown
}
