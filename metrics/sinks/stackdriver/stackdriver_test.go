package stackdriver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sd_api "google.golang.org/api/monitoring/v3"
	"k8s.io/heapster/metrics/core"
)

var (
	testProjectId = "test-project-id"
	zone          = "europe-west1-c"

	sink = &StackdriverSink{
		project:           testProjectId,
		zone:              zone,
		stackdriverClient: nil,
	}

	commonLabels = map[string]string{}
)

func generateIntMetric(value int64) core.MetricValue {
	return core.MetricValue{
		ValueType: core.ValueInt64,
		IntValue:  value,
	}
}

func generateFloatMetric(value float32) core.MetricValue {
	return core.MetricValue{
		ValueType:  core.ValueFloat,
		FloatValue: value,
	}
}

func deepCopy(source map[string]string) map[string]string {
	result := map[string]string{}
	for k, v := range source {
		result[k] = v
	}
	return result
}

func TestTranslateUptime(t *testing.T) {
	metricValue := generateIntMetric(30000)
	name := "uptime"
	timestamp := time.Now()

	ts := sink.TranslateMetric(timestamp, commonLabels, name, metricValue, timestamp)

	as := assert.New(t)
	as.Equal(ts.Metric.Type, "container.googleapis.com/container/uptime")
	as.Equal(len(ts.Points), 1)
	as.Equal(ts.Points[0].Value.DoubleValue, 30.0)
}

func TestTranslateCpuUsage(t *testing.T) {
	metricValue := generateIntMetric(3600000000000)
	name := "cpu/usage"
	timestamp := time.Now()
	createTime := timestamp.Add(-time.Second)

	ts := sink.TranslateMetric(timestamp, commonLabels, name, metricValue, createTime)

	as := assert.New(t)
	as.Equal(ts.Metric.Type, "container.googleapis.com/container/cpu/usage_time")
	as.Equal(len(ts.Points), 1)
	as.Equal(ts.Points[0].Value.DoubleValue, 3600.0)
}

func TestTranslateCpuLimit(t *testing.T) {
	metricValue := generateIntMetric(2000)
	name := "cpu/limit"
	timestamp := time.Now()
	createTime := timestamp.Add(-time.Second)

	ts := sink.TranslateMetric(timestamp, commonLabels, name, metricValue, createTime)

	as := assert.New(t)
	as.Equal(ts.Metric.Type, "container.googleapis.com/container/cpu/reserved_cores")
	as.Equal(len(ts.Points), 1)
	as.Equal(ts.Points[0].Value.DoubleValue, 2.0)
}

func TestTranslateMemoryLimitNode(t *testing.T) {
	metricValue := generateIntMetric(2048)
	name := "memory/limit"
	timestamp := time.Now()

	labels := deepCopy(commonLabels)
	labels["type"] = core.MetricSetTypeNode

	ts := sink.TranslateMetric(timestamp, labels, name, metricValue, timestamp)

	var expected *sd_api.TimeSeries = nil

	as := assert.New(t)
	as.Equal(ts, expected)
}

func TestTranslateMemoryLimitPod(t *testing.T) {
	metricValue := generateIntMetric(2048)
	name := "memory/limit"
	timestamp := time.Now()

	labels := deepCopy(commonLabels)
	labels["type"] = core.MetricSetTypePod

	ts := sink.TranslateMetric(timestamp, labels, name, metricValue, timestamp)

	as := assert.New(t)
	as.Equal(ts.Metric.Type, "container.googleapis.com/container/memory/bytes_total")
	as.Equal(len(ts.Points), 1)
	as.Equal(ts.Points[0].Value.Int64Value, int64(2048))
}

func TestTranslateMemoryNodeAllocatable(t *testing.T) {
	metricValue := generateIntMetric(2048)
	name := "memory/node_allocatable"
	timestamp := time.Now()

	ts := sink.TranslateMetric(timestamp, commonLabels, name, metricValue, timestamp)

	as := assert.New(t)
	as.Equal(ts.Metric.Type, "container.googleapis.com/container/memory/bytes_total")
	as.Equal(len(ts.Points), 1)
	as.Equal(ts.Points[0].Value.Int64Value, int64(2048))
}

func TestTranslateMemoryMajorPageFaults(t *testing.T) {
	metricValue := generateIntMetric(20)
	name := "memory/major_page_faults"
	timestamp := time.Now()

	ts := sink.TranslateMetric(timestamp, commonLabels, name, metricValue, timestamp)

	as := assert.New(t)
	as.Equal(ts.Metric.Type, "container.googleapis.com/container/memory/page_fault_count")
	as.Equal(len(ts.Points), 1)
	as.Equal(ts.Points[0].Value.Int64Value, int64(20))
	as.Equal(ts.Metric.Labels["fault_type"], "major")
}

func TestTranslateMemoryMinorPageFaults(t *testing.T) {
	metricValue := generateIntMetric(42)
	name := "memory/minor_page_faults"
	timestamp := time.Now()

	ts := sink.TranslateMetric(timestamp, commonLabels, name, metricValue, timestamp)

	as := assert.New(t)
	as.Equal(ts.Metric.Type, "container.googleapis.com/container/memory/page_fault_count")
	as.Equal(len(ts.Points), 1)
	as.Equal(ts.Points[0].Value.Int64Value, int64(42))
	as.Equal(ts.Metric.Labels["fault_type"], "minor")
}
