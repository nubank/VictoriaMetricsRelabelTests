package promrelabel

import (
	"testing"
)


func TestRelabelKubernetesDatomicPods(t *testing.T) {
  MetricsTest(t,

  `
    - action: replace
      regex: ^services_finagle_http_endpoint_logical_.+;(.+api\/render\/pagar\/).+$
      replacement: $1
      source_labels: [__name__, endpoint]
      target_label: endpoint
    - action: replace
      regex: ^services_finagle_http_endpoint_logical_.+;(.+api\/render\/boleto-cobranca\/).+$
      replacement: $1
      source_labels: [__name__, endpoint]
      target_label: endpoint
    - action: drop
      regex: .+\-RANDOM\-.+;.*TOKENS?\-REVOKED$
      source_labels: [group, topic]
    - action: drop
      regex: http.*
      source_labels: [endpoint]
    - action: drop
      regex: ^go_.+$
      source_labels: [__name__]
    - source_labels: [topic]
      target_label: topic
  `,

  `services_finagle_http_endpoint_logical_failures{endpoint="/api/render/pagar/:id",group="group-RANDOM-1234",topic="foo"}`,

  true,

  `services_finagle_http_endpoint_logical_failures{endpoint="/api/render/pagar/",group="group-RANDOM-1234",topic="foo"}`)
}

func TestRelabelKubernetesDatomicPodsDropRevokeToken(t *testing.T) {
  MetricsTest(t,

  `
    - action: replace
      regex: ^services_finagle_http_endpoint_logical_.+;(.+api\/render\/pagar\/).+$
      replacement: $1
      source_labels: [__name__, endpoint]
      target_label: endpoint
    - action: replace
      regex: ^services_finagle_http_endpoint_logical_.+;(.+api\/render\/boleto-cobranca\/).+$
      replacement: $1
      source_labels: [__name__, endpoint]
      target_label: endpoint
    - action: drop
      regex: .+\-RANDOM\-.+;.*TOKENS?\-REVOKED$
      source_labels: [group, topic]
    - action: drop
      regex: http.*
      source_labels: [endpoint]
    - action: drop
      regex: ^go_.+$
      source_labels: [__name__]
    - source_labels: [topic]
      target_label: topic
  `,

  `services_finagle_http_endpoint_logical_failures{endpoint="/api/render/pagar/:id",group="group-RANDOM-1234",topic="foo-TOKENS-REVOKED"}`,

  true,

  `{}`)
}

func TestRelabelKubernetesDatomicPodsDropGo(t *testing.T) {
  MetricsTest(t,

  `
    - action: replace
      regex: ^services_finagle_http_endpoint_logical_.+;(.+api\/render\/pagar\/).+$
      replacement: $1
      source_labels: [__name__, endpoint]
      target_label: endpoint
    - action: replace
      regex: ^services_finagle_http_endpoint_logical_.+;(.+api\/render\/boleto-cobranca\/).+$
      replacement: $1
      source_labels: [__name__, endpoint]
      target_label: endpoint
    - action: drop
      regex: .+\-RANDOM\-.+;.*TOKENS?\-REVOKED$
      source_labels: [group, topic]
    - action: drop
      regex: http.*
      source_labels: [endpoint]
    - action: drop
      regex: ^go_.+$
      source_labels: [__name__]
    - source_labels: [topic]
      target_label: topic
  `,

  `go_gc_duration_seconds{quantile="0"}`,

  true,

  `{}`)
}