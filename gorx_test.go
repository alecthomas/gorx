package gorx

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"testing"
	"time"

	"github.com/alecthomas/assert"
)

func add(a, b int) int {
	return a + b
}

//go:generate go run ./cmd/gorx/main.go --base-types -o gorx.go gorx

func TestFromAndToArray(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b := FromIntArray(a).ToArray()
	assert.Equal(t, a, b)
}

func TestJust(t *testing.T) {
	a := JustInt(1).ToArray()
	assert.Equal(t, []int{1}, a)
}

func TestDistinct(t *testing.T) {
	a := FromInts(1, 1, 2, 2, 3, 2, 4, 5).Distinct().ToArray()
	assert.Equal(t, []int{1, 2, 3, 4, 5}, a)
}

func TestResubscribe(t *testing.T) {
	expected := []int{1, 2, 3, 4}
	actual := FromIntArray(expected)
	assert.Equal(t, actual.ToArray(), expected)
	assert.Equal(t, actual.ToArray(), expected)
}

func TestElementAt(t *testing.T) {
	a := FromInts(1, 2, 3, 4).ElementAt(2).ToArray()
	assert.Equal(t, []int{3}, a)
}

func TestFilter(t *testing.T) {
	even := func(i int) bool { return i%2 == 0 }
	a := FromInts(1, 2, 3, 4, 5, 6, 7, 8).Filter(even).ToArray()
	assert.Equal(t, []int{2, 4, 6, 8}, a)
}

func TestFirst(t *testing.T) {
	a := FromInts(1, 2, 3, 4).First().ToArray()
	assert.Equal(t, a, []int{1})
}

func TestLast(t *testing.T) {
	a := FromInts(1, 2, 3, 4).Last().ToArray()
	assert.Equal(t, a, []int{4})
}

func TestMap(t *testing.T) {
	a := FromInts(1, 2, 3, 4).MapString(func(i int) string { return fmt.Sprintf("%d!", i) }).ToArray()
	assert.Equal(t, a, []string{"1!", "2!", "3!", "4!"})
}

func FromTestChannel(t *testing.T) {
	ch := make(chan int, 5)
	for i := 0; i < 5; i++ {
		ch <- i
	}
	close(ch)
	a := FromIntChannel(ch).ToArray()
	assert.Equal(t, []int{0, 1, 2, 3, 4}, a)
}

func TestSkipLast(t *testing.T) {
	a := FromInts(1, 2, 3, 4, 5).SkipLast(2).ToArray()
	assert.Equal(t, []int{1, 2, 3}, a)
}

func TestUnsubscribe(t *testing.T) {
	s := NewGenericSubscription()
	assert.False(t, s.Closed())
	s.Close()
	assert.True(t, s.Closed())
}

func TestAverageInt(t *testing.T) {
	a := FromInts(1, 2, 3, 4, 5).Average().ToArray()
	assert.Equal(t, []int{3}, a)
}

func TestAverageFloat32(t *testing.T) {
	b := FromFloat32s(1, 2, 3, 4).Average().ToArray()
	assert.Equal(t, []float32{2.5}, b)
}

func TestSumInt(t *testing.T) {
	a := FromInts(1, 2, 3, 4, 5).Sum().ToArray()
	assert.Equal(t, []int{15}, a)
}

func TestSumFloat32(t *testing.T) {
	a := FromFloat32s(1, 2, 3, 4.5).Sum().ToArray()
	assert.Equal(t, []float32{10.5}, a)
}

func TestCount(t *testing.T) {
	a := FromInts(1, 2, 3, 4, 5, 6, 7).Count().ToArray()
	assert.Equal(t, []int{7}, a)
}

func TestToOneWithError(t *testing.T) {
	_, err := FromInts(1, 2).ToOneWithError()
	assert.Error(t, err)
	value, err := FromInts(3).ToOneWithError()
	assert.NoError(t, err)
	assert.Equal(t, 3, value)
}

// func TestToOneWithErrorCancelsSubscription(t *testing.T) {
// 	var sub Subscription
// 	o := CreateInt(func(observer IntObserver, subscription Subscription) {
// 		sub = subscription
// 		observer.Next(1)
// 		observer.Next(2)
// 		observer.Complete()
// 	})
// 	_, err := o.ToOneWithError()
// 	assert.Error(t, err)
// 	assert.True(t, sub.Closed())
// }

