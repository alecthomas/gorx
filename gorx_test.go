package gorx

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func add(a, b int) int {
	return a + b
}

//go:generate go run ./cmd/gorx/main.go -o gorx.go gorx bool rune byte string uint int uint8 int8 uint16 int16 uint32 int32 uint64 int64 float32 float64 complex64 complex128 time.Time time.Duration

func TestFromAndToArray(t *testing.T) {
	t.Parallel()
	a := []int{1, 2, 3, 4, 5}
	b := FromIntArray(a).ToArray()
	assert.Equal(t, a, b)
}

func TestJust(t *testing.T) {
	t.Parallel()
	a := JustInt(1).ToArray()
	assert.Equal(t, []int{1}, a)
}

func TestDistinct(t *testing.T) {
	t.Parallel()
	a := FromInts(1, 1, 2, 2, 3, 2, 4, 5).Distinct().ToArray()
	assert.Equal(t, []int{1, 2, 3, 4, 5}, a)
}

func TestResubscribe(t *testing.T) {
	t.Parallel()
	expected := []int{1, 2, 3, 4}
	actual := FromIntArray(expected)
	assert.Equal(t, actual.ToArray(), expected)
	assert.Equal(t, actual.ToArray(), expected)
}

func TestElementAt(t *testing.T) {
	t.Parallel()
	a := FromInts(1, 2, 3, 4).ElementAt(2).ToArray()
	assert.Equal(t, []int{3}, a)
}

func TestFilter(t *testing.T) {
	t.Parallel()
	even := func(i int) bool { return i%2 == 0 }
	a := FromInts(1, 2, 3, 4, 5, 6, 7, 8).Filter(even).ToArray()
	assert.Equal(t, []int{2, 4, 6, 8}, a)
}

func TestFirst(t *testing.T) {
	t.Parallel()
	a := FromInts(1, 2, 3, 4).First().ToArray()
	assert.Equal(t, a, []int{1})
}

func TestLast(t *testing.T) {
	t.Parallel()
	a := FromInts(1, 2, 3, 4).Last().ToArray()
	assert.Equal(t, a, []int{4})
}

func TestMap(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	a := FromInts(1, 2, 3, 4, 5).SkipLast(2).ToArray()
	assert.Equal(t, []int{1, 2, 3}, a)
}

func TestUnsubscribe(t *testing.T) {
	t.Parallel()
	var s GenericSubscription
	s.Unsubscribe()
	assert.True(t, s.Unsubscribed())
}

func TestAverageInt(t *testing.T) {
	t.Parallel()
	a := FromInts(1, 2, 3, 4, 5).Average().ToArray()
	assert.Equal(t, []int{3}, a)
}

func TestAverageFloat32(t *testing.T) {
	t.Parallel()
	b := FromFloat32s(1, 2, 3, 4).Average().ToArray()
	assert.Equal(t, []float32{2.5}, b)
}

func TestSumInt(t *testing.T) {
	t.Parallel()
	a := FromInts(1, 2, 3, 4, 5).Sum().ToArray()
	assert.Equal(t, []int{15}, a)
}

func TestSumFloat32(t *testing.T) {
	t.Parallel()
	a := FromFloat32s(1, 2, 3, 4.5).Sum().ToArray()
	assert.Equal(t, []float32{10.5}, a)
}

func TestCount(t *testing.T) {
	t.Parallel()
	a := FromInts(1, 2, 3, 4, 5, 6, 7).Count().ToArray()
	assert.Equal(t, []int{7}, a)
}

func TestToOne(t *testing.T) {
	t.Parallel()
	_, err := FromInts(1, 2).ToOneWithError()
	assert.Error(t, err)
	value, err := FromInts(3).ToOneWithError()
	assert.NoError(t, err)
	assert.Equal(t, 3, value)
}

func TestMin(t *testing.T) {
	t.Parallel()
	value, err := FromInts(5, 4, 3, 2, 1, 2, 3, 4, 5).Min().ToOneWithError()
	assert.NoError(t, err)
	assert.Equal(t, 1, value)
}

func TestMax(t *testing.T) {
	t.Parallel()
	value, err := FromInts(4, 5, 4, 3, 2, 1, 2).Max().ToOneWithError()
	assert.NoError(t, err)
	assert.Equal(t, 5, value)
}

func TestToChannel(t *testing.T) {
	t.Parallel()
	expected := []int{1, 2, 3, 4, 5, 4, 3, 2, 1}
	a := FromIntArray(expected).ToChannel()
	b := []int{}
	for i := range a {
		b = append(b, i)
	}
	assert.Equal(t, expected, b)
}

func TestReduce(t *testing.T) {
	t.Parallel()
	a := FromInts(1, 2, 3, 4, 5).Reduce(0, add).ToOne()
	assert.Equal(t, 15, a)
}

