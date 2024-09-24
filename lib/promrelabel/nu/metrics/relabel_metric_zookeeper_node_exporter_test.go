package promrelabel

import (
	"testing"
)


func TestRelabelZookeeperNodeExporter(t *testing.T) {
  MetricsTest(t,

  `
    - action: drop
      regex: ^go_.+$
      source_labels: [__name__]
  `,

  `go_asd{foo="bar"}`,

  true,

  `{}`)
}