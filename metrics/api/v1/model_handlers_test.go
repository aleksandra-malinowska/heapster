// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"k8s.io/heapster/metrics/core"
	metricsink "k8s.io/heapster/metrics/sinks/metric"

	restful "github.com/emicklei/go-restful"
	"github.com/stretchr/testify/assert"
)

type FakeMetricSink struct {
	*metricsink.MetricSink
}

func (f *FakeMetricSink) GetContainersForPodFromNamespace(namespace, pod string) []string {
	return []string{"some container", "another con"}
}

func NewTestApi() *Api {
	gkeMetrics := make(map[string]core.MetricDescriptor)
	gkeLabels := make(map[string]core.LabelDescriptor)

	fakeMetricSink := &FakeMetricSink{&metricsink.MetricSink{}}

	return &Api{
		runningInKubernetes: true,
		metricSink:          fakeMetricSink,
		historicalSource:    nil,
		gkeMetrics:          gkeMetrics,
		gkeLabels:           gkeLabels,
	}
}

func TestAddClusterMetricsRoutes(t *testing.T) {
	fakeApi := NewTestApi()

	tests := []struct {
		name          string
		fun           func(request *restful.Request, response *restful.Response)
		pathParams    map[string]string
		expectedNames []string
	}{
		{
			name:          "get pod container list",
			fun:           fakeApi.podContainerList,
			pathParams:    map[string]string{},
			expectedNames: []string{"some container", "another con"},
		},
	}

	assert := assert.New(t)
	restful.DefaultResponseMimeType = restful.MIME_JSON

	for _, test := range tests {
		req := restful.NewRequest(&http.Request{})
		pathParams := req.PathParameters()
		for k, v := range test.pathParams {
			pathParams[k] = v
		}
		recorder := &fakeRespRecorder{
			data:    new(bytes.Buffer),
			headers: make(http.Header),
		}
		resp := restful.NewResponse(recorder)

		test.fun(req, resp)

		actualNames := []string{}
		err := json.Unmarshal(recorder.data.Bytes(), &actualNames)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		assert.Equal(http.StatusOK, recorder.status, "status should have been OK (200)")
		assert.Equal(test.expectedNames, actualNames, "should have gotten expected JSON")
	}
}
