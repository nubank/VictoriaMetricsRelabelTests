package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelIstioSidecar(t *testing.T) {
  s := func(labelsStr string) []string {
    labelsStr = strings.Trim(labelsStr, "{}")
    labelsList := strings.Split(labelsStr, ",")
    for i := range labelsList {
      labelsList[i] = strings.TrimSpace(labelsList[i])
    }
    sort.Strings(labelsList)
    return labelsList
  }

  f := func(config, metric string, isFinalize bool, resultExpected string) {
    t.Helper()
    pcs, err := pr.ParseRelabelConfigsData([]byte(config))
    if err != nil {
      t.Fatalf("cannot parse %q: %s", config, err)
    }
    labels := promutils.MustNewLabelsFromString(metric)
    resultLabels := pcs.Apply(labels.GetLabels(), 0)
    if isFinalize {
      resultLabels = pr.FinalizeLabels(resultLabels[:0], resultLabels)
    }
    pr.SortLabels(resultLabels)
    result := pr.LabelsToString(resultLabels)

    sortedResult := s(result)
		sortedExpected := s(resultExpected)

		if !reflect.DeepEqual(sortedResult, sortedExpected) {
			t.Fatalf("unexpected result; got\n%v\nwant\n%v", sortedResult, sortedExpected)
		}
  }

  f(`
    - action: replace
      regex: ([^:]+):(\d+)
      replacement: $${2}
      source_labels: [__address__]
      target_label: __port__
    - action: keep
      regex: 15090
      source_labels: [__port__]
    - action: keep
      regex: true
      source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
    - action: drop
      regex: .+
      source_labels: [__meta_kubernetes_pod_label_istio]
    - action: replace
      source_labels: [__meta_kubernetes_pod_name]
      target_label: kubernetes_pod_name
    - action: replace
      source_labels: [__meta_kubernetes_pod_label_nubank_com_br_name]
      target_label: service
    - action: replace
      source_labels: [__meta_kubernetes_pod_label_nubank_com_br_tier]
      target_label: tier
    - {action: labeldrop, regex: __meta_kubernetes_pod_label_nubank_com_br_name}
    - {action: labelmap, regex: __meta_kubernetes_pod_label_nubank_com_br_(.+)}
  `,

  `{__address__="100.72.107.150:15090",__meta_kubernetes_namespace="default",__meta_kubernetes_pod_annotation_iam_amazonaws_com_role="prod/shore/pink/prod-pink-shore-role",__meta_kubernetes_pod_annotation_istio_io_rev="default",__meta_kubernetes_pod_annotation_kubectl_kubernetes_io_default_container="nu-shore",__meta_kubernetes_pod_annotation_kubectl_kubernetes_io_default_logs_container="nu-shore",__meta_kubernetes_pod_annotation_nubank_com_br_last_cycle="1726238197",__meta_kubernetes_pod_annotation_nubank_com_br_last_cycle_author="alan.lamb",__meta_kubernetes_pod_annotation_prometheus_io_alerts="foo",__meta_kubernetes_pod_annotation_prometheus_io_path="/ops/prometheus/metrics",__meta_kubernetes_pod_annotation_prometheus_io_port="4443",__meta_kubernetes_pod_annotation_prometheus_io_scrape="true",__meta_kubernetes_pod_annotation_sidecar_istio_io_proxyCPU="1.0",__meta_kubernetes_pod_annotation_sidecar_istio_io_proxyCPULimit="1.0",__meta_kubernetes_pod_annotation_sidecar_istio_io_proxyMemory="1.0Gi",__meta_kubernetes_pod_annotation_sidecar_istio_io_proxyMemoryLimit="1.0Gi",__meta_kubernetes_pod_annotation_sidecar_istio_io_status="bar",__meta_kubernetes_pod_annotationpresent_iam_amazonaws_com_role="true",__meta_kubernetes_pod_annotationpresent_istio_io_rev="true",__meta_kubernetes_pod_annotationpresent_kubectl_kubernetes_io_default_container="true",__meta_kubernetes_pod_annotationpresent_kubectl_kubernetes_io_default_logs_container="true",__meta_kubernetes_pod_annotationpresent_nubank_com_br_last_cycle="true",__meta_kubernetes_pod_annotationpresent_nubank_com_br_last_cycle_author="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_alerts="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_path="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_port="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_scrape="true",__meta_kubernetes_pod_annotationpresent_sidecar_istio_io_proxyCPU="true",__meta_kubernetes_pod_annotationpresent_sidecar_istio_io_proxyCPULimit="true",__meta_kubernetes_pod_annotationpresent_sidecar_istio_io_proxyMemory="true",__meta_kubernetes_pod_annotationpresent_sidecar_istio_io_proxyMemoryLimit="true",__meta_kubernetes_pod_annotationpresent_sidecar_istio_io_status="true",__meta_kubernetes_pod_container_id="containerd://2844fa60f9e65155f392e4fb1a5e2594882945288186abf3582ed8c0e0eed88b",__meta_kubernetes_pod_container_image="193814090748.dkr.ecr.us-east-1.amazonaws.com/mirror/istio/proxyv2:1.22.0",__meta_kubernetes_pod_container_init="false",__meta_kubernetes_pod_container_name="istio-proxy",__meta_kubernetes_pod_container_port_name="http-envoy-prom",__meta_kubernetes_pod_container_port_number="15090",__meta_kubernetes_pod_container_port_protocol="TCP",__meta_kubernetes_pod_controller_kind="ReplicaSet",__meta_kubernetes_pod_controller_name="prod-s0-green-shore-deployment-5cb6b89758",__meta_kubernetes_pod_host_ip="10.0.88.44",__meta_kubernetes_pod_ip="100.72.107.150",__meta_kubernetes_pod_label_app_kubernetes_io_name="shore",__meta_kubernetes_pod_label_app_kubernetes_io_version="15296b1b52a49f0ffb52460e52ad33ada4e904e9",__meta_kubernetes_pod_label_nubank_com_br_business_unit="nucore-brazil",__meta_kubernetes_pod_label_nubank_com_br_environment="prod",__meta_kubernetes_pod_label_nubank_com_br_infosec_filter="true",__meta_kubernetes_pod_label_nubank_com_br_mtls_enabled="true",__meta_kubernetes_pod_label_nubank_com_br_name="shore",__meta_kubernetes_pod_label_nubank_com_br_operating_cost_center="120206",__meta_kubernetes_pod_label_nubank_com_br_prototype="s0",__meta_kubernetes_pod_label_nubank_com_br_squad="homepage-platform",__meta_kubernetes_pod_label_nubank_com_br_stack_id="green",__meta_kubernetes_pod_label_nubank_com_br_tier="critical",__meta_kubernetes_pod_label_pod_template_hash="5cb6b89758",__meta_kubernetes_pod_label_security_istio_io_tlsMode="istio",__meta_kubernetes_pod_label_service_istio_io_canonical_name="shore",__meta_kubernetes_pod_label_service_istio_io_canonical_revision="15296b1b52a49f0ffb52460e52ad33ada4e904e9",__meta_kubernetes_pod_label_sidecar_istio_io_inject="true",__meta_kubernetes_pod_labelpresent_app_kubernetes_io_name="true",__meta_kubernetes_pod_labelpresent_app_kubernetes_io_version="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_business_unit="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_environment="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_infosec_filter="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_mtls_enabled="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_name="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_operating_cost_center="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_prototype="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_squad="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_stack_id="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_tier="true",__meta_kubernetes_pod_labelpresent_pod_template_hash="true",__meta_kubernetes_pod_labelpresent_security_istio_io_tlsMode="true",__meta_kubernetes_pod_labelpresent_service_istio_io_canonical_name="true",__meta_kubernetes_pod_labelpresent_service_istio_io_canonical_revision="true",__meta_kubernetes_pod_labelpresent_sidecar_istio_io_inject="true",__meta_kubernetes_pod_name="prod-s0-green-shore-deployment-5cb6b89758-cs8nb",__meta_kubernetes_pod_node_name="ip-10-0-88-44.sa-east-1.compute.internal",__meta_kubernetes_pod_phase="Running",__meta_kubernetes_pod_ready="true",__meta_kubernetes_pod_uid="88384fd8-90a6-4a2b-9584-9261dde86128",__metrics_path__="/stats/prometheus",__scheme__="http",__scrape_interval__="30s",__scrape_timeout__="30s",job="istio-sidecar"}`,

  true,

  `{job="istio-sidecar",kubernetes_pod_name="prod-s0-green-shore-deployment-5cb6b89758-cs8nb",service="shore",tier="critical",business_unit="nucore-brazil",environment="prod",infosec_filter="true",mtls_enabled="true",operating_cost_center="120206",prototype="s0",squad="homepage-platform",stack_id="green"}`)
}