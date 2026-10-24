import type { Calculate } from '../types/calculator'

export class CalculationError extends Error {}

const clientCodes = new Set([
  'INVALID_REQUEST', 'INVALID_OPERANDS', 'UNSUPPORTED_OPERATION', 'UNSUPPORTED_MEDIA_TYPE',
  'DIVISION_BY_ZERO', 'INVALID_RESULT', 'PAYLOAD_TOO_LARGE', 'RATE_LIMIT_EXCEEDED', 'METHOD_NOT_ALLOWED',
])

export const calculate: Calculate = async request => {
  const controller = new AbortController()
  const timeout = setTimeout(() => controller.abort(), 10000)
  try {
    let response: Response
    try {
      response = await fetch('/api/v1/calculations', {
        method: 'POST', headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(request), signal: controller.signal,
      })
    } catch {
      throw new CalculationError(controller.signal.aborted
        ? 'The request timed out. Please try again.'
        : 'Could not reach the server. Please try again.')
    }
    if (response.status >= 500) throw new CalculationError('The server is unavailable. Please try again.')
    let data: unknown
    try { data = await response.json() } catch {
      throw new CalculationError('The server returned an invalid response.')
    }
    if (!response.ok) {
      if (data && typeof data === 'object' && 'error' in data) {
        const error = data.error
        if (error && typeof error === 'object' && 'code' in error && 'message' in error &&
          typeof error.code === 'string' && clientCodes.has(error.code) &&
          typeof error.message === 'string' && error.message.trim() && error.message.length <= 300) {
          throw new CalculationError(error.message)
        }
      }
      throw new CalculationError('Unable to calculate. Please try again.')
    }
    if (!data || typeof data !== 'object' || !('result' in data) ||
      typeof data.result !== 'number' || !Number.isFinite(data.result)) {
      throw new CalculationError('The server returned an invalid response.')
    }
    return data.result
  } finally { clearTimeout(timeout) }
}
