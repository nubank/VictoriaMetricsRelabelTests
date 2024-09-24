package promrelabel

import (
	"testing"
)


func TestRelabelMetricsPlatformStatus(t *testing.T) {
  MetricsTest(t,

  `
    - action: keep
      regex: (prometheus_build_info|vm_app_version)
      source_labels: [__name__]
  `,

  `prometheus_build_info{foo="bar"}`,

  true,

  `prometheus_build_info{foo="bar"}`)
}