func TestMin(t *testing.T) {
	value, err := FromInts(5, 4, 3, 2, 1, 2, 3, 4, 5).Min().ToOneWithError()
	assert.NoError(t, err)
	assert.Equal(t, 1, value)
}

func TestMax(t *testing.T) {
	value, err := FromInts(4, 5, 4, 3, 2, 1, 2).Max().ToOneWithError()
	assert.NoError(t, err)
	assert.Equal(t, 5, value)
}

func TestToChannel(t *testing.T) {
	expected := []int{1, 2, 3, 4, 5, 4, 3, 2, 1}
	a := FromIntArray(expected).ToChannel()
	b := []int{}
	for i := range a {
		b = append(b, i)
	}
	assert.Equal(t, expected, b)
}

func TestReduce(t *testing.T) {
	a := FromInts(1, 2, 3, 4, 5).Reduce(0, add).ToOne()
	assert.Equal(t, 15, a)
}

func TestDo(t *testing.T) {
	b := []int{}
	a := FromInts(1, 2, 3, 4, 5).Do(func(v int) {
		b = append(b, v)
	}).ToArray()
	assert.Equal(t, b, a)
}

func TestThrow(t *testing.T) {
	_, err := ThrowInt(errors.New("error")).ToArrayWithError()
	assert.Error(t, err)
}

func TestEmpty(t *testing.T) {
	a, err := EmptyInt().ToArrayWithError()
	assert.NoError(t, err)
	assert.Equal(t, []int{}, a)
}

func TestDoOnError(t *testing.T) {
	var oerr error
	_, err := ThrowInt(errors.New("error")).DoOnError(func(err error) { oerr = err }).ToArrayWithError()
	assert.Equal(t, err, oerr)
}

func TestDoOnComplete(t *testing.T) {
	complete := false
	a, err := EmptyInt().DoOnComplete(func() { complete = true }).ToArrayWithError()
	assert.NoError(t, err)
	assert.Equal(t, []int{}, a)
	assert.True(t, complete)
}

func TestReplay(t *testing.T) {
	ch := make(chan int, 5)
	for i := 0; i < 5; i++ {
		ch <- i
	}
	close(ch)
	s := FromIntChannel(ch).Replay(0, 0)
	a := s.ToArray()
	b := s.ToArray()
	expected := []int{0, 1, 2, 3, 4}
	assert.Equal(t, expected, a)
	assert.Equal(t, expected, b)
}

func TestReplayWithSize(t *testing.T) {
	ch := make(chan int, 5)
	for i := 0; i < 5; i++ {
		ch <- i
	}
	close(ch)
	s := FromIntChannel(ch).Replay(2, 0)
	a := s.ToArray()
	b := s.ToArray()
	assert.Equal(t, []int{0, 1, 2, 3, 4}, a)
	assert.Equal(t, []int{3, 4}, b)
	assert.Equal(t, []int{3, 4}, s.ToArray())
}

func TestReplayWithExpiry(t *testing.T) {
	ch := make(chan int)
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
			time.Sleep(time.Millisecond * 100)
		}
		close(ch)
	}()
	s := FromIntChannel(ch).Replay(0, time.Millisecond*600)
	time.Sleep(time.Millisecond * 500)
	assert.Equal(t, []int{0, 1, 2, 3, 4}, s.ToArray())
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, []int{1, 2, 3, 4}, s.ToArray())
}

func TestCreate(t *testing.T) {
	s := CreateInt(func(observer IntObserver, subscription Subscription) {
		observer.Next(0)
		observer.Next(1)
		observer.Next(2)
		observer.Complete()
	})
	a := s.ToArray()
	b := s.ToArray()
	assert.Equal(t, []int{0, 1, 2}, a)
	assert.Equal(t, []int{0, 1, 2}, b)
}

func TestRange(t *testing.T) {
	s := Range(1, 5)
	a := s.ToArray()
	b := s.ToArray()
	assert.Equal(t, []int{1, 2, 3, 4, 5}, a)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, b)
}

func TestRepeat(t *testing.T) {
	s := RepeatInt(5, 3)
	a := s.ToArray()
	b := s.ToArray()
	assert.Equal(t, []int{5, 5, 5}, a)
	assert.Equal(t, []int{5, 5, 5}, b)
}

