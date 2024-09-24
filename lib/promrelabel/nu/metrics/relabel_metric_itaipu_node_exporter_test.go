package promrelabel

import (
	"testing"
)


func TestRelabelItaipuNodeExporterDrop(t *testing.T) {
  MetricsTest(t,

  `
    - action: drop
      regex: ^go_.+$
      source_labels: [__name__]
  `,

  `go_gc_duration_seconds{quantile="1"}`,

  true,

  `{}`)
}

func TestRelabelItaipuNodeExporterDontDrop(t *testing.T) {
  MetricsTest(t,

  `
    - action: drop
      regex: ^go_.+$
      source_labels: [__name__]
  `,

  `node_cpu_seconds_total{cpu="0",mode="idle"}`,

  true,

  `node_cpu_seconds_total{cpu="0",mode="idle"}`)
}