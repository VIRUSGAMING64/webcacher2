package proxy

import (
	"fmt"
	"os"
	"time"
)

const file string = "/tmp/webcacher.log"

func Write(mot string, a ...any) {
	data := []byte(fmt.Sprintln(time.Now().Local().String(), mot, a))
	fd, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	fd.Write(data)
	fd.Close()
}

func Error(a ...any) {
	Write("ERROR", a)
}

func Log(a ...any) {
	Write("LOG", a)
}
