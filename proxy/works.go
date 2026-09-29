package proxy

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"os"
	"sync"
	"time"
	"webcacher2/cache"
	"webcacher2/config"
	"webcacher2/queue"
)

func Work(obj *queue.QueueObj) {
	queue.GQueue.Running.Add(1)
	defer queue.GQueue.Running.Add(-1)

	if obj == nil {
		time.Sleep(time.Second * 10)
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
	if resp == nil || resp.StatusCode != 200 {
		queue.GQueue.Push(obj)
	}
}

func Update(obj *queue.QueueObj) {
	queue.GQueue.Running.Add(1)
	defer queue.GQueue.Running.Add(-1)
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
	if resp == nil || resp.StatusCode != 200 {
		queue.GQueue.Push(obj)
	}
}

func MainWork() {
	go SaveAll()
	for {
		for queue.GQueue.Running.Load() < int32(queue.GQueue.Workers) {
			time.Sleep(time.Millisecond * 300)
			if config.Global.NoQueue {
				break
			}

			obj1 := queue.GQueue.Pop()

			//obj2 := queue.UGQueue.Pop()

			go Work(obj1)
			//go Update(obj2)
		}
		time.Sleep(time.Second)
		if queue.GQueue.Running.Load() != 0 {
			time.Sleep(time.Second * 9)
		}
	}
}

func SaveAll() {
	for {
		if time.Now().Second() != 0 {
			time.Sleep(time.Second)
			continue
		}
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
}
