/**
 * 格式化日期字符串为中文格式
 * @param dateString - 日期字符串，支持格式: "2026-04-02 22:55:45" 或 "2026-04-02T22:55:45"
 */
export function formatDate(dateString: string | undefined | null): string {
  if (!dateString) return '-'

  try {
    let date: Date

    if (dateString.includes(' ')) {
      date = new Date(dateString.replace(' ', 'T'))
    } else {
      date = new Date(dateString)
    }

    if (isNaN(date.getTime())) {
      return '-'
    }

    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
    })
  } catch {
    return '-'
  }
}

/**
 * 格式化日期时间为中文格式
 * @param dateString - 日期字符串
 */
export function formatDateTime(dateString: string | undefined | null): string {
  if (!dateString) return '-'

  try {
    let date: Date

    if (dateString.includes(' ')) {
      date = new Date(dateString.replace(' ', 'T'))
    } else {
      date = new Date(dateString)
    }

    if (isNaN(date.getTime())) {
      return '-'
    }

    return date.toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
    })
  } catch {
    return '-'
  }
}

/**
 * 格式化相对时间（如：刚刚、5分钟前、1天前）
 * @param dateString - 日期字符串
 */
export function formatRelativeTime(dateString: string | undefined | null): string {
  if (!dateString) return '-'

  try {
    let date: Date

    if (dateString.includes(' ')) {
      date = new Date(dateString.replace(' ', 'T'))
    } else {
      date = new Date(dateString)
    }

    if (isNaN(date.getTime())) {
      return '-'
    }

    const now = new Date()
    const diff = now.getTime() - date.getTime()

    const seconds = Math.floor(diff / 1000)
    const minutes = Math.floor(seconds / 60)
    const hours = Math.floor(minutes / 60)
    const days = Math.floor(hours / 24)

    if (seconds < 60) {
      return '刚刚'
    } else if (minutes < 60) {
      return `${minutes} 分钟前`
    } else if (hours < 24) {
      return `${hours} 小时前`
    } else if (days < 7) {
      return `${days} 天前`
    } else {
      return formatDate(dateString)
    }
  } catch {
    return '-'
  }
}
