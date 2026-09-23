const dateFormat = new Intl.DateTimeFormat('pt-BR')

export function formatDate(iso: string | undefined): string {
  return iso ? dateFormat.format(new Date(iso)) : '—'
}

const MINUTE = 60_000
const HOUR = 60 * MINUTE
const DAY = 24 * HOUR

export function formatDuration(fromIso: string, toIso: string): string {
  const elapsed = new Date(toIso).getTime() - new Date(fromIso).getTime()
  const days = Math.floor(elapsed / DAY)
  const hours = Math.floor((elapsed % DAY) / HOUR)
  const minutes = Math.floor((elapsed % HOUR) / MINUTE)

  if (days > 0) {
    const dayLabel = days === 1 ? '1 dia' : `${days} dias`
    return hours > 0 ? `${dayLabel} e ${hours} h` : dayLabel
  }
  if (hours > 0) {
    return minutes > 0 ? `${hours} h ${minutes} min` : `${hours} h`
  }
  return minutes > 0 ? `${minutes} min` : 'menos de 1 min'
}

const dateTimeFormat = new Intl.DateTimeFormat('pt-BR', { dateStyle: 'short', timeStyle: 'short' })

export function formatDateTime(iso: string | undefined): string {
  return iso ? dateTimeFormat.format(new Date(iso)) : '—'
}

const relative = new Intl.RelativeTimeFormat('pt-BR', { numeric: 'auto' })

/** "há 5 minutos", "ontem"... for recent moments, the date otherwise. */
export function formatRelative(iso: string | undefined, now: Date = new Date()): string {
  if (!iso) {
    return '—'
  }
  const elapsed = new Date(iso).getTime() - now.getTime()
  const minutes = Math.round(elapsed / MINUTE)
  if (minutes === 0) {
    return 'agora mesmo'
  }
  if (Math.abs(minutes) < 60) {
    return relative.format(minutes, 'minute')
  }
  const hours = Math.round(elapsed / HOUR)
  if (Math.abs(hours) < 24) {
    return relative.format(hours, 'hour')
  }
  const days = Math.round(elapsed / DAY)
  if (Math.abs(days) < 7) {
    return relative.format(days, 'day')
  }
  return formatDate(iso)
}
