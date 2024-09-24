package promrelabel

import (
	"testing"
)

func TestRelabelElasticacheInstancesLabels(t *testing.T) {
	MetricsTest(t,

		`
    - regex: (redis|memcached)_up
      replacement: $${1}_labels
      source_labels: [__name__]
      target_label: __name__
    - action: keep
      regex: (redis|memcached)_labels
      source_labels: [__name__]
  `,

		`redis_up{foo="foo",bar="bar"}`,

		true,

		`redis_labels{foo="foo",bar="bar"}`)
}
