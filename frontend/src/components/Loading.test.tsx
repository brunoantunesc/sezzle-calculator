import { act, cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, expect, it, vi } from 'vitest'
import '@testing-library/jest-dom/vitest'
import App from '../App'

afterEach(() => { cleanup(); localStorage.clear(); vi.useRealTimers() })

it('delays the equals spinner, keeps other keys clickable and ignores pending input', async () => {
 vi.useFakeTimers()
 let resolve!: (value: number) => void
 const calculate=vi.fn(() => new Promise<number>(done => {resolve=done}))
 render(<App calculate={calculate} />)
 for(const name of ['9','Square root']) fireEvent.click(screen.getByRole('button',{name}))
 const equals=screen.getByRole('button',{name:'Equals'})
 expect(equals).toBeDisabled()
 expect(equals).toHaveTextContent('=')
 expect(screen.queryByRole('progressbar')).not.toBeInTheDocument()
 expect(screen.queryByText('Calculating…')).not.toBeInTheDocument()
 act(()=>vi.advanceTimersByTime(199))
 expect(screen.queryByRole('progressbar')).not.toBeInTheDocument()
 act(()=>vi.advanceTimersByTime(1))
 expect(screen.getByRole('progressbar',{name:'Calculating'})).toBeInTheDocument()
 for(const name of ['All clear','Square root','7']) {
  const button=screen.getByRole('button',{name})
  expect(button).toBeEnabled()
  fireEvent.click(button)
 }
 fireEvent.keyDown(window,{key:'Escape'})
 expect(calculate).toHaveBeenCalledTimes(1)
 expect(screen.getByLabelText('Calculator display')).toHaveTextContent('9')
 await act(async()=>resolve(3))
 expect(screen.queryByRole('progressbar')).not.toBeInTheDocument()
 expect(equals).toBeEnabled()
 expect(screen.getByLabelText('Calculator display')).toHaveTextContent('3')
})

it('never flashes a spinner for a 100ms request or leaks its timer into the next request', async () => {
 vi.useFakeTimers()
 let resolve!: (value: number) => void
 render(<App calculate={() => new Promise<number>(done => {resolve=done})} />)
 fireEvent.click(screen.getByRole('button',{name:'Square root'}))
 act(()=>vi.advanceTimersByTime(100))
 await act(async()=>resolve(0))
 act(()=>vi.advanceTimersByTime(100))
 expect(screen.queryByRole('progressbar')).not.toBeInTheDocument()
 fireEvent.click(screen.getByRole('button',{name:'Square root'}))
 act(()=>vi.advanceTimersByTime(199))
 expect(screen.queryByRole('progressbar')).not.toBeInTheDocument()
 await act(async()=>resolve(0))
})
