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

// Test TranslateMetric

func testTranslateMetric(as *assert.Assertions, value int64, name string, labels map[string]string, expectedName string) *sd_api.TypedValue {
	metricValue := generateIntMetric(value)
	timestamp := time.Now()

	ts := sink.TranslateMetric(timestamp, labels, name, metricValue, timestamp)

	as.Equal(ts.Metric.Type, expectedName)
	as.Equal(len(ts.Points), 1)
	return ts.Points[0].Value
}

func TestTranslateUptime(t *testing.T) {
	as := assert.New(t)
	value := testTranslateMetric(as, 30000, "uptime", commonLabels,
		"container.googleapis.com/container/uptime")

	as.Equal(30.0, value.DoubleValue)
}

func TestTranslateCpuUsage(t *testing.T) {
	as := assert.New(t)
	value := testTranslateMetric(as, 3600000000000, "cpu/usage", commonLabels,
		"container.googleapis.com/container/cpu/usage_time")

	as.Equal(3600.0, value.DoubleValue)
}

func TestTranslateCpuLimit(t *testing.T) {
	as := assert.New(t)
	value := testTranslateMetric(as, 2000, "cpu/limit", commonLabels,
		"container.googleapis.com/container/cpu/reserved_cores")

	as.Equal(2.0, value.DoubleValue)
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
	as := assert.New(t)
	labels := deepCopy(commonLabels)
	labels["type"] = core.MetricSetTypePod
	value := testTranslateMetric(as, 2048, "memory/limit", labels,
		"container.googleapis.com/container/memory/bytes_total")

	as.Equal(int64(2048), value.Int64Value)
}

func TestTranslateMemoryNodeAllocatable(t *testing.T) {
	as := assert.New(t)
	value := testTranslateMetric(as, 2048, "memory/node_allocatable", commonLabels,
		"container.googleapis.com/container/memory/bytes_total")

	as.Equal(int64(2048), value.Int64Value)
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

func TestTranslateMemoryBytesUsed(t *testing.T) {
	as := assert.New(t)
	value := testTranslateMetric(as, 987, "memory/bytes_used", commonLabels,
		"container.googleapis.com/container/memory/bytes_used")

	as.Equal(int64(987), value.Int64Value)
}

// Test TranslateLabeledMetric

func TestTranslateFilesystemUsage(t *testing.T) {
	metric := core.LabeledMetric{
		MetricValue: generateIntMetric(10000),
		Labels: map[string]string{
			core.LabelResourceID.Key: "resource id",
		},
		Name: "filesystem/usage",
	}
	timestamp := time.Now()

	ts := sink.TranslateLabeledMetric(timestamp, commonLabels, metric, timestamp)

	as := assert.New(t)
	as.Equal(ts.Metric.Type, "container.googleapis.com/container/disk/bytes_used")
	as.Equal(len(ts.Points), 1)
	as.Equal(ts.Points[0].Value.Int64Value, int64(10000))
}

func TestTranslateFilesystemLimit(t *testing.T) {
	metric := core.LabeledMetric{
		MetricValue: generateIntMetric(30000),
		Labels: map[string]string{
			core.LabelResourceID.Key: "resource id",
		},
		Name: "filesystem/limit",
	}
	timestamp := time.Now()

	ts := sink.TranslateLabeledMetric(timestamp, commonLabels, metric, timestamp)

	as := assert.New(t)
	as.Equal(ts.Metric.Type, "container.googleapis.com/container/disk/bytes_total")
	as.Equal(len(ts.Points), 1)
	as.Equal(ts.Points[0].Value.Int64Value, int64(30000))
}
