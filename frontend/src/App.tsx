import { createTheme, CssBaseline, ThemeProvider } from '@mui/material'
import { Calculator } from './components/Calculator'
import type { Calculate } from './types/calculator'
import './App.css'

const theme = createTheme({
  palette: { mode: 'light', primary: { main: '#285e52' }, background: { default: '#eeeae2', paper: '#f7f5ef' }, text: { primary: '#252923' } },
  typography: { fontFamily: '"Segoe UI", system-ui, sans-serif', button: { textTransform: 'none' } },
  shape: { borderRadius: 16 },
})

export default function App({ calculate }: { calculate?: Calculate }) {
  return <ThemeProvider theme={theme}><CssBaseline />
    <main className="app-shell">
      <div className="calculator-stage"><Calculator calculate={calculate} /></div>
    </main>
  </ThemeProvider>
}
