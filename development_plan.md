# Development plan
Inspiration from: https://github.com/quii/learn-go-with-tests/blob/main/http-server.md (MIT License)

Code using Go Language (https://golang.org/)

## Features
See the challenge statement: https://github.com/a-canya/GoServer/coding-challenge.pdf

- F0. Get users
- F1. Sign up
- F2. Req friendship
- F3. Accept friendship
- F4. Decline friendship
- F5. List friends

## Feature building steps
These steps are based on TDD guidelines as described in wikipedia (https://en.wikipedia.org/wiki/Test-driven_development#Test-driven_development_cycle)

- A. Add tests (API call)
- _A'. Run tests_
- B. Code
- _B'. Run tests_
- C. Refactor

Note: each feature was meant to have it own ABC cycle. However I've realized during phases B and C that some tests/code was missing and added "subfeatures" named FXY (where X is the feature and Y the subfeature), and each subfeature has its won cycle.