import { Button, Drawer, IconButton, SvgIcon, useMediaQuery } from '@mui/material'
import type { HistoryEntry } from '../hooks/useHistory'

export function History({ open, onClose, entries, clear, reuse, pending }: {
  open: boolean; onClose: () => void; entries: HistoryEntry[]; clear: () => void;
  reuse: (result: number) => void; pending: boolean
}) {
  const desktop = useMediaQuery('(min-width:600px)')
  return <Drawer anchor={desktop ? 'right' : 'bottom'} open={open} onClose={onClose}
    slotProps={{ paper: {
      role: 'dialog', 'aria-modal': true, 'aria-label': 'Calculation history', id: 'calculation-history',
      sx: desktop
        ? { width: 'min(380px, calc(100vw - 32px))', top: 16, right: 16, bottom: 16, height: 'auto', maxHeight: 'calc(100dvh - 32px)', borderRadius: '20px' }
        : { maxHeight: '70dvh', borderRadius: '24px 24px 0 0' },
    } }}>
    <div className="history-heading"><h2>History</h2>
      <IconButton aria-label="Close history" onClick={onClose}><SvgIcon><path d="m6.4 5 5.6 5.6L17.6 5 19 6.4 13.4 12l5.6 5.6-1.4 1.4-5.6-5.6L6.4 19 5 17.6l5.6-5.6L5 6.4Z" /></SvgIcon></IconButton>
    </div>
    <div className="history-content">
      {entries.length === 0 ? <p>No calculations yet.</p> : <ol className="history-list">
        {entries.map((entry, index) => <li key={`${entry.timestamp}-${index}`}>
          <Button className="history-entry" disabled={pending} onClick={() => { reuse(entry.result); onClose() }}
            aria-label={`Reuse result ${entry.result}`}>
            <span className="history-expression">{entry.expression}</span>
            <strong>{String(entry.result)}</strong>
            <time dateTime={entry.timestamp}>{new Date(entry.timestamp).toLocaleString()}</time>
          </Button>
        </li>)}
      </ol>}
    </div>
    {entries.length > 0 && <div className="history-actions"><Button onClick={clear}>Clear history</Button></div>}
  </Drawer>
}
