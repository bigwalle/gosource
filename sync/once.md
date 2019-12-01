# go sync once 用法

> 控制函数只能被执行一次,不能被重复使用

```go
package main

import (
  "fmt"
  "sync"
  "time"
)

func main() {
  go do()

  go do()

  time.Sleep(time.Second*2)
}
var once  sync.Once
func do (){
  fmt.Println("start do ")
  once.Do(func() {
    fmt.Printf("doing something.")
  })
  fmt.Println("do end")
}


```

## once 源码
```go
package sync

import (
	"sync/atomic"
)

// Once is an object that will perform exactly one action.
//Once 是一个只执行一次的对象
type Once struct {
	// done indicates whether the action has been performed.
	// It is first in the struct because it is used in the hot path.
	// The hot path is inlined at every call site.
	// Placing done first allows more compact instructions on some architectures (amd64/x86),
	// and fewer instructions (to calculate offset) on other architectures.
	//定义标记,0 未调用,1 已调用
    done uint32 
    //锁标志
	m    Mutex  
}

// Do calls the function f if and only if Do is being called for the
// first time for this instance of Once. In other words, given
// 	var once Once
// if once.Do(f) is called multiple times, only the first call will invoke f,
// even if f has a different value in each invocation. A new instance of
// Once is required for each function to execute.
//
// Do is intended for initialization that must be run exactly once. Since f
// is niladic, it may be necessary to use a function literal to capture the
// arguments to a function to be invoked by Do:
// 	config.once.Do(func() { config.init(filename) })
//
// Because no call to Do returns until the one call to f returns, if f causes
// Do to be called, it will deadlock.
//
// If f panics, Do considers it to have returned; future calls of Do return
// without calling f.
//
//
func (o *Once) Do(f func()) {
	// Note: Here is an incorrect implementation of Do:
	//
	//	if atomic.CompareAndSwapUint32(&o.done, 0, 1) {
	//		f()
	//	}
	//
	// Do guarantees that when it returns, f has finished.
	// This implementation would not implement that guarantee:
	// given two simultaneous calls, the winner of the cas would
	// call f, and the second would return immediately, without
	// waiting for the first's call to f to complete.
	// This is why the slow path falls back to a mutex, and why
	// the atomic.StoreUint32 must be delayed until after f returns.
    //判断当前标志是否被执行过
	if atomic.LoadUint32(&o.done) == 0 {
		// Outlined slow-path to allow inlining of the fast-path.
		o.doSlow(f)
	}
}

func (o *Once) doSlow(f func()) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
        //函数f 调用之后,最后保存标志,标识函数已被调用
		defer atomic.StoreUint32(&o.done, 1)
		f()
	}
}

```
## 原理
> 相关的知识点
> mutex
Mutex 用于提供一种加锁机制（Locking Mechanism），可确保在某时刻只有一个协程在临界区运行，以防止出现竞态条件。

Mutex 可以在 sync 包内找到。Mutex 定义了两个方法：Lock 和 Unlock。所有在 Lock 和 Unlock 之间的代码，都只能由一个 Go 协程执行，于是就可以避免竞态条件。
```go
mutex.Lock()  
x = x + 1  
mutex.Unlock()
```

在上面的代码中，x = x + 1 只能由一个 Go 协程执行，因此避免了竞态条件。

如果有一个 Go 协程已经持有了锁（Lock），当其他协程试图获得该锁时，这些协程会被阻塞，直到 Mutex 解除锁定为止。
详见 [mutex](./mutex.md)

> atomic
>  原子操作,
Mutex由操作系统实现，而atomic包中的原子操作则由底层硬件直接提供支持。在 CPU 实现的指令集里，有一些指令被封装进了atomic包，这些指令在执行的过程中是不允许中断（interrupt）的，因此原子操作可以在lock-free的情况下保证并发安全，并且它的性能也能做到随 CPU 个数的增多而线性扩展。

>具体在 [atomic](/sync/atomic/atomic.md) 详细介绍


