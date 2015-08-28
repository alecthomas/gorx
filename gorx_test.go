package gorx

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

//go:generate go run ./cmd/gorx/main.go -o gorx.go gorx rune byte string uint int uint8 int8 uint16 int16 uint32 int32 uint64 int64 float32 float64 complex64 complex128 time.Time time.Duration

func TestToArray(t *testing.T) {
	a := []int{1, 2, 3, 4, 5}
	b := FromIntArray(a).ToArray()
	assert.Equal(t, a, b)
}

func TestDistinct(t *testing.T) {
	a := FromIntArray([]int{1, 1, 2, 2, 3, 2, 4, 5}).Distinct().ToArray()
	assert.Equal(t, []int{1, 2, 3, 4, 5}, a)
}

func TestResubscribe(t *testing.T) {
	expected := []int{1, 2, 3, 4}
	actual := FromIntArray(expected)
	assert.Equal(t, actual.ToArray(), expected)
	assert.Equal(t, actual.ToArray(), expected)
}

func TestElementAt(t *testing.T) {
	a := FromIntArray([]int{1, 2, 3, 4}).ElementAt(2).ToArray()
	assert.Equal(t, []int{3}, a)
}

func TestFilter(t *testing.T) {
	even := func(i int) bool { return i%2 == 0 }
	a := FromIntArray([]int{1, 2, 3, 4, 5, 6, 7, 8}).Filter(even).ToArray()
	assert.Equal(t, []int{2, 4, 6, 8}, a)
}

func TestFirst(t *testing.T) {
	a := FromIntArray([]int{1, 2, 3, 4}).First().ToArray()
	assert.Equal(t, a, []int{1})
}

func TestLast(t *testing.T) {
	a := FromIntArray([]int{1, 2, 3, 4}).Last().ToArray()
	assert.Equal(t, a, []int{4})
}

func TestMap(t *testing.T) {
	a := FromIntArray([]int{1, 2, 3, 4}).MapString(func(i int) string { return fmt.Sprintf("%d!", i) }).ToArray()
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
	a := FromIntArray([]int{1, 2, 3, 4, 5}).SkipLast(2).ToArray()
	assert.Equal(t, []int{1, 2, 3}, a)
}

func TestUnsubscribe(t *testing.T) {
	var s GenericSubscription
	s.Unsubscribe()
	assert.True(t, s.Unsubscribed())
}

func TestAverageInt(t *testing.T) {
	a := FromIntArray([]int{1, 2, 3, 4, 5}).Average().ToArray()
	assert.Equal(t, []int{3}, a)
}

func TestAverageFloat32(t *testing.T) {
	b := FromFloat32Array([]float32{1, 2, 3, 4}).Average().ToArray()
	assert.Equal(t, []float32{2.5}, b)
}

func TestSumInt(t *testing.T) {
	a := FromIntArray([]int{1, 2, 3, 4, 5}).Sum().ToArray()
	assert.Equal(t, []int{15}, a)
}

func TestSumFloat32(t *testing.T) {
	a := FromFloat32Array([]float32{1, 2, 3, 4.5}).Sum().ToArray()
	assert.Equal(t, []float32{10.5}, a)
}

func TestCount(t *testing.T) {
	a := FromIntArray([]int{1, 2, 3, 4, 5, 6, 7}).Count().ToArray()
	assert.Equal(t, []int{7}, a)
}

func TestToOne(t *testing.T) {
	_, err := FromIntArray([]int{1, 2}).ToOneWithError()
	assert.Error(t, err)
	value, err := FromIntArray([]int{3}).ToOneWithError()
	assert.NoError(t, err)
	assert.Equal(t, 3, value)
}

func TestMin(t *testing.T) {
	value, err := FromIntArray([]int{5, 4, 3, 2, 1, 2, 3, 4, 5}).Min().ToOneWithError()
	assert.NoError(t, err)
	assert.Equal(t, 1, value)
}

func TestMax(t *testing.T) {
	value, err := FromIntArray([]int{4, 5, 4, 3, 2, 1, 2}).Max().ToOneWithError()
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
	a := FromIntArray([]int{1, 2, 3, 4, 5}).Reduce(0, func(a int, b int) int { return a + b }).ToOne()
	assert.Equal(t, 15, a)
}

func TestDo(t *testing.T) {
	b := []int{}
	a := FromIntArray([]int{1, 2, 3, 4, 5}).Do(func(v int) {
		b = append(b, v)
	}).ToArray()
	assert.Equal(t, b, a)
}

func TestDoOnError(t *testing.T) {
	var oerr error
	_, err := FromIntError(errors.New("dead")).DoOnError(func(err error) { oerr = err }).ToArrayWithError()
	assert.Error(t, err)
	assert.Equal(t, err, oerr)
}

func TestDoOnComplete(t *testing.T) {
	complete := false
	a, err := FromIntComplete().DoOnComplete(func() { complete = true }).ToArrayWithError()
	assert.NoError(t, err)
	assert.Equal(t, []int{}, a)
	assert.True(t, complete)
}
