package promrelabel

import (
	"testing"
)


func TestRelabelOpsHealthTrigger(t *testing.T) {
  MetricsTest(t,

  `
    - action: drop
      regex: ^.+$
      source_labels: [instance]
  `,

  `metric_name{instance="test_instance"}`,

  true,

  `{}`)
}

func TestRelabelOpsHealthTriggerDontDrop(t *testing.T) {
  MetricsTest(t,

  `
    - action: drop
      regex: ^.+$
      source_labels: [instance]
  `,

  `metric_name{foo="bar"}`,

  true,

  `metric_name{foo="bar"}`)
}