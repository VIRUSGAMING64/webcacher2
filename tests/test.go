package main

import (
	"fmt"
	"webcacher2/cache"
	"webcacher2/proxy"
)

func TestCachefuntions() {
	ca := cache.NewUrlCache()
	ca.Load()
	data := []byte("morronga de caballo a la mierda de perrorerooqrqw")
	fmt.Println(ca.Push("hola", data))
	fmt.Println(ca.Push("h3la", data))
	fmt.Println(ca.Push("h2la", data))
	fmt.Println(ca.Push("hla", data))
	ca.Save()
	p := cache.NewUrlCache()
	p.Load()
	for k := range ca.MemKeys {
		data, err := p.Pop(ca.MemKeys[k])
		if err == nil {
			fmt.Println(ca.MemKeys[k], " : ", string(data))
		} else {
			fmt.Println(err)
		}
	}
	p.Save()
}

func main() {
	//TestCachefuntions()
	fmt.Println(proxy.HasInternet())
}
