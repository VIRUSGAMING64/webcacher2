package queue

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"
	"sync/atomic"
	wdebug "webcacher2/debug"
)

type QueueObj struct {
	Count   int    `json:"count"`
	Url     string `json:"url"`
	Host    string `json:"host"`
	Request []byte `json:"req"`
}

type Queue struct {
	first   *DataElem
	last    *DataElem
	Mtx     sync.Mutex
	Running atomic.Int32
	Size    int //* Solo para saberlo viendo el archivo
	Workers int
	Exists  map[string]bool
}

type DataElem struct {
	next *DataElem
	Data *QueueObj
	prev *DataElem
}

type SavedQueue struct {
	Size    int         `json:"size"` //* Solo para saberlo viendo el archivo
	Workers int         `json:"workers"`
	List    []*QueueObj `json:"list"`
}

var GQueue *Queue = NewQueue()

func NewQueue() *Queue {
	q := Queue{
		Workers: 32,
		Mtx:     sync.Mutex{},
		Running: atomic.Int32{},
		Exists:  make(map[string]bool),
	}
	q.first = nil
	q.last = nil
	return &q
}

func (q *Queue) Length() int {
	q.Mtx.Lock()
	defer q.Mtx.Unlock()
	return q.Size
}

func (q *Queue) Push(obj *QueueObj) {
	q.Mtx.Lock()
	defer q.Mtx.Unlock()
	if q.Exists[obj.Url] {
		return
	}
	q.Exists[obj.Url] = true
	elem := DataElem{}
	elem.Data = obj
	if q.first == nil {
		q.first = &elem
		q.last = &elem
	} else {
		tmp := q.last
		q.last = &elem
		tmp.next = &elem
		elem.prev = tmp
	}
	q.Size += 1
}

func (q *Queue) Pop() *QueueObj {
	q.Mtx.Lock()
	defer q.Mtx.Unlock()
	if q.first == nil {
		return nil
	}
	obj := q.first.Data
	q.first = q.first.next
	q.Size -= 1
	q.Exists[obj.Url] = false

	if q.first != nil {
		q.first.prev = nil
	}

	return obj
}

func (q *Queue) Json() []byte {
	q.Mtx.Lock()
	defer q.Mtx.Unlock()
	fr := q.first
	lis := []*QueueObj{}

	for fr != nil {
		lis = append(lis, fr.Data)
		fr = fr.next
	}
	filed := SavedQueue{Size: q.Size, Workers: q.Workers, List: lis}

	data, err := json.MarshalIndent(filed, "", "   ")

	if err != nil {
		return []byte{}
	}
	return data

}

func (q *Queue) Load(file string) {
	data, err := os.ReadFile(file)
	obj := SavedQueue{}
	if err != nil {
		wdebug.Error(" error loading queue: ", err)
		return
	}
	err = json.Unmarshal(data, &obj)
	for _, elem := range obj.List {
		q.Push(elem)
	}
	q.Workers = obj.Workers
	if err != nil {
		wdebug.Error(" error loading queue: ", err)
		return
	}

}

func (q *Queue) Save(file string) {
	data := q.Json()
	err := os.WriteFile(file, data, 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func NewObj(req *http.Request) *QueueObj {
	data, err := httputil.DumpRequest(req, true)

	if err != nil {
		return nil
	}
	o := QueueObj{
		Count:   0,
		Url:     req.URL.String(),
		Host:    req.Host,
		Request: data,
	}
	return &o
}
