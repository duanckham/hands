# Hands

#### “Dedicated to Brother Chang”

Hands is a process controller used to control the execution and return strategies of multiple goroutines.

## Getting started

### A simple example

```
n := 0
controller := hands.New()

controller.Do(func(ctx context.Context) error {
  n++
  return nil
})

err := controller.Run()
if err != nil {
  // ...
}

fmt.Println(n)

# Output:
# 1
```

Use the `Do` method to add a task, use the `Run` method to start the task(s).

## TaskOption

`TaskOption` is used to set some metadata for the task.

### `func Priority(priority int32) TaskOption`

Use the `Priority` method to set a priority for a task. The higher the priority, the higher the execution order.

(`hands.P()` is an alias for `hands.Priority()`.)

```
controller := New()

controller.Do(func(ctx context.Context) error {
  fmt.Println("3")
  return nil
}, hands.P(1))

controller.Do(func(ctx context.Context) error {
  fmt.Println("2")
  return nil
}, hands.P(2))

controller.Do(func(ctx context.Context) error {
  fmt.Println("1")
  return nil
}, hands.P(3))

controller.Run()

// Output:
// 1
// 2
// 3
```

## HandOption

HandOption is used to control the execution strategy of the task.

### `func Fastest() HandOption`

`Fastest()`: When a task is completed, return immediately.

```
n := 0
controller := hands.New()

controller.Do(func(ctx context.Context) error {
  time.Sleep(time.Duration(10) * time.Millisecond)
  n += 1
  return nil
})

controller.Do(func(ctx context.Context) error {
  n += 2
  return nil
})

controller.Run(hands.Fastest())

fmt.Println(n)

# Output:
# 2
```

### `func Percentage(percentage float32) HandOption`

When a certain percentage of tasks are executed, the results are returned.
```
n := 0
controller := hands.New()

controller.Do(func(ctx context.Context) error {
  n++
  return nil
})

controller.Do(func(ctx context.Context) error {
  n++
  return nil
})

controller.Do(func(ctx context.Context) error {
  n++
  return nil
})

controller.Do(func(ctx context.Context) error {
  n++
  return nil
})

controller.Run(hands.Percentage(0.5))

fmt.Println(n)

# Output:
# 2
```

### `func Between(l, r int32) HandOption`

`Between()`: Only execute tasks with a priority within the specified range.

```
n := 0
controller := hands.New()

controller.Do(func(ctx context.Context) error {
  n += 1
  return nil
}, hands.P(1))

controller.Do(func(ctx context.Context) error {
  n += 2
  return nil
}, hands.P(2))

controller.Do(func(ctx context.Context) error {
  n += 3
  return nil
}, hands.P(3))

controller.Do(func(ctx context.Context) error {
  n += 4
  return nil
}, hands.P(4))

controller.Run(hands.Between(2, 3))

fmt.Println(n)

# Output:
# 5
```

*Note*: If the use the `controller.Run()` method, tasks outside the `Between()` will not be executed, you can use the `controller.RunAll()` method to allow other priority tasks to be executed asynchronously.

```
...
controller.RunAll(hands.Between(2, 3))

fmt.Println(n)
time.Sleep(time.Duration(10) * time.Millisecond)
fmt.Println(n)

# Output:
# 5
# 10
```

### `func In(in []int32) HandOption`

`In()`: Only execute tasks in the specified priority list.

```
n := 0
controller := hands.New()

controller.Do(func(ctx context.Context) error {
  n += 1
  return nil
}, hands.P(1))

controller.Do(func(ctx context.Context) error {
  n += 2
  return nil
}, hands.P(2))

controller.Do(func(ctx context.Context) error {
  n += 3
  return nil
}, hands.P(3))

controller.Do(func(ctx context.Context) error {
  n += 4
  return nil
}, hands.P(4))

controller.Run(hands.In([]int32{2, 4}))

fmt.Println(n)

# Output:
# 6
```

Yes, the `controller.RunAll()` method can also be used here.

```
...
controller.RunAll(hands.In([]int32{2, 4}))

fmt.Println(n)
time.Sleep(time.Duration(10) * time.Millisecond)
fmt.Println(n)

# Output:
# 6
# 10
```

### `func WithContext(ctx context.Context) HandOption`

Make the task use the specified context.

```
c, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
defer cancel()

controller := hands.New()

controller.Do(func(ctx context.Context) error {
  time.Sleep(time.Duration(100) * time.Millisecond)
  return nil
})

err := controller.Run(hands.WithContext(c))

fmt.Println(err.Error())

# Output:
# context deadline exceeded
```

## Callback after all tasks have been executed

```
n := 0
controller := hands.New()

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

// Here.
controller.Done(func() {
  // `True`
  assert.Equal(t, n, 7)
})

controller.RunAll(Fastest())

// `True`
assert.Equal(t, n, 5)
time.Sleep(time.Duration(500) * time.Millisecond)
```
