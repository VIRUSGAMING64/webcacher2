package proxy

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
	"webcacher2/cache"
	"webcacher2/config"
	"webcacher2/queue"
)

func Work(obj *queue.QueueObj) {
	defer queue.GQueue.Running.Add(-1)
	if obj == nil {
		return
	}
	// {#685, 14} Esto esta hecho con copilot
	rootCAs, err := x509.SystemCertPool()
	if err != nil || rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	if certificate, err := os.ReadFile("public/proxy-ca.crt"); err == nil {
		rootCAs.AppendCertsFromPEM(certificate)
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(mustParseURL("http://localhost:8092")),
			TLSClientConfig: &tls.Config{RootCAs: rootCAs, MinVersion: tls.VersionTLS12},
		},
		Timeout: 3 * time.Second,
	}

	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(obj.Request)))
	req, err = http.NewRequest(req.Method, obj.Url, req.Body)
	req.Header.Add("webcacher-queue", "true")
	resp, err := client.Do(req)
	obj.Count += 1
	if resp == nil {
		queue.GQueue.Push(obj)
	}
}

func Update(obj *queue.QueueObj) {
	defer queue.GQueue.Running.Add(-1)
	if obj == nil {
		return
	}
	// {#685, 15} Esto esta hecho con copilot
	rootCAs, err := x509.SystemCertPool()
	if err != nil || rootCAs == nil {
		rootCAs = x509.NewCertPool()
	}
	if certificate, err := os.ReadFile("public/proxy-ca.crt"); err == nil {
		rootCAs.AppendCertsFromPEM(certificate)
	}
	client := &http.Client{
		Transport: &http.Transport{
			Proxy:           http.ProxyURL(mustParseURL("http://localhost:8092")),
			TLSClientConfig: &tls.Config{RootCAs: rootCAs, MinVersion: tls.VersionTLS12},
		},
		Timeout: 3 * time.Second,
	}

	/*
		*Aqui se esta haciendo update a todo,
		TODO verificar que sitios hacerles update
	*/

	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(obj.Request)))
	req, err = http.NewRequest(req.Method, obj.Url, req.Body)
	req.Header.Add("webcacher-update", "true")
	resp, err := client.Do(req)
	obj.Count += 1
	if resp == nil {
		queue.GQueue.Push(obj)
	}
}

func MainWork() {
	for {
		o := 0
		for queue.GQueue.Running.Load() < int32(queue.GQueue.Workers) {
			if config.Global.NoQueue {
				break
			}
			obj1 := queue.GQueue.Pop()
			obj2 := queue.UGQueue.Pop()
			go Work(obj1)
			go Update(obj2)
			queue.GQueue.Running.Add(1)
			o += 1
			if o == 512 {
				break
			}
		}
		SaveAll()
		fmt.Println("All data saved queue length: [", queue.GQueue.Length(), "]")
		time.Sleep(time.Second * 60)
	}
}

func SaveAll() {
	wg := sync.WaitGroup{}
	wg.Go(func() {
		queue.GQueue.Save("queue.json")
	})
	wg.Go(func() {
		cache.Global.Save()
	})
	wg.Go(func() {
		Pstats.Save("stats.json")
	})
	wg.Wait()
}
