
# AGENTS.md

## Project Overview

Full-stack calculator application built with React, TypeScript, Vite, Material UI, and Go.

The project prioritizes correctness, security, maintainability, testability, and simplicity.

## General Guidelines

- Follow the requirements and acceptance criteria defined in the corresponding GitHub Issues.
- Prefer simple, readable, and maintainable solutions.
- Avoid unnecessary dependencies, abstractions, and frameworks.
- Follow established conventions for each language.
- Keep changes focused on the requested task.
- Do not introduce unrelated features.
- Do not modify API contracts without explicit approval.

## Backend

- Use idiomatic Go and prefer the standard library.
- Keep arithmetic logic independent of HTTP handling.
- Validate all incoming requests server-side.
- Handle errors explicitly and consistently.
- Treat all client input as untrusted.
- Keep the backend stateless.

## Frontend

- Use React with TypeScript and Material UI.
- Keep components small and focused.
- Separate UI state from API communication.
- All mathematical operations must be performed by the Go backend.
- Use localStorage only for calculation history.
- Maintain accessibility and responsive design.

## Testing

- Follow TDD where practical, especially for calculator logic.
- Add or update tests alongside implementation changes.
- Test observable behavior rather than implementation details.
- Cover relevant edge cases and error scenarios.
- Ensure existing tests pass before completing a task.

## Workflow

- Upon receiving a task prompt, log the original prompt verbatim in `AI_USAGE.md` before making any code changes.
- Include a timestamp in ISO 8601 format using the `America/Sao_Paulo` timezone (UTC-03:00).
- Associate each prompt with its corresponding GitHub Issue when applicable.
- Work within the scope of the assigned GitHub Issue.
- Explain significant architectural decisions.
- Run relevant tests, type checks, and formatting tools.
- Report any limitations or incomplete requirements.
- Never claim tests have passed without executing them.
- Do not create commits or push changes unless explicitly requested.
