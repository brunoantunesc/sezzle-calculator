import { cleanup, fireEvent, render, screen, waitFor } from '@testing-library/react'
import { afterEach, expect, it, vi } from 'vitest'
import '@testing-library/jest-dom/vitest'
import App from './App'

afterEach(()=>{cleanup();localStorage.clear();vi.unstubAllGlobals()})

it('renders HTTP results and saves only successful requests to history', async()=>{
 const fetch=vi.fn()
  .mockResolvedValueOnce(new Response('{"result":5}',{status:200}))
  .mockResolvedValueOnce(new Response('{"error":{"code":"DIVISION_BY_ZERO","message":"Cannot divide by zero"}}',{status:422}))
 vi.stubGlobal('fetch',fetch)
 render(<App />)
 for(const key of ['1','0','/','2','Enter']) fireEvent.keyDown(window,{key})
 await waitFor(()=>expect(screen.getByLabelText('Calculator display')).toHaveTextContent('5'))
 expect(JSON.parse(fetch.mock.calls[0][1].body)).toEqual({operation:'divide',operands:[10,2]})
 const output=screen.getByRole('status')
 expect(output).toHaveAttribute('aria-live','polite')
 expect(output).toHaveTextContent('5')
 for(const key of ['/','0','Enter']) fireEvent.keyDown(window,{key})
 expect(await screen.findByRole('alert')).toHaveTextContent('Cannot divide by zero')
 expect(JSON.parse(localStorage.getItem('calculator-history')!)).toEqual([
  expect.objectContaining({expression:'10 ÷ 2',result:5,timestamp:expect.any(String)}),
 ])
 fireEvent.click(screen.getByRole('button',{name:'History'}))
 expect(screen.getByRole('dialog',{name:'Calculation history'})).toHaveTextContent('10 ÷ 2')
 expect(screen.getAllByRole('button',{name:'Reuse result 5'})).toHaveLength(1)
})

it.each([
 ['+','add'],['-','subtract'],['*','multiply'],['/','divide'],['^','power'],
])('maps the %s keyboard operator to %s',async(key,operation)=>{
 const fetch=vi.fn().mockResolvedValue(new Response('{"result":123}',{status:200}))
 vi.stubGlobal('fetch',fetch)
 render(<App />)
 for(const input of ['2',key,'3','Enter']) fireEvent.keyDown(window,{key:input})
 await waitFor(()=>expect(screen.getByRole('status')).toHaveTextContent('123'))
 expect(JSON.parse(fetch.mock.calls[0][1].body)).toEqual({operation,operands:[2,3]})
})
