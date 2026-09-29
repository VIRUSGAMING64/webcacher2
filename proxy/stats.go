package proxy

import (
	"encoding/json"
	"math"
	"net/http"
	"os"
	"sync"
	"time"
)

type UrlStat struct {
	Url    string `json:"url"`
	Method string `json:"method"`
	Time   int64  `json:"time"`
}

type Stats struct {
	Total      int64     `json:"total"`
	CacheUse   int64     `json:"cache_used"`
	Downloaded int64     `json:"downloaded"`
	Hints      int       `json:"hints"`
	Bypass     int       `json:"bypass"`
	Length     int       `json:"hlength"`
	History    []UrlStat `json:"History"`
	mtx        sync.Mutex
}

var Pstats *Stats = NewStats()

func NewStats() *Stats {
	s := Stats{}
	return &s
}

func (s *Stats) Copy() Stats {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.Clean()
	ns := Stats{
		Total:      s.Total,
		CacheUse:   s.CacheUse,
		Downloaded: s.Downloaded,
		Hints:      s.Hints,
		Bypass:     s.Bypass,
		Length:     s.Length,
		History:    s.History,
	}
	return ns
}

func (s *Stats) Load(file string) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &s)
}

func (s *Stats) Clean() {
	his := make([]UrlStat, 0)
	for _, elem := range s.History {
		if time.Since(time.Unix(0, elem.Time)).Seconds() >= 60 {
			continue
		}
		his = append(his, elem)
	}
	s.History = his
}

func (s *Stats) Save(file string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.Clean()
	data, err := json.MarshalIndent(&s, "", "   ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, data, 0644)
}

func (s *Stats) AddHint(resp *http.Response) {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.Hints += 1
	size := int64(math.Max(0, float64(resp.ContentLength)))
	s.CacheUse += size
	st := UrlStat{
		Url:    resp.Request.URL.String(),
		Method: resp.Request.Method,
		Time:   time.Now().UnixNano(),
	}
	s.History = append(s.History, st)
	s.Length = len(s.History)
}

func (s *Stats) AddBypass(resp *http.Response) {

	s.mtx.Lock()
	defer s.mtx.Unlock()
	size := int64(math.Max(0, float64(resp.ContentLength)))
	s.Bypass += 1
	s.Total += size
	s.CacheUse += size
	s.Downloaded += size
	st := UrlStat{
		Url:    resp.Request.URL.String(),
		Method: resp.Request.Method,
		Time:   time.Now().UnixNano(),
	}
	s.History = append(s.History, st)
	s.Length = len(s.History)

}
