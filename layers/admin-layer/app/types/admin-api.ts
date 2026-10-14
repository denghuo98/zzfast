/**
 * 管理后台 API 类型定义 —— 与 zzfast 后端 zschema/zapi 对齐。
 * 后端 JSON 字段统一为 camelCase。
 */

/** 分页请求 */
export interface AdminPageReq {
  page?: number
  pageSize?: number
  pagination?: boolean
}

/** 分页响应 */
export interface AdminPageRes {
  page: number
  pageSize: number
  pagination: boolean
  pageCount: number
  totalCount: number
}

/** 分页列表响应 */
export interface AdminPageList<T> extends AdminPageRes {
  list: T[]
}

// ==================== 认证 ====================

export interface AdminLoginResponse {
  id: number
  username: string
  token: string
  expires: number
}

export interface AdminCaptchaResponse {
  id: string
  base64: string
}

export interface AdminUserInfo {
  id: number
  username: string
  realName: string
  avatar: string
  sex: number
  email: string
  mobile: string
}

// ==================== 用户（Member） ====================

export interface AdminRoleRef {
  id: number
  name: string
  key: string
  remark: string
  sort: number
  status: number
  createdAt?: string
  updatedAt?: string
}

export interface AdminMemberItem {
  id: number
  deptId: number
  realName: string
  username: string
  avatar: string
  sex: number
  email: string
  mobile: string
  lastActiveAt?: string | null
  remark: string
  status: number
  createdAt?: string
  updatedAt?: string
  roleIds: number[]
  roles: AdminRoleRef[]
}

export interface AdminMemberEditParams {
  id?: number
  username: string
  password?: string
  realName: string
  roleIds: number[]
  avatar?: string
  sex?: number
  email?: string
  mobile?: string
  remark?: string
  status?: number
}

export interface AdminMemberListParams extends AdminPageReq {
  username?: string
  realName?: string
  email?: string
  mobile?: string
  status?: number
}

export interface AdminProfileUpdateParams {
  id: number
  avatar?: string
  realName: string
  email?: string
  mobile?: string
  sex?: number
}

export interface AdminPasswordUpdateParams {
  id: number
  oldPassword: string
  newPassword: string
}

// ==================== 角色（Role） ====================

export interface AdminRoleItem {
  id: number
  name: string
  key: string
  remark: string
  sort: number
  status: number
  createdAt?: string
  updatedAt?: string
  menuIds: number[]
}

export interface AdminRoleEditParams {
  id?: number
  name: string
  key: string
  remark?: string
  sort?: number
  status?: number
  menuIds?: number[]
}

export interface AdminRoleListParams extends AdminPageReq {
  name?: string
  key?: string
  status?: number
}

// ==================== 菜单（Menu） ====================

export interface AdminMenuItem {
  id: number
  pid: number
  title: string
  name: string
  path: string
  icon: string
  type: number
  redirect?: string
  component?: string
  isFrame?: number
  frameSrc?: string
  hidden?: number
  sort: number
  remark?: string
  status: number
  permissions?: string
  createdAt?: string
  updatedAt?: string
  children?: AdminMenuItem[]
}

export interface AdminMenuListParams extends AdminPageReq {
  name?: string
}

// ==================== 登录日志（LoginLog） ====================

export interface AdminLoginLogItem {
  id: number
  reqId: string
  memberId: number
  username: string
  response?: unknown
  loginAt?: string
  loginIp: string
  province: string
  city: string
  userAgent: string
  errMsg: string
  status: number
  os?: string
  browser?: string
  createdAt?: string
  updatedAt?: string
}

export interface AdminLoginLogListParams extends AdminPageReq {
  username?: string
  status?: number
  loginIp?: string
  loginAt?: [string, string]
}
