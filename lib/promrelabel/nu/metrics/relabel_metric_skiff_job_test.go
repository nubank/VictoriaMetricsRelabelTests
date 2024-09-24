package promrelabel

import (
	"testing"
)

func skiff_get_relabel_config() string {
  return `
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
    regex: .+\-RANDOM\-.+;.*TOKENS?\-REVOKED(\-CUSTOMER-FACING)?$
    source_labels: [group, topic]
  - action: drop
    regex: http.*
    source_labels: [endpoint]
  - source_labels: [topic]
    target_label: topic
  - action: drop
    regex: ^.+$
    source_labels: [__meta_kubernetes_pod_label_istio]`
}

func TestRelabelSkiffJob(t *testing.T) {
  MetricsTest(t,

  skiff_get_relabel_config(),

  `services_finagle_http_endpoint_logical_failures{endpoint="/api/render/pagar/:id",group="group-RANDOM-1234",topic="foo"}`,

  true,

  `services_finagle_http_endpoint_logical_failures{endpoint="/api/render/pagar/",group="group-RANDOM-1234",topic="foo"}`)
}

func TestRelabelSkiffJobBoleto(t *testing.T) {
  MetricsTest(t,

  skiff_get_relabel_config(),

  `services_finagle_http_endpoint_logical_failures{endpoint="/api/render/boleto-cobranca/asdf",group="group-RANDOM-1234",topic="foo"}`,

  true,

  `services_finagle_http_endpoint_logical_failures{endpoint="/api/render/boleto-cobranca/",group="group-RANDOM-1234",topic="foo"}`)
}

func TestRelabelSkiffJobDropRandom(t *testing.T) {
  MetricsTest(t,

  skiff_get_relabel_config(),

  `services_finagle_http_endpoint_logical_failures{group="group-RANDOM-1234",topic="asd-TOKENS-REVOKED-CUSTOMER-FACING"}`,

  true,

  `{}`)
}

func TestRelabelSkiffJobDropHttp(t *testing.T) {
  MetricsTest(t,

  skiff_get_relabel_config(),

  `foo{endpoint="http-asdf"}`,

  true,

  `{}`)
}

func TestRelabelSkiffJobDropK8sPod(t *testing.T) {
  MetricsTest(t,

  skiff_get_relabel_config(),

  `foo{__meta_kubernetes_pod_label_istio="aaa"}`,

  true,

  `{}`)
}