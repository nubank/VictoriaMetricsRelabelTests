package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelKeda(t *testing.T) {
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
      regex: keda-operator
      source_labels: [__meta_kubernetes_pod_container_name]
    - {action: labelmap, regex: __meta_kubernetes_pod_label_(.+)}
    - action: replace
      source_labels: [__meta_kubernetes_namespace]
      target_label: kubernetes_namespace
    - action: replace
      source_labels: [__meta_kubernetes_pod_name]
      target_label: kubernetes_pod_name
  `,

  `{__address__="100.73.236.246:8080",__meta_kubernetes_namespace="keda",__meta_kubernetes_pod_annotation_iam_amazonaws_com_role="prod/keda-operator/green/prod-green-keda-operator-role",__meta_kubernetes_pod_annotationpresent_iam_amazonaws_com_role="true",__meta_kubernetes_pod_container_id="containerd://5b2560e2996b768d1328a2730a9e6cb02c158306cd1af24249b854e364d5d743",__meta_kubernetes_pod_container_image="ghcr.io/kedacore/keda:2.13.0",__meta_kubernetes_pod_container_init="false",__meta_kubernetes_pod_container_name="keda-operator",__meta_kubernetes_pod_container_port_name="http",__meta_kubernetes_pod_container_port_number="8080",__meta_kubernetes_pod_container_port_protocol="TCP",__meta_kubernetes_pod_controller_kind="ReplicaSet",__meta_kubernetes_pod_controller_name="keda-operator-864cc96c97",__meta_kubernetes_pod_host_ip="10.0.201.101",__meta_kubernetes_pod_ip="100.73.236.246",__meta_kubernetes_pod_label_app="keda-operator",__meta_kubernetes_pod_label_name="keda-operator",__meta_kubernetes_pod_label_nubank_com_br_name="keda-operator",__meta_kubernetes_pod_label_pod_template_hash="864cc96c97",__meta_kubernetes_pod_labelpresent_app="true",__meta_kubernetes_pod_labelpresent_name="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_name="true",__meta_kubernetes_pod_labelpresent_pod_template_hash="true",__meta_kubernetes_pod_name="keda-operator-864cc96c97-srjt2",__meta_kubernetes_pod_node_name="ip-10-0-201-101.sa-east-1.compute.internal",__meta_kubernetes_pod_phase="Running",__meta_kubernetes_pod_ready="true",__meta_kubernetes_pod_uid="cc325b4d-ea9a-484c-9ed6-f72dae4dae55",__metrics_path__="/metrics",__scheme__="http",__scrape_interval__="30s",__scrape_timeout__="30s",job="keda"}`,

  true,

  `{job="keda",kubernetes_namespace="keda",kubernetes_pod_name="keda-operator-864cc96c97-srjt2",app="keda-operator",name="keda-operator",nubank_com_br_name="keda-operator",pod_template_hash="864cc96c97"`)
}