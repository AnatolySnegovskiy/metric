package agent

import (
	"context"
	"fmt"
	"github.com/shirou/gopsutil/mem"
	"reflect"
	"runtime"
	"time"
)

var m runtime.MemStats
var runtimeEntityArray = []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse",
	"HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse",
	"MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse",
	"StackSys", "Sys", "TotalAlloc"}

func (a *Agent) updateStoragePeriodically(ctx context.Context) error {
	runtime.ReadMemStats(&m)

	gauge, err := a.storage.GetMetricType("gauge")
	if err != nil {
		return fmt.Errorf("error getting gauge: %w", err)
	}

	counter, err := a.storage.GetMetricType("counter")
	if err != nil {
		return fmt.Errorf("error getting counter: %w", err)
	}

	if counter.Process(ctx, "PollCount", "1") != nil {
		return fmt.Errorf("error while processing field: PollCount")
	}
	if gauge.Process(ctx, "RandomValue", fmt.Sprintf("%v", time.Now().UnixNano())) != nil {
		return fmt.Errorf("error while processing field: RandomValue")
	}

	for _, field := range runtimeEntityArray {
		if gauge.Process(ctx, field, fmt.Sprintf("%v", reflect.ValueOf(m).FieldByName(field))) != nil {
			return fmt.Errorf("error while processing field: %s", field)
		}
	}

	return nil
}
func (a *Agent) updateGopsutil(ctx context.Context) error {
	gauge, err := a.storage.GetMetricType("gauge")
	if err != nil {
		return fmt.Errorf("error getting gauge: %w", err)
	}
	v, _ := mem.VirtualMemory()

	if err := gauge.Process(ctx, "TotalMemory", fmt.Sprintf("%v", v.Total)); err != nil {
		return err
	}

	if err := gauge.Process(ctx, "FreeMemory", fmt.Sprintf("%v", v.Free)); err != nil {
		return err
	}

	if err := gauge.Process(ctx, "CPUutilization1", fmt.Sprintf("%v", v.UsedPercent)); err != nil {
		return err
	}

	return err
}
