package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelKedaMetricsAdapter(t *testing.T) {
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
    - action: keep
      regex: keda-metrics-apiserver
      source_labels: [__meta_kubernetes_pod_container_name]
    - {action: labelmap, regex: __meta_kubernetes_pod_label_(.+)}
    - action: replace
      source_labels: [__meta_kubernetes_namespace]
      target_label: kubernetes_namespace
    - action: replace
      source_labels: [__meta_kubernetes_pod_name]
      target_label: kubernetes_pod_name
    - regex: (.*):[0-9]+
      replacement: $${1}:9022
      source_labels: [__address__]
      target_label: address
  `,

  `{__address__="100.72.112.228:8080",__meta_kubernetes_namespace="keda",__meta_kubernetes_pod_container_id="containerd://0a22312593f6bd55543254416269e1a98ea419a9ad4fa9e5c14783516c9c810c",__meta_kubernetes_pod_container_image="ghcr.io/kedacore/keda-metrics-apiserver:2.13.0",__meta_kubernetes_pod_container_init="false",__meta_kubernetes_pod_container_name="keda-metrics-apiserver",__meta_kubernetes_pod_container_port_name="http",__meta_kubernetes_pod_container_port_number="8080",__meta_kubernetes_pod_container_port_protocol="TCP",__meta_kubernetes_pod_controller_kind="ReplicaSet",__meta_kubernetes_pod_controller_name="keda-metrics-apiserver-68676b8c95",__meta_kubernetes_pod_host_ip="10.0.81.90",__meta_kubernetes_pod_ip="100.72.112.228",__meta_kubernetes_pod_label_app="keda-metrics-apiserver",__meta_kubernetes_pod_label_nubank_com_br_name="keda-metrics-apiserver",__meta_kubernetes_pod_label_pod_template_hash="68676b8c95",__meta_kubernetes_pod_labelpresent_app="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_name="true",__meta_kubernetes_pod_labelpresent_pod_template_hash="true",__meta_kubernetes_pod_name="keda-metrics-apiserver-68676b8c95-tfdqq",__meta_kubernetes_pod_node_name="ip-10-0-81-90.sa-east-1.compute.internal",__meta_kubernetes_pod_phase="Running",__meta_kubernetes_pod_ready="true",__meta_kubernetes_pod_uid="87848d75-edbe-4b1e-b4b4-43bc3249d5b9",__metrics_path__="/metrics",__scheme__="http",__scrape_interval__="30s",__scrape_timeout__="30s",job="keda-metrics-adapter"}`,

  true,

  `{job="keda-metrics-adapter",kubernetes_namespace="keda",kubernetes_pod_name="keda-metrics-apiserver-68676b8c95-tfdqq",app="keda-metrics-apiserver",nubank_com_br_name="keda-metrics-apiserver",pod_template_hash="68676b8c95",address="100.72.112.228:9022"}`)
}