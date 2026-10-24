import { useEffect, useState } from 'react'

export interface HistoryEntry { expression: string; result: number; timestamp: string }
const storageKey = 'calculator-history'
function readHistory(): HistoryEntry[] {
  try {
    const stored: unknown = JSON.parse(localStorage.getItem(storageKey) ?? '[]')
    if (!Array.isArray(stored)) return []
    return stored.filter((item): item is HistoryEntry => item && typeof item === 'object' &&
      typeof item.expression === 'string' && item.expression.length <= 200 &&
      typeof item.result === 'number' && Number.isFinite(item.result) &&
      typeof item.timestamp === 'string' && item.timestamp.length <= 40 && Number.isFinite(Date.parse(item.timestamp)))
      .slice(0, 10)
  } catch { return [] }
}

export function useHistory() {
  const [entries, setEntries] = useState(readHistory)
  useEffect(() => {
    try { localStorage.setItem(storageKey, JSON.stringify(entries)) } catch { /* History remains available for this session. */ }
  }, [entries])
  function add(expression: string, result: number) {
    setEntries(previous => [{ expression, result, timestamp: new Date().toISOString() }, ...previous].slice(0, 10))
  }
  return { entries, add, clear: () => setEntries([]) }
}
