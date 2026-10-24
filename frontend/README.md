# Sezzle calculator frontend

From this directory:

```sh
npm ci
npm run dev
npm test
npm run lint
npm run build
```

The SPA uses React, TypeScript, MUI and Vite. `dist/` contains static production
assets. Run `npm run preview` to inspect a production build locally.

## Current scope

Calculations use the Go API through a dedicated typed service, with relative
`/api/v1/calculations` requests. No arithmetic is performed in the frontend.
Start the Go server with `go run .` in `backend`, then start Vite here. The Vite
proxy forwards `/api` to `http://127.0.0.1:8080`; production uses the same-origin
`/api` endpoint. `npm run preview` only serves static assets, so API-backed preview
requires a same-origin reverse proxy.

The service validates response shapes and finite results, parses recognized API
error codes/messages, and handles network failures and a 10-second timeout.
Server failures use a generic message rather than exposing internal details.
Errors appear in bottom-center notices with a close button and six-second expiry
(paused while hovered or the browser is inactive by MUI). Input remains editable.

`useCalculator` accepts an asynchronous `Calculate` function; tests inject mock
responses to verify requests, chaining, loading, failures, duplicate prevention and
stale-response handling. Unmounting invalidates pending work; late responses cannot replace
new input or enter history.

History keeps the last 10 successful calculations in localStorage, newest first,
with expression, result and timestamp. Entries can be reused and history cleared.
Invalid or unavailable storage is handled without preventing calculations.
The MUI Drawer opens from the right at 600px and above, or from the bottom on
mobile. It is size-limited and scrollable, with focus trapping/restoration and
close-button, backdrop and Escape dismissal. Calculator shortcuts are suspended
while history is open; closing it preserves input.

## Interface and input

- `Calculator`, `Display` and `Keypad` separate composition from presentation.
- Inputs remain strings until submission, preserving `0.` and decimal zeros.
- Entry is limited to 16 digits for usability. This does not guarantee decimal
  precision: submission uses JavaScript numbers to match the backend float64 contract.
- Returned numbers use their original JavaScript string representation, without
  rounding or truncation. Long values shrink and wrap to avoid layout overflow.
- Binary operations chain sequentially, without expression parsing or precedence.
  Selecting another operator before entry replaces the pending operator. Equals
  waits for a second operand and does not repeat the last calculation.
- Square root submits the current displayed operand immediately. When editing
  a second operand, its returned root becomes that operand.
- Digits after a result begin a new entry; an operator reuses the result.
- While pending, equals is disabled. Other keys retain tactile feedback but all input, including AC and square root, is ignored until the request settles. Failed requests preserve editable input.

Keyboard: digits, `.`, `+`, `-`, `*` or `x`, `/`, `^`, `r` for square root,
Enter or `=` for equals, Backspace, Escape or Delete for AC. Modifiers and editable
fields are left alone. Enter/Space on a focused button retain native activation.
Keyboard-held keys receive the same visual depression as pointer-held buttons;
release or window blur removes it. Tab navigation has visible focus outlines.

Layout was checked at 320px and 390px mobile widths and a 1280×800 desktop viewport.
Narrow/short screens can scroll vertically. Reduced-motion preferences are respected.

Vitest and React Testing Library run in jsdom. The unused nested `frontend/`
scaffold is excluded from this app's lint scope, without modifying its contents.

Loading uses a spinner beside the centered equals sign only after 200ms. Faster requests never show it; each request owns a cancellable delay. There is no artificial minimum request duration.
