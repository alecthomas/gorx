# Reactive eXtensions for Go

This is a package and tool providing [Reactive eXtensions](http://reactivex.io) for Go.

The main package provides Rx operators for common Go builtin types, and the
tool generates Go code for arbitrary types.

# Why?

Yes, good question. Mostly as an exercise to see if it was feasible/possible.
It is, largely, except for operators that produce sequence types (ie. arrays
or observables of T).

# Installation

```
go get github.com/alecthomas/gorx github.com/alecthomas/gorx/cmd/gorx
```

# Usage

To use the package:

```go
import "github.com/alecthomas/gorx"

gorx.
  FromTimeChannel(time.Tick(time.Second)).
  Take(5).
  Do(func(t time.Time) { fmt.Printf("%s\n", t) }).
  Wait()
```

To generate Rx operators for custom types:

```
gorx --import=gopkg.in/alecthomas/kingpin.v2 kingpinrx '*kingpin.CmdClause' '*kingpin.FlagClause'
```

# Examples

Simple:

```go
gorx.FromStrings("Ben", "George").Do(func(s string) { fmt.Println(s) }).Wait()
```

# Operators

Following are a list of the core operators as [defined by reactivex.io](http://reactivex.io/documentation/operators.html) that have been implemented and will (probably) be implemented soon:

## Create operators

- Create
- Empty
- Never
- Throw
- Just
- Range
- Repeat
- Start

Not implemented:

- Defer
- Timer

## Transformations

- Map
- Reduce
- Scan

Not implemented:

- Buffer
- FlatMap
- GroupBy
- Window

Note: These operators are currently not implemented because each distinct
observable type requires quite a lot of boilerplate code, and these operators
produce new types. eg. `.Buffer(2)` would transform `T` to a stream of `[]T`,
`.FlatMap(f)` would transform `T` to a `TStreamStream`, etc. One "solution" is
to only generate these operators if the user explicitly requests these
resultant types.

## Filters

- Distinct
- ElementAt
- Filter
- First
- Last
- Skip
- SkipLast
- Take
- TakeLast
- IgnoreElements
- Sample
- Debounce

# Combining

- Merge

Not implemented:

- MergeDelayError
- CombineLatest
- And / Then / When
- Zip
- Join
- StartWith
- Switch
- Zip

Note: See note above for Transformations for why these are not implemented.

# Error handling

- Catch
- Retry

# Mathematics and Aggregation

- Concat
- Average
- Count
- Min
- Max
- Reduce
- Sum

# Utility

- Do
- Subscribe

Not implemented:

- Delay
- Timeout
- Timestamp
- Materialize / Dematerialize
- Serialize
- TimeInterval

# Conditional and Boolean

Not implemented:

- All
- Amb
- Contains
- DefaultIfEmpty
- SequenceEqual
- SkipUntil
- SkipWhile
- TakeUntil
- TakeWhile

# Conversion

- To (one, array, channel)