func TestStart(t *testing.T) {
	s := StartInt(func() (int, error) { return 42, nil })
	a := s.ToArray()
	b := s.ToArray()
	assert.Equal(t, []int{42}, a)
	assert.Equal(t, []int{42}, b)
}

func TestStartWithError(t *testing.T) {
	s := StartInt(func() (int, error) { return 0, errors.New("error") })
	_, err := s.ToArrayWithError()
	assert.Error(t, err)
}

func TestScan(t *testing.T) {
	a := FromInts(1, 2, 3, 4, 5).Scan(0, add).ToArray()
	assert.Equal(t, []int{1, 3, 6, 10, 15}, a)
}

func TestSubscribeNext(t *testing.T) {
	wait := make(chan bool)
	a := []int{}
	_ = FromInts(1, 2, 3, 4, 5).
		DoOnComplete(func() { wait <- true }).
		SubscribeNext(func(v int) { a = append(a, v) })
	<-wait
	assert.Equal(t, []int{1, 2, 3, 4, 5}, a)
}

func TestTake(t *testing.T) {
	s := FromInts(1, 2, 3, 4, 5)
	a := s.Take(2).ToArray()
	b := s.Take(3).ToArray()
	assert.Equal(t, []int{1, 2}, a)
	assert.Equal(t, []int{1, 2, 3}, b)
}

func TestTakeLast(t *testing.T) {
	s := FromInts(1, 2, 3, 4, 5)
	a := s.TakeLast(2).ToArray()
	b := s.TakeLast(3).ToArray()
	assert.Equal(t, []int{4, 5}, a)
	assert.Equal(t, []int{3, 4, 5}, b)
}

func TestIgnoreElements(t *testing.T) {
	s := FromInts(1, 2, 3, 4, 5)
	a := s.IgnoreElements().ToArray()
	assert.Equal(t, []int{}, a)
}

type tiStruct struct {
	v int
	e bool // True if elapsed time was >= 10ms
}

func TestInterval(t *testing.T) {
	seen := []tiStruct{}
	last := time.Now()
	wait := make(chan bool)
	Interval(time.Millisecond * 10).
		Take(5).
		DoOnComplete(func() { wait <- true }).
		SubscribeNext(func(n int) { seen = append(seen, tiStruct{n, time.Now().Sub(last) >= 10*time.Millisecond}) })
	<-wait
	assert.Equal(t, []tiStruct{
		{0, true},
		{1, true},
		{2, true},
		{3, true},
		{4, true},
	}, seen)
}

func TestSample(t *testing.T) {
	a := Interval(time.Millisecond * 90).Sample(time.Millisecond * 200).Take(3).ToArray()
	assert.Equal(t, []int{1, 3, 5}, a)
}

func TestDebounce(t *testing.T) {
	s := CreateInt(func(observer IntObserver, subscription Subscription) {
		time.Sleep(100 * time.Millisecond)
		observer.Next(1)
		time.Sleep(300 * time.Millisecond)
		observer.Next(2)
		time.Sleep(80 * time.Millisecond)
		observer.Next(3)
		time.Sleep(110 * time.Millisecond)
		observer.Next(4)
		observer.Complete()
	})
	a := s.Debounce(time.Millisecond * 100).ToArray()
	assert.Equal(t, []int{1, 3, 4}, a)
}

func TestMerge(t *testing.T) {
	sa := CreateInt(func(observer IntObserver, subscription Subscription) {
		time.Sleep(10 * time.Millisecond)
		observer.Next(1)
		time.Sleep(10 * time.Millisecond)
		observer.Next(3)
		observer.Complete()
	})
	sb := CreateInt(func(observer IntObserver, subscription Subscription) {
		time.Sleep(5 * time.Millisecond)
		observer.Next(0)
		time.Sleep(10 * time.Millisecond)
		observer.Next(2)
		observer.Complete()
	})
	a := sa.Merge(sb).ToArray()
	assert.Equal(t, []int{0, 1, 2, 3}, a)
}

