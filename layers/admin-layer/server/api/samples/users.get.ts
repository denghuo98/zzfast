/**
 * GET /api/samples/users — 为高级表格示例生成随机用户数据
 *
 * 查询参数：
 *   page    — 页码 (默认 1)
 *   pageSize — 每页条数 (默认 10)
 *   keyword  — 搜索关键词 (可选)
 */
export default defineEventHandler((event) => {
  const query = getQuery(event)
  const page = Math.max(1, parseInt(String(query.page)) || 1)
  const pageSize = Math.min(100, Math.max(1, parseInt(String(query.pageSize)) || 10))
  const keyword = String(query.keyword || '').toLowerCase()

  const firstNames = ['Alice', 'Bob', 'Charlie', 'Diana', 'Eve', 'Frank', 'Grace', 'Henry', 'Iris', 'Jack', 'Kate', 'Leo', 'Mia', 'Noah', 'Olivia', 'Peter', 'Quinn', 'Rose', 'Sam', 'Tina']
  const lastNames = ['Johnson', 'Smith', 'Brown', 'Prince', 'Williams', 'Miller', 'Lee', 'Davis', 'Wilson', 'Moore']
  const roles = ['Admin', 'Editor', 'User', 'Viewer', 'Manager']
  const statuses = ['Active', 'Inactive', 'Pending']

  // 随机种子，保证同一页的数据稳定
  let seed = page * 1000

  function pseudoRandom(): number {
    seed = (seed * 16807 + 0) % 2147483647
    return (seed - 1) / 2147483646
  }

  function pick<T>(arr: T[]): T { return arr[Math.floor(pseudoRandom() * arr.length)] }
  function randDate(): string {
    const d = new Date(2025, Math.floor(pseudoRandom() * 12), Math.floor(pseudoRandom() * 28) + 1)
    return d.toISOString().slice(0, 10)
  }

  const total = 856
  const items = []

  for (let i = 0; i < pageSize; i++) {
    const idx = (page - 1) * pageSize + i
    if (idx >= total) break

    const firstName = pick(firstNames)
    const lastName = pick(lastNames)
    const name = `${firstName} ${lastName}`
    const email = `${firstName.toLowerCase()}.${lastName.toLowerCase()}@example.com`

    const item = {
      id: idx + 1,
      name,
      email,
      role: pick(roles),
      status: pick(statuses),
      lastLogin: randDate(),
    }

    if (!keyword) {
      items.push(item)
    } else if (
      item.name.toLowerCase().includes(keyword) ||
      item.email.toLowerCase().includes(keyword) ||
      item.role.toLowerCase().includes(keyword)
    ) {
      items.push(item)
    }
  }

  return {
    list: items,
    total: keyword ? items.length : total,
    page,
    itemsPerPage: pageSize,
  }
})
