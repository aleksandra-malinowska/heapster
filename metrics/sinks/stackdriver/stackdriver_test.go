package stackdriver

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