func TestDo(t *testing.T) {
	t.Parallel()
	b := []int{}
	a := FromInts(1, 2, 3, 4, 5).Do(func(v int) {
		b = append(b, v)
	}).ToArray()
	assert.Equal(t, b, a)
}

func TestThrow(t *testing.T) {
	t.Parallel()
	_, err := ThrowInt(errors.New("error")).ToArrayWithError()
	assert.Error(t, err)
}

func TestEmpty(t *testing.T) {
	t.Parallel()
	a, err := EmptyInt().ToArrayWithError()
	assert.NoError(t, err)
	assert.Equal(t, []int{}, a)
}

func TestDoOnError(t *testing.T) {
	t.Parallel()
	var oerr error
	_, err := ThrowInt(errors.New("error")).DoOnError(func(err error) { oerr = err }).ToArrayWithError()
	assert.Equal(t, err, oerr)
}

func TestDoOnComplete(t *testing.T) {
	t.Parallel()
	complete := false
	a, err := EmptyInt().DoOnComplete(func() { complete = true }).ToArrayWithError()
	assert.NoError(t, err)
	assert.Equal(t, []int{}, a)
	assert.True(t, complete)
}

func TestReplay(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
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
	t.Parallel()
	s := Range(0, 5)
	a := s.ToArray()
	b := s.ToArray()
	assert.Equal(t, []int{0, 1, 2, 3, 4}, a)
	assert.Equal(t, []int{0, 1, 2, 3, 4}, b)
}

func TestRepeat(t *testing.T) {
	t.Parallel()
	s := RepeatInt(5, 3)
	a := s.ToArray()
	b := s.ToArray()
	assert.Equal(t, []int{5, 5, 5}, a)
	assert.Equal(t, []int{5, 5, 5}, b)
}

func TestStart(t *testing.T) {
	t.Parallel()
	s := StartInt(func() int { return 42 })
	a := s.ToArray()
	b := s.ToArray()
	assert.Equal(t, []int{42}, a)
	assert.Equal(t, []int{42}, b)
}

func TestScan(t *testing.T) {
	t.Parallel()
	a := FromInts(1, 2, 3, 4, 5).Scan(0, add).ToArray()
	assert.Equal(t, []int{1, 3, 6, 10, 15}, a)
}

func TestSubscribeNext(t *testing.T) {
	t.Parallel()
	wait := make(chan bool)
	a := []int{}
	_ = FromInts(1, 2, 3, 4, 5).
		DoOnComplete(func() { wait <- true }).
		SubscribeNext(func(v int) { a = append(a, v) })
	<-wait
	assert.Equal(t, []int{1, 2, 3, 4, 5}, a)
}

func TestTake(t *testing.T) {
	t.Parallel()
	s := FromInts(1, 2, 3, 4, 5)
	a := s.Take(2).ToArray()
	b := s.Take(3).ToArray()
	assert.Equal(t, []int{1, 2}, a)
	assert.Equal(t, []int{1, 2, 3}, b)
}

func TestTakeLast(t *testing.T) {
	t.Parallel()
	s := FromInts(1, 2, 3, 4, 5)
	a := s.TakeLast(2).ToArray()
	b := s.TakeLast(3).ToArray()
	assert.Equal(t, []int{4, 5}, a)
	assert.Equal(t, []int{3, 4, 5}, b)
}

func TestIgnoreElements(t *testing.T) {
	t.Parallel()
	s := FromInts(1, 2, 3, 4, 5)
	a := s.IgnoreElements().ToArray()
	assert.Equal(t, []int{}, a)
}

type tiStruct struct {
	v int
	e bool // True if elapsed time was >= 10ms
}

func TestInterval(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
	a := Interval(time.Millisecond * 90).Sample(time.Millisecond * 200).Take(3).ToArray()
	assert.Equal(t, []int{1, 3, 5}, a)
}

func TestDebounce(t *testing.T) {
	t.Parallel()
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
	t.Parallel()
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

func TestConcat(t *testing.T) {
	t.Parallel()
	a := []int{1, 2, 3}
	b := []int{4, 5}
	c := []int{6, 7}
	s := FromIntArray(a).Concat(FromIntArray(b)).Concat(FromIntArray(c)).ToArray()
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, s)
	s = FromIntArray(a).Concat(FromIntArray(b), FromIntArray(c)).ToArray()
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, s)
}

func TestRecover(t *testing.T) {
	t.Parallel()
	merged := FromInts(1, 2, 3).
		Concat(ThrowInt(errors.New("error"))).
		Catch(FromInts(4, 5))
	assert.Equal(t, []int{1, 2, 3, 4, 5}, merged.ToArray())
}

func TestRetry(t *testing.T) {
	t.Parallel()
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
