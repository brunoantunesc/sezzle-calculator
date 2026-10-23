
## 2026-10-23T20:00:29.0492863-03:00 — GitHub Issue #1

The calculator logic's tests are all passing, so lets go to the new step, which is creating the HTTP endpoints and integrating the calculator.
The issue already has the REST contract for the API, including the calculations endpoint and a health check, which is good measure for tests / docker validation.

Take utmost care in the requisition validation, as in bad JSON, missing, unexpected or incorrectly typed fields, payload limit and forbidden HTTP methods. Error codes are important.

Keep the TDD approach for this implementation.

When finished, serve me with an overview of the changes and decisions.

## 2026-10-23T20:11:34.6203853-03:00 — GitHub Issue #1

I reviewed the implementation, and there are a few details that should be improved. requestFields has to consume correctly the JSON closing symbol, validate the quantity of operands as soon as possible, and also add a few tests for empty operation, empty array and payload above the limit, to make the code more secure against regressions.

## 2026-10-23T20:21:13.6907996-03:00 — GitHub Issue #1

Now, to finish the backend, we need to focus on configurating the server correctly and adding some security protections.

Implement server timeouts, header limits, configurable CORS, security headers and a basic rate limiting of 60 reqs/min/ip. This should make sure that, despite being a simple app, it is secure enough against any malicious intent. Be careful about concurrency and memory issues in rate limiting.

Keep the TDD model where possible, and make sure that the tests are more deterministic, avoiding to use time.Sleep which will make the tests take a lot of time to run, instead inject a clock function in the test.

Also make sure to test if different IPs are not being limit together and if there are no memory leaks, making sure to always clean old entries. No need for a complex abstraction tho.

## 2026-10-23T20:36:42.9878435-03:00 — GitHub Issue #1

Implementation is almost done, but there are a few things that need improvement.

The rate limiter is searching all registered IPs with the allow() method. This can induce heavy load under pressure, as it is O(n) per request, and you should fix it.

The best way to handle this is adding a new field nextCleanup, and only check the map when now >= nextCleanup. After cleanup, always put this field to now +1 min. Expired entries should also be treated correctly, because if the IP is already scheduled to be cleaned up, it should be allowed to add a request.

Also, Retry-After must be improved. The middleware is always return 60 when a client exceeds the limit, but it should return the time until the windows expires for that IP. This will also require the allow() method to return not only bool, but it is a simple change.

Lastly, the CSP currently will not work in the future when the Go server needs to serve the static files for React. You can keep it for the API responses, but do not enforce it to the frontend. No need to implement the serving now, but organized the security headers to make it easier to differ both uses.

Make sure to also update the tests to reflect these changes.
