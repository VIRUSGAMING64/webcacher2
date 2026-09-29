package views

import (
	"runtime"
	"strconv"
	"strings"
)

func RenderServer() string {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	rows := []string{

		metric("CPUs", strconv.Itoa(runtime.NumCPU())),
		metric("Memory used", formatBytes(int64(m.Alloc+m.HeapAlloc))),
	}
	return PanelStyle.Render(strings.Join(rows, "\n"))
}
