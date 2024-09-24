package promrelabel

import (
	"testing"
)

func TestRelabelAuroraNodeExporterDrop(t *testing.T) {
	MetricsTest(t,

		`
    - action: drop
      regex: ^go_.+$
      source_labels: [__name__]
  `,

		`go_gc_duration_seconds{quantile="0"}`,

		true,

		`{}`)
}

func TestRelabelAuroraNodeExporterDontDrop(t *testing.T) {
	MetricsTest(t,

		`
    - action: drop
      regex: ^go_.+$
      source_labels: [__name__]
  `,

		`node_disk_io_now{device="dm-0"}`,

		true,

		`node_disk_io_now{device="dm-0"}`)
}
