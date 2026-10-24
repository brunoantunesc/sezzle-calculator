import { afterEach, expect, it, vi } from 'vitest'
import { calculate } from './calculatorApi'
afterEach(() => { vi.unstubAllGlobals(); vi.useRealTimers() })
it('posts the contract to the relative API and returns a finite result', async () => {
 const fetch=vi.fn().mockResolvedValue(new Response(JSON.stringify({result:5}),{status:200}))
 vi.stubGlobal('fetch',fetch)
 await expect(calculate({operation:'divide',operands:[10,2]})).resolves.toBe(5)
 expect(fetch).toHaveBeenCalledWith('/api/v1/calculations',expect.objectContaining({method:'POST',headers:{'Content-Type':'application/json'},body:'{"operation":"divide","operands":[10,2]}'}))
})
it('parses domain error messages and normalizes unsafe or broken responses', async () => {
 const fetch=vi.fn();vi.stubGlobal('fetch',fetch)
 fetch.mockResolvedValueOnce(new Response(JSON.stringify({error:{code:'DIVISION_BY_ZERO',message:'Cannot divide by zero'}}),{status:422}))
 await expect(calculate({operation:'divide',operands:[1,0]})).rejects.toThrow('Cannot divide by zero')
 for(const body of ['not json','{"result":null}','{"result":"5"}']) {
 fetch.mockResolvedValueOnce(new Response(body,{status:200}))
 await expect(calculate({operation:'sqrt',operands:[4]})).rejects.toThrow('The server returned an invalid response.')
 }
 fetch.mockResolvedValueOnce(new Response('{"error":{"code":"INTERNAL_ERROR","message":"secret"}}',{status:500}))
 await expect(calculate({operation:'sqrt',operands:[4]})).rejects.toThrow('The server is unavailable. Please try again.')
 fetch.mockRejectedValueOnce(new TypeError('fetch failed'))
 await expect(calculate({operation:'sqrt',operands:[4]})).rejects.toThrow('Could not reach the server. Please try again.')
})

it('times out a stalled request without waiting for real time', async () => {
 vi.useFakeTimers()
 vi.stubGlobal('fetch',vi.fn((_url,options:RequestInit)=>new Promise((_resolve,reject)=>{
  options.signal!.addEventListener('abort',()=>reject(new DOMException('Aborted','AbortError')))
 })))
 const outcome=expect(calculate({operation:'sqrt',operands:[4]})).rejects.toThrow('The request timed out. Please try again.')
 await vi.advanceTimersByTimeAsync(10000)
 await outcome
 expect(vi.getTimerCount()).toBe(0)
})

it.each([
 null,
 {error:null},
 {error:{code:'UNKNOWN',message:'private details'}},
 {error:{code:'INVALID_OPERANDS',message:''}},
 {error:{code:'INVALID_OPERANDS',message:'x'.repeat(301)}},
])('does not expose unexpected error payloads: %j', async body=>{
 vi.stubGlobal('fetch',vi.fn().mockResolvedValue(new Response(JSON.stringify(body),{status:400})))
 await expect(calculate({operation:'add',operands:[]})).rejects.toThrow('Unable to calculate. Please try again.')
})

it('preserves the rate-limit message', async()=>{
 vi.stubGlobal('fetch',vi.fn().mockResolvedValue(new Response(JSON.stringify({error:{code:'RATE_LIMIT_EXCEEDED',message:'Rate limit exceeded'}}),{status:429})))
 await expect(calculate({operation:'sqrt',operands:[4]})).rejects.toThrow('Rate limit exceeded')
})
