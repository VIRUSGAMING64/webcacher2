package queue

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
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
	First   *DataElem
	Last    *DataElem
	Mtx     sync.Mutex
	Running atomic.Int32
	Size    int //* Solo para saberlo viendo el archivo
	Workers int
	Exists  map[string]bool
}

type DataElem struct {
	Next *DataElem
	Data *QueueObj
	Prev *DataElem
}

type SavedQueue struct {
	Size    int         `json:"size"` //* Solo para saberlo viendo el archivo
	Workers int         `json:"workers"`
	List    []*QueueObj `json:"list"`
}

var GQueue *Queue = NewQueue()
var UGQueue *Queue = NewQueue()

func NewQueue() *Queue {
	q := Queue{
		Workers: runtime.NumCPU(),
		Mtx:     sync.Mutex{},
		Running: atomic.Int32{},
		Exists:  make(map[string]bool),
	}
	q.First = nil
	q.Last = nil
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
	if q.First == nil {
		q.First = &elem
		q.Last = &elem
	} else {
		tmp := q.Last
		q.Last = &elem
		tmp.Next = &elem
		elem.Prev = tmp
	}
	q.Size += 1
}

func (q *Queue) Pop() *QueueObj {
	q.Mtx.Lock()
	defer q.Mtx.Unlock()
	if q.First == nil {
		return nil
	}
	obj := q.First.Data
	q.First = q.First.Next
	q.Size -= 1
	q.Exists[obj.Url] = false

	if q.First != nil {
		q.First.Prev = nil
	}

	return obj
}

func (q *Queue) Json() []byte {
	q.Mtx.Lock()
	defer q.Mtx.Unlock()
	fr := q.First
	lis := []*QueueObj{}

	for fr != nil {
		lis = append(lis, fr.Data)
		fr = fr.Next
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
