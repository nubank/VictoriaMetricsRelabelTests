package promrelabel

import (
	"testing"
)


func TestRelabelKubernetesCadvisor(t *testing.T) {
  MetricsTest(t,

  `
    - regex: (.+)
      source_labels: [container]
      target_label: container_name
    - regex: (.+)
      source_labels: [pod]
      target_label: pod_name
    - regex: nu-(.+)
      replacement: $1
      source_labels: [container_name]
      target_label: service
    - regex: nu-datomic;staging-global-[a-zA-Z0-9]+-(.+)-deployment.*
      source_labels: [container, pod]
      target_label: service
  `,

  `container_cpu_usage_seconds_total{container="nu-datomic",pod="staging-global-asdf-datomic-deployment",cpu="total"}`,

  true,

  `container_cpu_usage_seconds_total{container="nu-datomic",container_name="nu-datomic",pod="staging-global-asdf-datomic-deployment",pod_name="staging-global-asdf-datomic-deployment",service="datomic",cpu="total"}`)
}