func TestMergeDelayError(t *testing.T) {
	sa := CreateInt(func(observer IntObserver, subscription Subscription) {
		time.Sleep(10 * time.Millisecond)
		observer.Next(1)
		observer.Error(errors.New("error"))
	})
	sb := CreateInt(func(observer IntObserver, subscription Subscription) {
		time.Sleep(5 * time.Millisecond)
		observer.Next(0)
		time.Sleep(10 * time.Millisecond)
		observer.Next(2)
		observer.Complete()
	})
	a := sa.MergeDelayError(sb).ToArray()
	assert.Equal(t, []int{0, 1, 2}, a)
}

func TestConcat(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{4, 5}
	c := []int{6, 7}
	s := FromIntArray(a).Concat(FromIntArray(b)).Concat(FromIntArray(c)).ToArray()
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, s)
	s = FromIntArray(a).Concat(FromIntArray(b), FromIntArray(c)).ToArray()
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, s)
}

func TestRecover(t *testing.T) {
	merged := FromInts(1, 2, 3).
		Concat(ThrowInt(errors.New("error"))).
		Catch(FromInts(4, 5))
	assert.Equal(t, []int{1, 2, 3, 4, 5}, merged.ToArray())
}

func TestRetry(t *testing.T) {
	errored := false
	a := CreateInt(func(observer IntObserver, subscription Subscription) {
		observer.Next(1)
		observer.Next(2)
		observer.Next(3)
		if errored {
			observer.Complete()
		} else {
			observer.Error(errors.New("error"))
			errored = true
		}
	})
	b := a.Retry().ToArray()
	assert.Equal(t, []int{1, 2, 3, 1, 2, 3}, b)
	assert.True(t, errored)
}

func TestLinkedSubscription(t *testing.T) {
	linked := NewLinkedSubscription()
	sub := NewGenericSubscription()
	assert.False(t, linked.Closed())
	assert.False(t, sub.Closed())
	linked.Link(sub)
	assert.Panics(t, func() { linked.Link(sub) })
	linked.Close()
	assert.True(t, sub.Closed())
	assert.True(t, linked.Closed())
}

func TestLinkedSubscriptionUnsubscribesTargetOnLink(t *testing.T) {
	linked := NewLinkedSubscription()
	sub := NewGenericSubscription()
	linked.Close()
	assert.True(t, linked.Closed())
	assert.False(t, sub.Closed())
	linked.Link(sub)
	assert.True(t, linked.Closed())
	assert.True(t, sub.Closed())
}

func TestChannelSubscription(t *testing.T) {
	done := make(chan bool)
	unsubscribed := false
	var s Subscription = NewChannelSubscription()
	events, ok := s.(SubscriptionEvents)
	assert.True(t, ok)
	events.OnUnsubscribe(func() { unsubscribed = true; done <- true })
	assert.False(t, s.Closed())
	s.Close()
	assert.True(t, s.Closed())
	<-done
	assert.True(t, unsubscribed)
}

func TestFlatMap(t *testing.T) {
	actual := Range(1, 2).FlatMap(func(n int) IntObservable { return Range(n, 2) }).ToArray()
	sort.Ints(actual)
	assert.Equal(t, []int{1, 2, 2, 3}, actual)
}

func TestTimeout(t *testing.T) {
	wg := sync.WaitGroup{}
	start := time.Now()
	wg.Add(1)
	actual, err := CreateInt(func(observer IntObserver, subscription Subscription) {
		observer.Next(1)
		time.Sleep(time.Millisecond * 500)
		assert.True(t, subscription.Closed())
		wg.Done()
	}).
		Timeout(time.Millisecond * 250).
		ToArrayWithError()
	elapsed := time.Now().Sub(start)
	assert.Error(t, err)
	assert.Equal(t, ErrTimeout, err)
	assert.True(t, elapsed > time.Millisecond*250 && elapsed < time.Millisecond*500)
	assert.Equal(t, []int{1}, actual)
	wg.Wait()
}

func TestFork(t *testing.T) {
	ch := make(chan int, 30)
	s := FromIntChannel(ch).Fork()
	a := []int{}
	b := []int{}
	s.SubscribeNext(func(n int) { a = append(a, n) })
	s.SubscribeNext(func(n int) { b = append(b, n) })
	ch <- 1
	ch <- 2
	ch <- 3
	close(ch)
	s.Wait()
	assert.Equal(t, []int{1, 2, 3}, a)
	assert.Equal(t, []int{1, 2, 3}, b)
}
