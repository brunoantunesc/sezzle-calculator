
export function Display({ display, expression, pending }: {
  display: string; expression: string; pending: boolean
}) {
  return <div className="display" aria-label="Calculator display" aria-busy={pending}>
    <div className="display-top" aria-hidden="true" />
    <div className="expression">{expression || '\u00a0'}</div>
    <output className={`display-value ${display.length > 11 ? 'long-value' : ''}`} aria-live="polite" aria-atomic="true">{display}</output>
    <div className="display-message" aria-hidden="true">{'\u00a0'}</div>
  </div>
}
