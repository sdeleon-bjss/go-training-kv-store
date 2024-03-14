# BJSS Go Academy - KV Store

## Table of Contents
1. [Introduction](#introduction)
2. [How to run, test and benchmark the application](#how-to-run-test-and-benchmark-the-application)
    1. [Run the application](#run-the-application)
    2. [Test the application](#test-the-application)
    3. [Benchmark the application](#benchmark-the-application)

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


## How to run, test and benchmark the application
### Run the application
There are two ways from the root of project:

Using the binary
```shell
go run .
```

Using the `main.go` file:
```shell
go run main.go
```

### Test the application
```shell
go test -v ./store/...
```
- this will run the tests within `store/` for `/store/kv_test.go`
Example output:
```shell
$ go test -v ./store/...
=== RUN   TestCommandSet_Apply
--- PASS: TestCommandSet_Apply (0.00s)     
=== RUN   TestCommandSet_ApplyError        
--- PASS: TestCommandSet_ApplyError (0.00s)
=== RUN   TestCommandGet_Apply
--- PASS: TestCommandGet_Apply (0.00s)
=== RUN   TestCommandGet_ApplyError
--- PASS: TestCommandGet_ApplyError (0.00s)
=== RUN   TestCommandDelete_Apply
--- PASS: TestCommandDelete_Apply (0.00s)
=== RUN   TestCommandDelete_ApplyError
--- PASS: TestCommandDelete_ApplyError (0.00s)
=== RUN   TestActor_Set
--- PASS: TestActor_Set (1.01s)
=== RUN   TestActor_SetError
--- PASS: TestActor_SetError (2.02s)
=== RUN   TestActor_ConcurrentSet
--- PASS: TestActor_ConcurrentSet (0.00s)
=== RUN   TestActor_ConcurrentSetSameKey
--- PASS: TestActor_ConcurrentSetSameKey (0.00s)
=== RUN   TestActor_Get
--- PASS: TestActor_Get (0.00s)
=== RUN   TestActor_GetError
--- PASS: TestActor_GetError (0.00s)
=== RUN   TestActor_ConcurrentGet
2024/03/14 09:35:00 Error applying command: key hello does not exist
--- PASS: TestActor_ConcurrentGet (0.00s)
=== RUN   TestActor_Delete
--- PASS: TestActor_Delete (1.01s)
=== RUN   TestActor_DeleteError
--- PASS: TestActor_DeleteError (1.01s)
=== RUN   TestActor_ConcurrentDelete
2024/03/14 09:35:02 Error applying command: key hello does not exist
2024/03/14 09:35:02 Error applying command: key hello does not exist
2024/03/14 09:35:02 Error applying command: key hello does not exist
2024/03/14 09:35:02 Error applying command: key hello does not exist
2024/03/14 09:35:02 Error applying command: key hello does not exist
2024/03/14 09:35:02 Error applying command: key hello does not exist
2024/03/14 09:35:02 Error applying command: key hello does not exist
2024/03/14 09:35:02 Error applying command: key hello does not exist
2024/03/14 09:35:02 Error applying command: key hello does not exist
--- PASS: TestActor_ConcurrentDelete (0.00s)
=== RUN   TestActor_List
--- PASS: TestActor_List (0.00s)
=== RUN   TestNewActor
--- PASS: TestNewActor (0.00s)
PASS
ok      github.com/sdeleon-bjss/store   (cached)
```

### Benchmark the application
```shell
go test -bench=. ./store/...
```

Example output:
```shell
goos: windows
goarch: amd64
pkg: github.com/sdeleon-bjss/store
cpu: 11th Gen Intel(R) Core(TM) i7-11800H @ 2.30GHz
BenchmarkActor_Set
BenchmarkActor_Set-16                             207596              5691 ns/op

BenchmarkActor_SetConcurrent
BenchmarkActor_SetConcurrent-16                   424105              3155 ns/op

BenchmarkActor_SetConcurrentSameKey
BenchmarkActor_SetConcurrentSameKey-16            432754              2924 ns/op

BenchmarkActor_GetConcurrent
BenchmarkActor_GetConcurrent-16                   249266              6274 ns/op

BenchmarkActor_GetConcurrentSameKey
BenchmarkActor_GetConcurrentSameKey-16           1000000              1022 ns/op

BenchmarkActor_DeleteConcurrent
BenchmarkActor_DeleteConcurrent-16                375348             18734 ns/op

BenchmarkActor_DeleteConcurrentSameKey
BenchmarkActor_DeleteConcurrentSameKey-16         870126             15059 ns/op

BenchmarkActor_ConcurrentMixedOperations
BenchmarkActor_ConcurrentMixedOperations-16       247268             33598 ns/op

BenchmarkActor_StressTestSet
BenchmarkActor_StressTestSet-16                   169128             24904 ns/op
         401.53 MB/s
PASS
```
