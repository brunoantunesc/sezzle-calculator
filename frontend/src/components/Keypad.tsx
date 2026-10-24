import { useEffect, useState } from 'react'
import { Button, CircularProgress } from '@mui/material'

// Each pending request mounts its own indicator and cancellable delay.
function LoadingIndicator() {
  const [visible, setVisible] = useState(false)
  useEffect(() => {
    const timer = setTimeout(() => setVisible(true), 200)
    return () => clearTimeout(timer)
  }, [])
  return visible ? <CircularProgress className="equals-spinner" size={20} color="inherit" aria-label="Calculating" /> : null
}

const keys = [
  { value: 'clear', label: 'AC', name: 'All clear', kind: 'action' },
  { value: 'sign', label: '±', name: 'Toggle sign', kind: 'action' },
  { value: 'sqrt', label: '√', name: 'Square root', kind: 'action' },
  { value: 'divide', label: '÷', name: 'Divide', kind: 'operator' },
  { value: '7', label: '7' }, { value: '8', label: '8' }, { value: '9', label: '9' },
  { value: 'multiply', label: '×', name: 'Multiply', kind: 'operator' },
  { value: '4', label: '4' }, { value: '5', label: '5' }, { value: '6', label: '6' },
  { value: 'subtract', label: '−', name: 'Subtract', kind: 'operator' },
  { value: '1', label: '1' }, { value: '2', label: '2' }, { value: '3', label: '3' },
  { value: 'add', label: '+', name: 'Add', kind: 'operator' },
  { value: 'power', label: 'xʸ', name: 'Power', kind: 'action' },
  { value: '0', label: '0' }, { value: '.', label: '.', name: 'Decimal point' },
  { value: 'backspace', label: '⌫', name: 'Backspace', kind: 'action' },
  { value: '=', label: '=', name: 'Equals', kind: 'equals' },
]

export function Keypad({ press, pressed, operation, pending }: {
  press: (key: string) => void; pressed: Set<string>; operation: string | null; pending: boolean
}) {
  return <div className="keypad" aria-label="Calculator keypad">
    {keys.map(key => <Button key={key.value} type="button" disableRipple
      aria-label={key.name ?? key.label}
      aria-pressed={key.kind === 'operator' || key.value === 'power' ? operation === key.value : undefined}
      disabled={pending && key.value === '='}
      aria-busy={key.value === '=' ? pending : undefined}
      className={`calc-key ${key.kind ?? 'number'} ${pressed.has(key.value) ? 'is-pressed' : ''} ${operation === key.value ? 'is-selected' : ''}`}
      onClick={() => press(key.value)}>{key.label}{key.value === '=' && pending && <LoadingIndicator />}</Button>)}
  </div>
}
