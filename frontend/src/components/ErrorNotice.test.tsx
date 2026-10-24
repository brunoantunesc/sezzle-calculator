import { act, cleanup, fireEvent, render, screen } from '@testing-library/react'
import { afterEach, expect, it, vi } from 'vitest'
import { ErrorNotice } from './ErrorNotice'
afterEach(()=>{cleanup();vi.useRealTimers()})
it('expires after six seconds and supports manual close', () => {
 vi.useFakeTimers()
 const dismiss=vi.fn()
 render(<ErrorNotice message="Cannot divide by zero" dismiss={dismiss} />)
 act(()=>vi.advanceTimersByTime(5999))
 expect(dismiss).not.toHaveBeenCalled()
 act(()=>vi.advanceTimersByTime(1))
 expect(dismiss).toHaveBeenCalledTimes(1)
 fireEvent.click(screen.getByRole('button',{name:'Dismiss error'}))
 expect(dismiss).toHaveBeenCalledTimes(2)
})
