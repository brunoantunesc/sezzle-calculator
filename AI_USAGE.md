
## 2026-10-23T20:00:29.0492863-03:00 — GitHub Issue #1

The calculator logic's tests are all passing, so lets go to the new step, which is creating the HTTP endpoints and integrating the calculator.
The issue already has the REST contract for the API, including the calculations endpoint and a health check, which is good measure for tests / docker validation.

Take utmost care in the requisition validation, as in bad JSON, missing, unexpected or incorrectly typed fields, payload limit and forbidden HTTP methods. Error codes are important.

Keep the TDD approach for this implementation.

When finished, serve me with an overview of the changes and decisions.

## 2026-10-23T20:11:34.6203853-03:00 — GitHub Issue #1

I reviewed the implementation, and there are a few details that should be improved. requestFields has to consume correctly the JSON closing symbol, validate the quantity of operands as soon as possible, and also add a few tests for empty operation, empty array and payload above the limit, to make the code more secure against regressions.
