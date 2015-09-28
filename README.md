# Reactive eXtensions for Go

This is a package and tool providing [Reactive eXtensions](http://reactivex.io) for Go.

The main package provides Rx operators for common Go builtin types, and the
tool generates Go code for arbitrary types.

# Why?

Yes, good question. Mostly as an exercise to see if it was feasible/possible.
It is, largely, except for operators that produce sequence types (ie. arrays
or observables of T).

That said, the Rx operators *do* provide some nice functionality that would
otherwise have to be implemented by hand.

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

A very basic example creating an observable from a set of strings and printing
them:

```go
gorx.FromStrings("Ben", "George").Do(func(s string) { fmt.Println(s) }).Wait()
```

A more [complex example](examples/complex/main.go). Try retrieving article from cache, otherwise fetch
original from Wikipedia, all with a timeout.

```go
func GetCached(url string) *ResponseStream {
  fmt.Printf("No cache entry for %s\n", url)
  return ThrowResponse(errors.New("not implemented"))
}
func SetCached(response *http.Response) {
  fmt.Printf("Caching %s\n", response.Request.URL)
}

func Get(url string) *ResponseStream {
  return StartResponse(func() (*http.Response, error) {
    response, err := http.Get(url)
    if err == nil && (response.StatusCode < 200 || response.StatusCode > 299) {
      return nil, errors.New(http.StatusText(response.StatusCode))
    }
    return response, err
  })
}

func URLForArticle(article string) string {
  return "http://en.wikipedia.org/wiki/" + article
}

func LogError(err error) {
  fmt.Printf("error: %s\n", err)
}

func GetWikipediaArticles(timeout time.Duration, articles ...string) *ResponseStream {
  // Try cached URL first, then recover with remote URL and
  // finally recover with an empty stream.
  return FromStringArray(articles).
    Map(URLForArticle).
    FlatMapResponse(func(url string) ResponseObservable {
      remote := Get(url).
        Timeout(timeout).
        Do(SetCached).
        DoOnError(LogError).
        Catch(EmptyResponse())
      return GetCached(url).
        Catch(remote)
    })
}
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
- FlatMap

Not implemented:

- Buffer
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
- MergeDelayError

Not implemented:

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
