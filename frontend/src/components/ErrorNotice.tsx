import { Alert, IconButton, Snackbar, SvgIcon } from '@mui/material'

export function ErrorNotice({ message, dismiss }: { message: string; dismiss: () => void }) {
  return <Snackbar key={message} open={Boolean(message)} autoHideDuration={6000}
    anchorOrigin={{ vertical: 'bottom', horizontal: 'center' }}
    onClose={(_, reason) => { if (reason !== 'clickaway') dismiss() }}>
    <Alert severity="error" variant="filled" className="error-notice"
      action={<IconButton size="small" color="inherit" aria-label="Dismiss error" onClick={dismiss}>
        <SvgIcon fontSize="small"><path d="m6.4 5 5.6 5.6L17.6 5 19 6.4 13.4 12l5.6 5.6-1.4 1.4-5.6-5.6L6.4 19 5 17.6l5.6-5.6L5 6.4Z" /></SvgIcon>
      </IconButton>}>{message}</Alert>
  </Snackbar>
}
