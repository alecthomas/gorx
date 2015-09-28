package main

//go:generate go run ../../cmd/gorx/main.go -o rx/rx.go --import=net/http rx *http.Response string

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/alecthomas/gorx/example/complex/rx"
)

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

func main() {
	GetWikipediaArticles(time.Second*5, "MinHash", "Streaming_algorithm", "A MISSING PAGE", "ANOTHER MISSING PAGE").
		Do(func(response *http.Response) {
		defer response.Body.Close()
		content, _ := ioutil.ReadAll(response.Body)
		fmt.Printf("Retrieved %d bytes from %s\n", len(content), response.Request.URL)
	}).
		Wait()
}
