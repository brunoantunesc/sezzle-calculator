import { act, cleanup, renderHook } from '@testing-library/react'
import { afterEach, expect, it, vi } from 'vitest'
import { useHistory } from './useHistory'
afterEach(() => {cleanup();localStorage.clear();vi.restoreAllMocks()})
it('caps history at 10, reloads it and filters malformed entries', () => {
 const {result,unmount}=renderHook(()=>useHistory())
 act(()=>{for(let i=0;i<12;i++) result.current.add(`value ${i}`,i)})
 expect(result.current.entries).toHaveLength(10)
 expect(result.current.entries.map(entry=>entry.result)).toEqual([11,10,9,8,7,6,5,4,3,2])
 unmount()
 const restored=renderHook(()=>useHistory())
 expect(restored.result.current.entries).toHaveLength(10)
 restored.unmount()
 localStorage.setItem('calculator-history','[null,{"result":"3"}]')
 expect(renderHook(()=>useHistory()).result.current.entries).toEqual([])
})
it('tolerates invalid JSON and unavailable storage', () => {
 localStorage.setItem('calculator-history','broken')
 expect(renderHook(()=>useHistory()).result.current.entries).toEqual([])
 vi.spyOn(Storage.prototype,'getItem').mockImplementation(()=>{throw new Error('blocked')})
 vi.spyOn(Storage.prototype,'setItem').mockImplementation(()=>{throw new Error('blocked')})
 const {result}=renderHook(()=>useHistory())
 act(()=>result.current.add('1 + 2',3))
 expect(result.current.entries[0].result).toBe(3)
})
