/**
 * 表格全局配置
 *
 * 提供表格与后端接口的分页字段映射，支持多种常见后端响应格式。
 */
export const tableConfig = {
  /** 响应数据 — 列表数据字段名（按优先级查找） */
  recordFields: ['list', 'data', 'records', 'items', 'result', 'rows'],
  /** 响应数据 — 总条数字段名 */
  totalFields: ['total', 'count', 'totalCount'],
  /** 响应数据 — 当前页码字段名 */
  currentFields: ['current', 'page', 'pageNum'],
  /** 响应数据 — 每页大小字段名 */
  sizeFields: ['size', 'pageSize', 'limit', 'itemsPerPage'],

  /** 请求参数 — 分页参数名 */
  paginationKey: {
    current: 'page',
    size: 'pageSize',
  },
}
