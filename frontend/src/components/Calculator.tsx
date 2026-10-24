import { useEffect, useRef, useState } from 'react'
import { IconButton, Paper, SvgIcon } from '@mui/material'
import { useCalculator } from '../hooks/useCalculator'
import type { Calculate } from '../types/calculator'
import { Display } from './Display'
import { Keypad } from './Keypad'
import { History } from './History'
import { ErrorNotice } from './ErrorNotice'
import { useHistory } from '../hooks/useHistory'

const keyboard: Record<string, string> = {
  '+': 'add', '-': 'subtract', '*': 'multiply', 'x': 'multiply', 'X': 'multiply', '/': 'divide',
  '^': 'power', 'r': 'sqrt', 'R': 'sqrt', '=': '=', Enter: '=', Backspace: 'backspace', Escape: 'clear', Delete: 'clear', '.': '.',
}

export function Calculator({ calculate }: { calculate?: Calculate }) {
  const history = useHistory()
  const calculator = useCalculator(calculate, history.add)
  const [historyOpen, setHistoryOpen] = useState(false)
  const press = useRef(calculator.press)
  useEffect(() => { press.current = calculator.press })
  const [pressed, setPressed] = useState<Set<string>>(new Set())
  useEffect(() => {
    const held = new Map<string, string>()
    function down(event: KeyboardEvent) {
      if (historyOpen || event.defaultPrevented) return
      if (event.ctrlKey || event.metaKey || event.altKey || event.target instanceof HTMLElement && event.target.closest('input, textarea, select, [contenteditable="true"]')) return
      // Preserve native Enter/Space activation for focused buttons.
      if ((event.key === 'Enter' || event.key === ' ') && event.target instanceof HTMLElement && event.target.closest('button')) return
      const key = /^\d$/.test(event.key) ? event.key : keyboard[event.key]
      if (!key) return
      event.preventDefault()
      if (event.repeat) return
      held.set(event.code || event.key, key)
      setPressed(new Set(held.values()))
      press.current(key)
    }
    function up(event: KeyboardEvent) {
      held.delete(event.code || event.key)
      setPressed(new Set(held.values()))
    }
    function reset() { held.clear(); setPressed(new Set()) }
    window.addEventListener('keydown', down)
    window.addEventListener('keyup', up)
    window.addEventListener('blur', reset)
    return () => { window.removeEventListener('keydown', down); window.removeEventListener('keyup', up); window.removeEventListener('blur', reset) }
  }, [historyOpen])
  return <Paper component="section" elevation={0} className="calculator" aria-label="Calculator">
    <div className="calculator-toolbar">
      <IconButton className="history-button" aria-label="History" title="History"
        aria-expanded={historyOpen} aria-controls={historyOpen ? 'calculation-history' : undefined}
        onClick={() => { setPressed(new Set()); setHistoryOpen(true) }}>
        <SvgIcon><path d="M13 3a9 9 0 0 0-8.5 6H2l3.5 4L9 9H6.6A7 7 0 1 1 6 15H3.9A9 9 0 1 0 13 3Zm-1 4v6l5 3 1-1.7-4-2.3V7Z" /></SvgIcon>
      </IconButton>
    </div>
    <Display {...calculator} />
    <History open={historyOpen} onClose={() => setHistoryOpen(false)} entries={history.entries}
      clear={history.clear} reuse={calculator.reuse} pending={calculator.pending} />
    <ErrorNotice message={calculator.error} dismiss={calculator.dismissError} />
    <Keypad press={calculator.press} pressed={pressed} operation={calculator.operation} pending={calculator.pending} />
  </Paper>
}
