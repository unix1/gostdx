# gostdx

Just few extended standard library functions for Golang using generics.

Prerequsites:
* Go version 1.18 or greater due to generics.

# usage

## slices

<details>
  <summary>Expand for slices examples</summary>


```go
import "github.com/unix1/gostdx/slices"
```

### fold

Generic sequential fold:

```go
list := []int{1, 2, 3, 4, 5}
sumFunc := func(elem, sum int) int { return sum + elem }
sum := slices.Fold(sumFunc, 0, list)
fmt.Println("sum is:", sum) // sum is 15
```

### concurrent fold

Generic concurrent fold:

#### lock-free

```go
acc := int64(0)
concurrency := 5
list := []int64{1, 2, 3, 4, 5}
sumFunc := func(elem int64, acc *int64) *int64 {
    atomic.AddInt64(acc, elem)
    return acc
}
sum := slices.FoldC(sumFunc, &acc, list, concurrency)
fmt.Println("sum is:", *sum) // sum is 15
```

#### with locks

Folds a list of tuples to a map

```go
type tuple struct {
    v1 string
    v2 string
}

type mapAcc struct {
    sync.Mutex
    m map[string]string
}

acc = &mapAcc{m: map[string]string{}}
concurrency := 2
list := []tuple{{"k1", "v1"}, {"k2", "v2"}, {"k3", "v3"}}
F := func(e tuple, acc *mapAcc) *mapAcc {
    acc.Lock()
    defer acc.Unlock()
    acc.m[e.v1] = e.v2
    return acc
}
m := slices.FoldC(F, acc, list, concurrency)
fmt.Println("map is:", m.m) // map is: map[k1:v1 k2:v2 k3:v3]
```

</details>

## maps

<details>
  <summary>Expand for maps examples</summary>


```go
import "github.com/unix1/gostdx/maps"
```

### fold

Generic sequential fold:

```go
m := map[int]int{1: 10, 2: 20, 3: 30}
sumFunc := func(k int, v int, acc int) int { return acc + k*v }
sum := maps.Fold(sumFunc, 0, m)
fmt.Println("sum of k*v is", sum) // sum of k*v is 140
```

### concurrent fold

Generic concurrent fold:

```go
acc := int64(0)
concurrency := 3
m := map[int64]int64{1: 10, 2: 20, 3: 30}
sumFunc := func(k int64, v int64, acc *int64) *int64 {
    atomic.AddInt64(acc, k*v)
    return acc
}
sum := maps.FoldC(sumFunc, &acc, m, concurrency)
fmt.Println("sum of k*v is", *sum) // sum of k*v is 140
```

</details>
