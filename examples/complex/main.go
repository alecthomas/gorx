package main

//go:generate go run ../../cmd/gorx/main.go -o rx/rx.go --import=net/http rx *http.Response string

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/alecthomas/gorx/examples/complex/rx"
)

func GetCached(url string) *rx.ResponseStream {
	fmt.Printf("No cache entry for %s\n", url)
	return rx.ThrowResponse(errors.New("not implemented"))
}
func SetCached(response *http.Response) {
	fmt.Printf("Caching %s\n", response.Request.URL)
}

func Get(url string) *rx.ResponseStream {
	return rx.CreateResponse(func(observer rx.ResponseObserver, subscription rx.Subscription) {
		response, err := http.Get(url)
		if err != nil {
			observer.Error(err)
		}
		if response.StatusCode < 200 || response.StatusCode > 299 {
			observer.Error(errors.New(http.StatusText(response.StatusCode)))
			return
		}
		observer.Next(response)
		observer.Complete()
	})
}

func URLForArticle(article string) string {
	return "http://en.wikipedia.org/wiki/" + article
}

func LogError(err error) {
	fmt.Printf("error: %s\n", err)
}

func GetWikipediaArticles(timeout time.Duration, articles ...string) *rx.ResponseStream {
	// Try cached URL first, then recover with remote URL and
	// finally recover with an empty stream.
	return rx.FromStringArray(articles).
		Map(URLForArticle).
		FlatMapResponse(func(url string) rx.ResponseObservable {
			remote := Get(url).
				Timeout(timeout).
				Do(SetCached).
				DoOnError(LogError).
				Catch(rx.EmptyResponse())
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
