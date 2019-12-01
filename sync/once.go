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

var  once1  sync.Once
func do (){
  fmt.Println("start do ")

  once1.Do(func() {
    fmt.Println("doing something.")
  })
  once.Do(func() {
    fmt.Println("doing something.")
  })
  fmt.Println("do end")
}
