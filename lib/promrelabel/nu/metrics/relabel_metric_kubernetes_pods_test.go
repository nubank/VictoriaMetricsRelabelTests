package promrelabel

import (
	"testing"
)

func _getRule() string {
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
      regex: 
        ^beggar;(services_finagle_http_endpoint_logical_request_latency_ms|services_finagle_http_endpoint_logical_requests|services_finagle_http_endpoint_backups_backups_won|services_finagle_http_endpoint_logical_success|services_finagle_http_endpoint_backups_backups_sent);.+http.+$
      source_labels: [service, __name__, endpoint]
    - action: replace
      regex: ^auth;services_finagle_.+;nubank.okta.com:443:\S.+$
      replacement: nubank.okta.com
      source_labels: [service, __name__, endpoint]
      target_label: endpoint
    - action: replace
      regex: ^maat;services_finagle_.+;hooks.slack.com:443:\S.+$
      replacement: hooks.slack.com
      source_labels: [service, __name__, endpoint]
      target_label: endpoint
    - action: replace
      regex: ^webapp-proxy-webhooks;services_finagle_.+;(\S.+):([0-9]+):\S.+$
      replacement: $1
      source_labels: [service, __name__, endpoint]
      target_label: endpoint
    - action: drop
      regex: ^services_pathom_operation_latency_seconds_bucket;(7\.5|3\.0|30\.0|10\.0)$
      source_labels: [__name__, le]
    - action: drop
      regex: 
        ^pathom_api_call_latency_milliseconds_bucket;(7000\.0|9000\.0|13000\.0|17000\.0|20000.0)$
      source_labels: [__name__, le]
    - action: drop
      regex: ^services_finagle_http_endpoint_logical_request_latency_ms;0.9999?$
      source_labels: [__name__, quantile]
    - action: drop
      regex: ^take-it-easy;.+$
      source_labels: [service, nuinvest_service]
    - action: drop
      regex: .+\-RANDOM\-.+;.*TOKENS?\-REVOKED(\-CUSTOMER-FACING)?$
      source_labels: [group, topic]
    - action: drop
      regex: http.*
      source_labels: [endpoint]
    - action: drop
      regex: ^ops_.+;proximo$
      source_labels: [__name__, service]
    - action: drop
      regex: ^ops_.+;pushgateway$
      source_labels: [__name__, service]
    - action: drop
      regex: ^ops_.+;gauss$
      source_labels: [__name__, service]
    - source_labels: [topic]
      target_label: topic
    - action: drop
      regex: 
        ^(mobile_navigation_pop_received|mobile_navigation_push_received|mobile_navigation_replace_received)$
      source_labels: [__name__]
    - action: drop
      regex: ^mobile_scanner_performed$
      source_labels: [__name__]
    - action: drop
      regex: ^income_tax_report_sent_reports_total;mufasa$
      source_labels: [__name__, service]
    - action: drop
      regex: ^mobile_flow_loaded|mobile_flow_route_change_counter$
      source_labels: [__name__]
    - action: drop
      regex: ^services_http_requests_total;backoffice-proxy$
      source_labels: [__name__, service]
    - action: drop
      regex: ^.+$
      source_labels: [__meta_kubernetes_pod_label_istio]
`
}

func TestRelabelKubernetesPods(t *testing.T) {
  MetricsTest(t,

  _getRule(),

  `services_finagle_http_endpoint_logical_failures{endpoint="/api/render/pagar/:id",group="group-RANDOM-1234",topic="foo"}`,

  true,

  `services_finagle_http_endpoint_logical_failures{endpoint="/api/render/pagar/",group="group-RANDOM-1234",topic="foo"}`)
}

func TestRelabelKubernetesPodsBoletoCobranca(t *testing.T) {
  MetricsTest(t,

  _getRule(),

  `services_finagle_http_endpoint_logical_failures{endpoint="/api/render/boleto-cobranca/asdf",group="group-RANDOM-1234",topic="foo"}`,

  true,

  `services_finagle_http_endpoint_logical_failures{endpoint="/api/render/boleto-cobranca/",group="group-RANDOM-1234",topic="foo"}`)
}

func TestRelabelKubernetesPodsDropBeggar(t *testing.T) {
  MetricsTest(t,

  _getRule(),

  `services_finagle_http_endpoint_logical_request_latency_ms{service="beggar",endpoint="foo-http-bar"}`,

  true,

  `{}`)
}

func TestRelabelKubernetesPodsOkta(t *testing.T) {
  MetricsTest(t,

  _getRule(),

  `services_finagle_http_requests{endpoint="nubank.okta.com:443:/api/login",service="auth"}`,

  true,

  `services_finagle_http_requests{endpoint="nubank.okta.com",service="auth"}`)
}

func TestRelabelKubernetesPodsSlack(t *testing.T) {
  MetricsTest(t,

  _getRule(),

  `services_finagle_http_requests{endpoint="hooks.slack.com:443:/webhook/notify",service="maat"}`,

  true,

  `services_finagle_http_requests{endpoint="hooks.slack.com",service="maat"}`)
}

func TestRelabelKubernetesPodsWebApp(t *testing.T) {
  MetricsTest(t,

  _getRule(),

  `services_finagle_http_requests{endpoint="10.0.0.1:443:/api/callback",service="webapp-proxy-webhooks"}`,

  true,

  `services_finagle_http_requests{endpoint="10.0.0.1",service="webapp-proxy-webhooks"}`)
}

func TestRelabelKubernetesPodsDropBucketLatency(t *testing.T) {
  MetricsTest(t,

  _getRule(),

  `services_pathom_operation_latency_seconds_bucket{le="7.5"}`,

  true,

  `{}`)
}

func TestRelabelKubernetesPodsDropBucket2Latency(t *testing.T) {
  MetricsTest(t,

  _getRule(),

  `pathom_api_call_latency_milliseconds_bucket{le="7000.0"}`,

  true,

  `{}`)
}

func TestRelabelKubernetesPodsDropBucket3Latency(t *testing.T) {
  MetricsTest(t,

  _getRule(),

  `services_finagle_http_endpoint_logical_request_latency_ms{quantile="0.999"}`,

  true,

  `{}`)
}

func TestRelabelKubernetesPodsDropTakeItEasy(t *testing.T) {
  MetricsTest(t,

  _getRule(),

  `foo{service="take-it-easy",nuinvest_service="asdf"}`,

  true,

  `{}`)
}

func TestRelabelKubernetesPodsDropRandom(t *testing.T) {
  MetricsTest(t,

  _getRule(),

  `foo{group="aaa-RANDOM-bbb",topic="cccTOKEN-REVOKED"}`,

  true,

  `{}`)
}

func TestRelabelKubernetesPodsDropHttp(t *testing.T) {
  MetricsTest(t,

  _getRule(),

  `foo{endpoint="http://localhost"}`,

  true,

  `{}`)
}

func TestRelabelKubernetesPodsDropProximo(t *testing.T) {
  MetricsTest(t,

  _getRule(),

  `ops_foo_bar{service="proximo"}`,

  true,

  `{}`)
}

func TestRelabelKubernetesPodsDropPushgateway(t *testing.T) {
  MetricsTest(t,

  _getRule(),

  `ops_foo_bar{service="pushgateway"}`,

  true,

  `{}`)
}

func TestRelabelKubernetesPodsDropGauss(t *testing.T) {
  MetricsTest(t,

  _getRule(),

  `ops_foo_bar{service="gauss"}`,

  true,

  `{}`)
}

func TestRelabelKubernetesPodsRenameTopic(t *testing.T) {
  MetricsTest(t,

  _getRule(),

  `ops_foo_bar{topic="a"}`,

  true,

  `ops_foo_bar{topic="a"}`)
}

func TestRelabelKubernetesPodsDropMobile(t *testing.T) {
  MetricsTest(t,

  _getRule(),

  `mobile_navigation_push_received{foo="bar"}`,

  true,

  `{}`)
}

func TestRelabelKubernetesPodsDropMobileScanner(t *testing.T) {
  MetricsTest(t,

  _getRule(),

  `mobile_scanner_performed{foo="bar"}`,

  true,

  `{}`)
}

func TestRelabelKubernetesPodsDropIncomeTax(t *testing.T) {
  MetricsTest(t,

  _getRule(),

  `income_tax_report_sent_reports_total{service="mufasa"}`,

  true,

  `{}`)
}

func TestRelabelKubernetesPodsDropMobileFlow(t *testing.T) {
  MetricsTest(t,

  _getRule(),

  `mobile_flow_route_change_counter{foo="bar"}`,

  true,

  `{}`)
}

func TestRelabelKubernetesPodsDropHttpTotal(t *testing.T) {
  MetricsTest(t,

  _getRule(),

  `services_http_requests_total{service="backoffice-proxy"}`,

  true,

  `{}`)
}

func TestRelabelKubernetesPodsDropIstio(t *testing.T) {
  MetricsTest(t,

  _getRule(),

  `foo{__meta_kubernetes_pod_label_istio="backoffice-proxy"}`,

  true,

  `{}`)
}
