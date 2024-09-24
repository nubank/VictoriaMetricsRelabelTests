package promrelabel

import (
	"testing"
)


func TestRelabelKubernetesKubestatemetricsOld(t *testing.T) {
  MetricsTest(t,

  `
    - action: keep
      regex: 
        (kube_pod_container_resource_requests_cpu_cores|kube_pod_container_resource_limits_cpu_cores|kube_pod_container_resource_requests_memory_bytes|kube_pod_container_resource_limits_memory_bytes|kube_node_status_capacity_pods|kube_node_status_capacity_cpu_cores|kube_node_status_capacity_memory_bytes|kube_node_status_allocatable_pods|kube_node_status_allocatable_cpu_cores|kube_node_status_allocatable_memory_bytes|kube_hpa_.*)
      source_labels: [__name__]
  `,

  `kube_pod_container_resource_requests_cpu_cores{,foo="foo",bar="bar"}`,

  true,

  `kube_pod_container_resource_requests_cpu_cores{foo="foo",bar="bar"}`)
}