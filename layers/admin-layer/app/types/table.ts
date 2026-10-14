export interface TablePagination {
  /** 当前页码 */
  page: number
  /** 每页条数 */
  pageSize: number
  /** 总条数 */
  total: number
}

export interface ColumnOption<T = any> {
  /** 列类型 */
  type?: 'selection' | 'expand' | 'index' | 'globalIndex'
  /** 列属性名 */
  accessorKey?: string
  /** 列标题 */
  header?: string
  /** 列宽度 */
  width?: string | number
  /** 最小列宽度 */
  minWidth?: string | number
  /** 固定列 */
  fixed?: boolean | 'left' | 'right'
  /** 是否可排序 */
  sortable?: boolean
  /** 过滤器选项 */
  filters?: any[]
  /** 过滤方法 */
  filterMethod?: (value: any, row: any) => boolean
  /** 过滤器位置 */
  filterPlacement?: string
  /** 是否禁用 */
  disabled?: boolean
  /** 是否显示列 */
  visible?: boolean
  /** 是否选中显示 */
  checked?: boolean
  /** 自定义渲染函数 */
  formatter?: (row: T) => any
  /** 是否使用插槽渲染内容 */
  useSlot?: boolean
  /** 插槽名称（默认为 prop 值） */
  slotName?: string
  /** 是否使用表头插槽 */
  useHeaderSlot?: boolean
  /** 表头插槽名称（默认为 `${prop}-header`） */
  headerSlotName?: string
  /** 固定宽度列是否启用超长省略号 + UTooltip 悬浮全文（未设置时自动：无父插槽则自动启用） */
  ellipsis?: boolean
  [key: string]: any
}
