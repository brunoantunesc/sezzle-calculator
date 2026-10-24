import { useEffect, useRef, useState } from 'react'
import { calculate as calculateApi, CalculationError } from '../services/calculatorApi'
import { symbols } from '../types/calculator'
import type { Calculate, Operation } from '../types/calculator'

interface State {
  display: string
  expression: string
  left: string | null
  operation: Operation | null
  fresh: boolean
  result: boolean
  pending: boolean
  error: string
}
const initial: State = { display: '0', expression: '', left: null, operation: null, fresh: false, result: false, pending: false, error: '' }

const binary = new Set(['add', 'subtract', 'multiply', 'divide', 'power'])

export function useCalculator(calculate: Calculate = calculateApi, onSuccess?: (expression: string, result: number) => void) {
  const [state, setState] = useState(initial)
  const current = useRef(initial)
  const version = useRef(0)
  useEffect(() => () => { version.current++ }, [])

  function update(next: State) { current.current = next; setState(next) }

  async function submit(operation: Operation, operands: string[], nextOperation: Operation | null = null) {
    const id = ++version.current
    const before = current.current
    const expression = operation === 'sqrt' ? `√(${operands[0]})` : `${operands[0]} ${symbols[operation]} ${operands[1]}`
    update({ ...before, expression, pending: true, error: '' })
    try {
      const numbers = operands.map(Number)
      if (!numbers.every(Number.isFinite)) throw new Error('Invalid input')
      const result = await calculate({ operation, operands: numbers })
      if (typeof result !== 'number' || !Number.isFinite(result)) throw new Error('Invalid response')
      if (id !== version.current) return
      onSuccess?.(expression, result)
      const display = String(result)
      if (operation === 'sqrt' && before.operation) {
        update({ ...before, display, pending: false, fresh: false, result: false, error: '' }); return
      }
      update({ ...initial, display, result: true, fresh: true,
        left: nextOperation ? display : null, operation: nextOperation,
        expression: nextOperation ? `${display} ${symbols[nextOperation]}` : `${expression} =`,
      })
    } catch (error) {
      if (id !== version.current) return
      update({ ...before, pending: false, error: error instanceof CalculationError ? error.message : 'Unable to calculate. Please try again.' })
    }
  }

  function press(key: string) {
    const s = current.current
    if (s.pending) return
    if (key === 'clear') { version.current++; update({ ...initial }); return }
    if (/^\d$/.test(key) || key === '.') {
      const reset = s.fresh || s.result
      let display = reset ? '0' : s.display
      if (key === '.') { if (!display.includes('.')) display += '.' }
      else if (display === '0' || display === '-0') display = `${display.startsWith('-') ? '-' : ''}${key}`
      else display += key
      // Bound entry length for display usability; arithmetic stays on the backend.
      if (display.replace(/[-.]/g, '').length > 16) return
      update({ ...s, display, fresh: false, result: false, error: '', expression: s.operation ? s.expression : '' })
      return
    }
    if (key === 'sign') {
      const display = s.fresh && s.operation ? '-0' : s.display.startsWith('-') ? s.display.slice(1) : `-${s.display}`
      update({ ...s, display, fresh: false, result: false, error: '', expression: s.operation ? s.expression : '' }); return
    }
    if (key === 'backspace') {
      if (s.fresh) return
      const shortened = s.display.slice(0, -1)
      update({ ...s, display: shortened === '' || shortened === '-' ? '0' : shortened, error: '' }); return
    }
    if (binary.has(key)) {
      const operation = key as Operation
      if (s.operation && s.left !== null && !s.fresh) { void submit(s.operation, [s.left, s.display], operation); return }
      update({ ...s, left: s.display, operation, fresh: true, result: false, expression: `${s.display} ${symbols[operation]}`, error: '' }); return
    }
    if (key === 'sqrt') {
      void submit('sqrt', [s.display]); return
    }
    if (key === '=' && s.operation && s.left !== null && !s.fresh) void submit(s.operation, [s.left, s.display])
  }
  function reuse(result: number) {
    if (current.current.pending || !Number.isFinite(result)) return
    version.current++
    update({ ...initial, display: String(result), result: true, fresh: true })
  }
  return { ...state, press, reuse, dismissError: () => update({ ...current.current, error: '' }) }
}
