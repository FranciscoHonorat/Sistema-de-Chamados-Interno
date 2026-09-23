import { formatDuration, formatRelative } from './format'

describe('formatDuration', () => {
  const start = '2026-09-22T10:00:00Z'

  it.each([
    ['2026-09-22T10:00:30Z', 'menos de 1 min'],
    ['2026-09-22T10:45:00Z', '45 min'],
    ['2026-09-22T12:15:00Z', '2 h 15 min'],
    ['2026-09-22T13:00:00Z', '3 h'],
    ['2026-09-25T14:30:00Z', '3 dias e 4 h'],
    ['2026-09-23T10:00:00Z', '1 dia'],
  ])('from 10:00 to %s is %s', (end, expected) => {
    expect(formatDuration(start, end)).toBe(expected)
  })
})

describe('formatRelative', () => {
  const now = new Date('2026-09-23T12:00:00Z')

  it.each([
    ['2026-09-23T11:59:50Z', 'agora mesmo'],
    ['2026-09-23T11:55:00Z', 'há 5 minutos'],
    ['2026-09-23T09:00:00Z', 'há 3 horas'],
    ['2026-09-22T12:00:00Z', 'ontem'],
    ['2026-09-01T12:00:00Z', '01/09/2026'],
    [undefined, '—'],
  ])('%s → %s', (iso, expected) => {
    expect(formatRelative(iso, now)).toBe(expected)
  })
})
