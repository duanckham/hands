package hands

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimpleCase(t *testing.T) {
	n := 0
	controller := New()

	controller.Do(func(ctx context.Context) error {
		n++
		return nil
	})

	controller.Run()

	assert.Equal(t, n, 1)
}

func TestFastestOnly(t *testing.T) {
	var n int32
	controller := New()

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(10) * time.Millisecond)
		atomic.AddInt32(&n, 1)
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(10) * time.Millisecond)
		atomic.AddInt32(&n, 1)
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		atomic.AddInt32(&n, 5)
		return nil
	})

	controller.Run(Fastest())

	assert.Equal(t, atomic.LoadInt32(&n), int32(5))
}

func TestFastestAndRunAll(t *testing.T) {
	var n int32
	controller := New()

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(10) * time.Millisecond)
		atomic.AddInt32(&n, 1)
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(10) * time.Millisecond)
		atomic.AddInt32(&n, 1)
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		atomic.AddInt32(&n, 5)
		return nil
	})

	controller.RunAll(Fastest())

	assert.Equal(t, atomic.LoadInt32(&n), int32(5))
	time.Sleep(time.Duration(100) * time.Millisecond)
	assert.Equal(t, atomic.LoadInt32(&n), int32(7))
}

func TestPercentageOnly(t *testing.T) {
	var n int32
	controller := New()

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(10) * time.Millisecond)
		atomic.AddInt32(&n, 1)
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(20) * time.Millisecond)
		atomic.AddInt32(&n, 1)
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(30) * time.Millisecond)
		atomic.AddInt32(&n, 1)
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(40) * time.Millisecond)
		atomic.AddInt32(&n, 1)
		return nil
	})

	controller.Run(Percentage(0.5))

	assert.Equal(t, atomic.LoadInt32(&n), int32(2))
}

func TestPercentageAndRunAll(t *testing.T) {
	var n int32
	controller := New()

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(10) * time.Millisecond)
		atomic.AddInt32(&n, 1)
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(20) * time.Millisecond)
		atomic.AddInt32(&n, 1)
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(30) * time.Millisecond)
		atomic.AddInt32(&n, 1)
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(40) * time.Millisecond)
		atomic.AddInt32(&n, 1)
		return nil
	})

	controller.Run(Percentage(0.5))

	assert.Equal(t, atomic.LoadInt32(&n), int32(2))
	time.Sleep(time.Duration(100) * time.Millisecond)
	assert.Equal(t, atomic.LoadInt32(&n), int32(4))
}

func TestBetweenOnly(t *testing.T) {
	var n int32
	controller := New()

	for i := 0; i < 5; i++ {
		controller.Do(func(ctx context.Context) error {
			atomic.AddInt32(&n, 1)
			return nil
		}, P(int32(i)))
	}

	controller.Run(Between(2, 3))

	assert.Equal(t, atomic.LoadInt32(&n), int32(2))
}

func TestBetweenAndRunAll(t *testing.T) {
	var n int32
	controller := New()

	for i := 0; i < 5; i++ {
		controller.Do(func(ctx context.Context) error {
			atomic.AddInt32(&n, 1)
			return nil
		}, P(int32(i)))
	}

	controller.RunAll(Between(2, 3))

	assert.Equal(t, atomic.LoadInt32(&n), int32(2))
	time.Sleep(time.Duration(10) * time.Millisecond)
	assert.Equal(t, atomic.LoadInt32(&n), int32(5))
}

func TestInOnly(t *testing.T) {
	var n int32
	controller := New()

	for i := 0; i < 5; i++ {
		controller.Do(func(ctx context.Context) error {
			atomic.AddInt32(&n, 1)
			return nil
		}, P(int32(i)))
	}

	controller.Run(In([]int32{2, 4}))

	assert.Equal(t, atomic.LoadInt32(&n), int32(2))
}

func TestInAndRunAll(t *testing.T) {
	var n int32
	controller := New()

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(10) * time.Millisecond)
		atomic.AddInt32(&n, 2)
		return nil
	}, P(1))

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(20) * time.Millisecond)
		atomic.AddInt32(&n, 3)
		return nil
	}, P(2))

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(30) * time.Millisecond)
		atomic.AddInt32(&n, 4)
		return nil
	}, P(3))

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(40) * time.Millisecond)
		atomic.AddInt32(&n, 5)
		return nil
	}, P(4))

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(50) * time.Millisecond)
		atomic.AddInt32(&n, 6)
		return nil
	}, P(5))

	controller.RunAll(In([]int32{2, 4}))

	assert.Equal(t, atomic.LoadInt32(&n), int32(8))
	time.Sleep(time.Duration(100) * time.Millisecond)
	assert.Equal(t, atomic.LoadInt32(&n), int32(20))
}

func TestWithContextButTimeout(t *testing.T) {
	c, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	var n int32
	controller := New()

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(100) * time.Millisecond)
		atomic.AddInt32(&n, 1)
		return nil
	}, Name("sleep 100ms"))

	err := controller.Run(WithContext(c))

	assert.Equal(t, err.Error(), "context deadline exceeded")
	assert.Equal(t, atomic.LoadInt32(&n), int32(0))
}

func TestDone(t *testing.T) {
	var n int32

	controller := New()

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(10) * time.Millisecond)
		atomic.AddInt32(&n, 1)
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(10) * time.Millisecond)
		atomic.AddInt32(&n, 1)
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		atomic.AddInt32(&n, 5)
		return nil
	})

	controller.Done(func() {
		assert.Equal(t, atomic.LoadInt32(&n), int32(7))
	})

	controller.RunAll(Fastest())

	assert.Equal(t, atomic.LoadInt32(&n), int32(5))
	time.Sleep(time.Duration(500) * time.Millisecond)
}
