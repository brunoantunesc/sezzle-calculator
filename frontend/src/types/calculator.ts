export type Operation = 'add' | 'subtract' | 'multiply' | 'divide' | 'power' | 'sqrt'
export interface CalculationRequest { operation: Operation; operands: number[] }
// Supply an API service here when backend integration is enabled.
export type Calculate = (request: CalculationRequest) => Promise<number>
export const symbols: Record<Operation, string> = {
  add: '+', subtract: '−', multiply: '×', divide: '÷', power: '^', sqrt: '√',
}
