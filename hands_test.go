package hands

import (
	"context"
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
	n := 0
	controller := New()

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(10) * time.Millisecond)
		n++
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(10) * time.Millisecond)
		n++
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		n += 5
		return nil
	})

	controller.Run(Fastest())

	assert.Equal(t, n, 5)
}

func TestFastestAndRunAll(t *testing.T) {
	n := 0
	controller := New()

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(10) * time.Millisecond)
		n++
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(10) * time.Millisecond)
		n++
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		n += 5
		return nil
	})

	controller.RunAll(Fastest())

	assert.Equal(t, n, 5)
	time.Sleep(time.Duration(100) * time.Millisecond)
	assert.Equal(t, n, 7)
}

func TestBetweenOnly(t *testing.T) {
	n := 0
	controller := New()

	for i := 0; i < 5; i++ {
		controller.Do(func(ctx context.Context) error {
			n++
			return nil
		}, P(int32(i)))
	}

	controller.Run(Between(2, 3))

	assert.Equal(t, n, 2)
}

func TestBetweenAndRunAll(t *testing.T) {
	n := 0
	controller := New()

	for i := 0; i < 5; i++ {
		controller.Do(func(ctx context.Context) error {
			n++
			return nil
		}, P(int32(i)))
	}

	controller.RunAll(Between(2, 3))

	assert.Equal(t, n, 2)
	time.Sleep(time.Duration(10) * time.Millisecond)
	assert.Equal(t, n, 5)
}

func TestInOnly(t *testing.T) {
	n := 0
	controller := New()

	for i := 0; i < 5; i++ {
		controller.Do(func(ctx context.Context) error {
			n++
			return nil
		}, P(int32(i)))
	}

	controller.Run(In([]int32{2, 4}))

	assert.Equal(t, n, 2)
}

func TestInAndRunAll(t *testing.T) {
	n := 0
	controller := New()

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(10) * time.Millisecond)
		n += 2
		return nil
	}, P(1))

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(20) * time.Millisecond)
		n += 3
		return nil
	}, P(2))

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(30) * time.Millisecond)
		n += 4
		return nil
	}, P(3))

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(40) * time.Millisecond)
		n += 5
		return nil
	}, P(4))

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(50) * time.Millisecond)
		n += 6
		return nil
	}, P(5))

	controller.RunAll(In([]int32{2, 4}))

	assert.Equal(t, n, 8)
	time.Sleep(time.Duration(100) * time.Millisecond)
	assert.Equal(t, n, 20)
}

func TestPercentageOnly(t *testing.T) {
	n := 0
	controller := New()

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(10) * time.Millisecond)
		n++
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(20) * time.Millisecond)
		n++
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(30) * time.Millisecond)
		n++
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(40) * time.Millisecond)
		n++
		return nil
	})

	controller.Run(Percentage(0.5))

	assert.Equal(t, n, 2)
}

func TestPercentageAndRunAll(t *testing.T) {
	n := 0
	controller := New()

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(10) * time.Millisecond)
		n++
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(20) * time.Millisecond)
		n++
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(30) * time.Millisecond)
		n++
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(40) * time.Millisecond)
		n++
		return nil
	})

	controller.Run(Percentage(0.5))

	assert.Equal(t, n, 2)
	time.Sleep(time.Duration(100) * time.Millisecond)
	assert.Equal(t, n, 4)
}

func TestWithContextButTimeout(t *testing.T) {
	c, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	n := 0
	controller := New()

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(100) * time.Millisecond)
		n++
		return nil
	}, Name("sleep 100ms"))

	err := controller.Run(WithContext(c))

	assert.Equal(t, err.Error(), "context deadline exceeded")
	assert.Equal(t, n, 0)
}

func TestDone(t *testing.T) {
	n := 0
	controller := New()

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(10) * time.Millisecond)
		n++
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		time.Sleep(time.Duration(10) * time.Millisecond)
		n++
		return nil
	})

	controller.Do(func(ctx context.Context) error {
		n += 5
		return nil
	})

	controller.Done(func() {
		assert.Equal(t, n, 7)
	})

	controller.RunAll(Fastest())

	assert.Equal(t, n, 5)
	time.Sleep(time.Duration(500) * time.Millisecond)
}
