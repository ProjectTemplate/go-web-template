package metric

import (
	"context"
	"runtime"

	"go-web-template/base/lib/config"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// RecordCpuInfo 记录 CPU 信息
func RecordCpuInfo(metricConfig config.Metric) error {
	meter := NewMeter(metricConfig.ServiceName)

	// CPU 总量
	cpuTotal, err := meter.Float64ObservableGauge("cpu_total", metric.WithDescription("the percentage of CPU total"))
	if err != nil {
		return err
	}

	// CPU 系统使用量
	cpuSystem, err := meter.Float64ObservableGauge("cpu_system", metric.WithDescription("the percentage of CPU system"))
	if err != nil {
		return err
	}

	// CPU 用户使用量
	cpuUser, err := meter.Float64ObservableGauge("cpu_user", metric.WithDescription("the percentage of CPU user"))
	if err != nil {
		return err
	}

	// CPU 空闲量
	cpuIdle, err := meter.Float64ObservableGauge("cpu_idle", metric.WithDescription("the percentage of CPU idle"))
	if err != nil {
		return err
	}

	// CPU 窃取量
	cpuSteal, err := meter.Float64ObservableGauge("cpu_steal", metric.WithDescription("the percentage of CPU steal"))
	if err != nil {
		return err
	}

	// CPU 窃取量
	cpuUsage, err := meter.Float64ObservableGauge("cpu_usage", metric.WithDescription("the percentage of CPU usage"))
	if err != nil {
		return err
	}

	obs := []metric.Observable{cpuTotal, cpuSystem, cpuUser, cpuIdle, cpuSteal, cpuUsage}
	_, err = meter.RegisterCallback(func(ctx context.Context, observer metric.Observer) error {
		times, errInner := cpu.Times(false)
		if errInner != nil {
			return errInner
		}
		percent, errInner := cpu.Percent(0, false)
		if errInner != nil {
			return errInner
		}
		stat := times[0]

		observer.ObserveFloat64(cpuTotal, stat.Total()/stat.Total()*100, initCommonAttributes(metricConfig))
		observer.ObserveFloat64(cpuSystem, stat.System/stat.Total()*100, initCommonAttributes(metricConfig))
		observer.ObserveFloat64(cpuUser, stat.User/stat.Total()*100, initCommonAttributes(metricConfig))
		observer.ObserveFloat64(cpuIdle, stat.Idle/stat.Total()*100, initCommonAttributes(metricConfig))
		observer.ObserveFloat64(cpuSteal, stat.Steal/stat.Total()*100, initCommonAttributes(metricConfig))
		observer.ObserveFloat64(cpuUsage, percent[0], initCommonAttributes(metricConfig))

		return nil
	}, obs...)

	return err
}

// RecordLoadInfo 记录负载信息
func RecordLoadInfo(metricConfig config.Metric) error {
	meter := NewMeter(metricConfig.ServiceName)

	// 总进程数量
	loadProcsTotal, err := meter.Int64ObservableGauge("load_procs_total", metric.WithDescription("the number of total procs"))
	if err != nil {
		return err
	}

	// 正在运行的进程数量
	loadProcsRunning, err := meter.Int64ObservableGauge("load_procs_running", metric.WithDescription("the number of running procs"))
	if err != nil {
		return err
	}

	// 被阻塞的进程数量
	loadProcsBlocked, err := meter.Int64ObservableGauge("load_procs_blocked", metric.WithDescription("the number of blocked procs"))
	if err != nil {
		return err
	}

	// 创建的进程数量
	loadProcsCreated, err := meter.Int64ObservableGauge("load_procs_created", metric.WithDescription("the number of created procs"))
	if err != nil {
		return err
	}

	obs := []metric.Observable{loadProcsTotal, loadProcsRunning, loadProcsBlocked, loadProcsCreated}

	_, err = meter.RegisterCallback(func(ctx context.Context, observer metric.Observer) error {
		misc, errInner := load.Misc()
		if errInner != nil {
			return err
		}

		observer.ObserveInt64(loadProcsTotal, int64(misc.ProcsTotal), initCommonAttributes(metricConfig))
		observer.ObserveInt64(loadProcsRunning, int64(misc.ProcsRunning), initCommonAttributes(metricConfig))
		observer.ObserveInt64(loadProcsBlocked, int64(misc.ProcsBlocked), initCommonAttributes(metricConfig))
		observer.ObserveInt64(loadProcsCreated, int64(misc.ProcsCreated), initCommonAttributes(metricConfig))

		return nil
	}, obs...)

	return err
}

// RecordMemoryInfo 记录内存信息
func RecordMemoryInfo(metricConfig config.Metric) error {
	meter := NewMeter(metricConfig.ServiceName)

	// 申请内存空间
	memoryTotal, err := meter.Int64ObservableGauge("memory_total", metric.WithDescription("Bytes of total heap memory"))
	if err != nil {
		return err
	}

	//占用内存空间
	memoryUsed, err := meter.Int64ObservableGauge("memory_used", metric.WithDescription("Bytes of used heap memory"))
	if err != nil {
		return err
	}

	obs := []metric.Observable{memoryTotal, memoryUsed}
	_, err = meter.RegisterCallback(func(ctx context.Context, observer metric.Observer) error {
		memory, errInner := mem.VirtualMemory()
		if errInner != nil {
			return errInner
		}

		observer.ObserveInt64(memoryTotal, int64(memory.Total), initCommonAttributes(metricConfig))
		observer.ObserveInt64(memoryUsed, int64(memory.Used), initCommonAttributes(metricConfig))

		return nil
	}, obs...)

	return err
}

// RecordGoRoutineCount 记录 goroutine 数量
func RecordGoRoutineCount(metricConfig config.Metric) error {
	meter := NewMeter(metricConfig.ServiceName)

	// goroutine 数量
	goroutineCount, err := meter.Int64ObservableGauge("goroutine_count", metric.WithDescription("Number of active goroutines"))
	if err != nil {
		return err
	}

	obs := []metric.Observable{goroutineCount}
	_, err = meter.RegisterCallback(func(ctx context.Context, observer metric.Observer) error {

		observer.ObserveInt64(goroutineCount, int64(runtime.NumGoroutine()), initCommonAttributes(metricConfig))

		return nil
	}, obs...)

	return err
}

func initCommonAttributes(metricConfig config.Metric) metric.MeasurementOption {
	return metric.WithAttributes(
		attribute.String("cluster", metricConfig.Cluster),
		attribute.String("service_name", metricConfig.ServiceName),
		attribute.String("pod", hostName),
	)
}
