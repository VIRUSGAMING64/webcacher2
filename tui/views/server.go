package views

import (
	"runtime"
	"runtime/metrics"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

var processCPUState struct {
	sync.Mutex
	lastCPU  time.Duration
	lastWall time.Time
}

func RenderServer() string {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	rows := []string{
		metric("CPUs", strconv.Itoa(runtime.NumCPU())),
		metric("GOMAXPROCS", strconv.Itoa(runtime.GOMAXPROCS(0))),
		metric("Goroutines", strconv.Itoa(runtime.NumGoroutine())),
		metric("Memory allocated", formatBytes(int64(m.Alloc))),
		metric("Memory in use", formatBytes(int64(m.HeapInuse))),
		metric("Heap allocated", formatBytes(int64(m.HeapAlloc))),
		metric("Heap idle", formatBytes(int64(m.HeapIdle))),
		metric("Heap released", formatBytes(int64(m.HeapReleased))),
		metric("Heap system", formatBytes(int64(m.HeapSys))),
		metric("Stack in use", formatBytes(int64(m.StackInuse))),
		metric("Stack system", formatBytes(int64(m.StackSys))),
		metric("Runtime system", formatBytes(int64(m.Sys))),
		metric("Total allocated", formatBytes(int64(m.TotalAlloc))),
		metric("Allocations", strconv.FormatUint(m.Mallocs, 10)),
		metric("Frees", strconv.FormatUint(m.Frees, 10)),
		metric("Heap objects", strconv.FormatUint(m.HeapObjects, 10)),
		metric("GC cycles", strconv.FormatUint(uint64(m.NumGC), 10)),
		metric("Forced GC", strconv.FormatUint(uint64(m.NumForcedGC), 10)),
		metric("GC CPU", strconv.FormatFloat(m.GCCPUFraction*100, 'f', 2, 64)+"%"),
		metric("Last GC", formatLastGC(m.LastGC)),
		metric("Next GC target", formatBytes(int64(m.NextGC))),
		metric("User CPU", formatFloatMetric("/cpu/classes/user:cpu-seconds")),
		metric("GC CPU total", formatFloatMetric("/cpu/classes/gc/total:cpu-seconds")),
		metric("Scheduler goroutines", formatUintMetric("/sched/goroutines:goroutines")),
		metric("GC total cycles", formatUintMetric("/gc/cycles/total:gc-cycles")),
	}
	cpuPercent, cpuTime := processCPUUsage()
	rows = append(rows,
		metric("Process CPU", cpuPercent),
		metric("Process CPU time", cpuTime),
	)
	return PanelStyle.Render(strings.Join(rows, "\n"))
}

func processCPUUsage() (string, string) {
	var usage syscall.Rusage
	if err := syscall.Getrusage(syscall.RUSAGE_SELF, &usage); err != nil {
		return "n/a", "n/a"
	}

	cpuTime := time.Duration(usage.Utime.Sec)*time.Second +
		time.Duration(usage.Utime.Usec)*time.Microsecond +
		time.Duration(usage.Stime.Sec)*time.Second +
		time.Duration(usage.Stime.Usec)*time.Microsecond
	now := time.Now()

	processCPUState.Lock()
	defer processCPUState.Unlock()
	if processCPUState.lastWall.IsZero() {
		processCPUState.lastCPU = cpuTime
		processCPUState.lastWall = now
		return "warming up", formatDuration(cpuTime)
	}

	wallDelta := now.Sub(processCPUState.lastWall)
	cpuDelta := cpuTime - processCPUState.lastCPU
	processCPUState.lastCPU = cpuTime
	processCPUState.lastWall = now
	if wallDelta <= 0 || cpuDelta < 0 {
		return "n/a", formatDuration(cpuTime)
	}

	percent := float64(cpuDelta) / float64(wallDelta) * 100
	return strconv.FormatFloat(percent, 'f', 2, 64) + "%", formatDuration(cpuTime)
}

func formatDuration(value time.Duration) string {
	return value.Round(time.Millisecond).String()
}

func formatLastGC(value uint64) string {
	if value == 0 {
		return "never"
	}
	return time.Since(time.Unix(0, int64(value))).Round(time.Millisecond).String() + " ago"
}

func formatUintMetric(name string) string {
	samples := []metrics.Sample{{Name: name}}
	metrics.Read(samples)
	if samples[0].Value.Kind() != metrics.KindUint64 {
		return "n/a"
	}
	return strconv.FormatUint(samples[0].Value.Uint64(), 10)
}

func formatFloatMetric(name string) string {
	samples := []metrics.Sample{{Name: name}}
	metrics.Read(samples)
	if samples[0].Value.Kind() != metrics.KindFloat64 {
		return "n/a"
	}
	return strconv.FormatFloat(samples[0].Value.Float64(), 'f', 2, 64) + " s"
}
