# BJSS Go Academy - KV Store

## Introduction
Write a RESTful concurrent key-value store using the actor pattern and Go's communicating sequential processes (CSP).

The service should use the URL path as a key, accept any data passed in the POST body as a value, and use the GET or DELETE verbs to retrieve/delete values. POST and DELETE actions should not wait for their actions to be completed before returning but should be concurrently safe.

### Required

1. The solution should follow the established idiomatic Golang code styles.
2. The solution should be simple and readable.
3. The solution should include meaningful unit tests that prove the solution works, including corner cases. There is no need to include tests purely for coverage.
Optional

### _Optional_
* The service should be meaningfully performance-tested and benchmarked.
* When a shutdown signal is received, all existing connections should be fulfilled before the application shuts down.