package proxy

import (
	"fmt"
	"os"
	"time"
)

func Log(a ...any) {
	data := []byte(fmt.Sprintln(time.Now().Local().String(), "LOG", a))
	fd, err := os.OpenFile("webcacher.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	fd.Write(data)
	fd.Close()
}
