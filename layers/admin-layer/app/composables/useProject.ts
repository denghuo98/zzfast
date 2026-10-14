/**
 * useProject — 选中的项目状态管理
 *
 * 保存到 localStorage，页面刷新后保持选中。
 * 上层应用可在请求头中携带 selectedProject.code。
 */
const STORAGE_KEY = 'admin-selected-project'

export interface ProjectInfo {
  code: string
  name: string
  costCode?: string
  costName?: string
}

export const useProject = () => {
  const selectedProject = ref<ProjectInfo | null>(null)

  // 初始化：从 localStorage 恢复
  if (import.meta.client) {
    try {
      const stored = localStorage.getItem(STORAGE_KEY)
      if (stored) selectedProject.value = JSON.parse(stored)
    } catch { /* ignore */ }
  }

  // 选择项目
  function selectProject(project: ProjectInfo | null) {
    selectedProject.value = project
    if (import.meta.client) {
      if (project) {
        localStorage.setItem(STORAGE_KEY, JSON.stringify(project))
      } else {
        localStorage.removeItem(STORAGE_KEY)
      }
    }
  }

  return {
    selectedProject,
    selectProject,
  }
}
