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
)

func TestTranslateUptime(t *testing.T) {
	metricValue := core.MetricValue{
		ValueType: core.ValueInt64,
		IntValue:  30000,
	}
	labels := map[string]string{}
	name := "uptime"
	timestamp := time.Now()
	createTime := timestamp

	ts := sink.TranslateMetric(timestamp, labels, name, metricValue, createTime)

	as := assert.New(t)

	as.Equal(ts.Metric.Type, "container.googleapis.com/container/uptime")
	as.Equal(len(ts.Points), 1)
	as.Equal(ts.Points[0].Value.DoubleValue, 30.0)
}
