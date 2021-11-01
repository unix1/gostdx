# gostdx

Just few extended standard library functions for Golang using generics.

Prerequsites:
* Go version 1.18 or greater due to generics.

## lists

Import

```go
import "github.com/unix1/gostdx/lists"
```

### fold (aka reduce)

Generic sequential fold/reduce example:

```go
list := []int{1,2,3,4,5}
sumFunc := func(elem, sum int) int { return sum + elem }
sum := lists.Fold(sumFunc, 0, list)
fmt.Println("sum is:", sum) // sum is 15
```

### concurrent fold

Generic concurrent examples:

#### lock-free

```go
acc := int64(0)
concurrency := 5
list := []int64{1,2,3,4,5}
sumFunc := func(elem int64, acc *int64) *int64 {
    atomic.AddInt64(acc, elem)
    return acc
}
sum := lists.FoldC(sumFunc, &acc, list, concurrency)
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
m := lists.FoldC(F, acc, list, concurrency)
fmt.Println("map is:", m.m) // map is: map[k1:v1 k2:v2 k3:v3]
```
