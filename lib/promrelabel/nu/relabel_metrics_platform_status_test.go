package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelMetricsPlatformStatus(t *testing.T) {
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
      regex: .*
      replacement: global
      source_labels: [__address__]
      target_label: prototype
    - action: keep
      regex: true
      source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
    - action: drop
      regex: .+
      source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_high_latency]
    - action: replace
      regex: (.+)
      source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_path]
      target_label: metrics_path
    - action: replace
      regex: ^(https?)$
      source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scheme]
      target_label: scheme
    - action: replace
      regex: ([^:]+)(?::\d+)?;(\d+)
      replacement: $1:$2
      source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
      target_label: address
    - action: replace
      source_labels: [__meta_kubernetes_namespace]
      target_label: kubernetes_namespace
    - action: replace
      source_labels: [__meta_kubernetes_pod_name]
      target_label: kubernetes_pod_name
    - action: replace
      source_labels: [__meta_kubernetes_pod_label_nubank_com_br_name]
      target_label: service
    - action: replace
      source_labels: [__meta_kubernetes_pod_label_nubank_com_br_tier]
      target_label: tier
    - {action: labelmap, regex: __meta_kubernetes_pod_label_nubank_com_br_(.+)}
  `,

  `{__address__="100.71.188.49:8401",__meta_kubernetes_namespace="monitoring",__meta_kubernetes_pod_annotation_iam_amazonaws_com_role="prod/vicmetrics-storage/pink/prod-pink-vicmetrics-storage-role",__meta_kubernetes_pod_annotation_prometheus_io_path="/metrics",__meta_kubernetes_pod_annotation_prometheus_io_port="8482",__meta_kubernetes_pod_annotation_prometheus_io_scheme="http",__meta_kubernetes_pod_annotation_prometheus_io_scrape="true",__meta_kubernetes_pod_annotationpresent_iam_amazonaws_com_role="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_path="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_port="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_scheme="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_scrape="true",__meta_kubernetes_pod_container_id="containerd://67ed3c7b563558b131dd8448d06a85cb0b69882af2a3ec5dc3beed19281f3d69",__meta_kubernetes_pod_container_image="193814090748.dkr.ecr.us-east-1.amazonaws.com/nu-vicmetrics-storage:46a07ffea81c30332eee68d33cb63a775e04a001",__meta_kubernetes_pod_container_init="false",__meta_kubernetes_pod_container_name="nu-vicmetrics-storage",__meta_kubernetes_pod_container_port_name="query",__meta_kubernetes_pod_container_port_number="8401",__meta_kubernetes_pod_container_port_protocol="TCP",__meta_kubernetes_pod_controller_kind="StatefulSet",__meta_kubernetes_pod_controller_name="prod-s0-green-vicmetrics-storage",__meta_kubernetes_pod_host_ip="10.0.233.208",__meta_kubernetes_pod_ip="100.71.188.49",__meta_kubernetes_pod_label_apps_kubernetes_io_pod_index="19",__meta_kubernetes_pod_label_controller_revision_hash="prod-s0-green-vicmetrics-storage-59f6b8c49d",__meta_kubernetes_pod_label_nubank_com_br_environment="prod",__meta_kubernetes_pod_label_nubank_com_br_name="vicmetrics-storage",__meta_kubernetes_pod_label_nubank_com_br_prototype="s0",__meta_kubernetes_pod_label_nubank_com_br_squad="reliability-metrics",__meta_kubernetes_pod_label_nubank_com_br_stack_id="green",__meta_kubernetes_pod_label_nubank_com_br_tier="critical",__meta_kubernetes_pod_label_spotinst_io_restrict_scale_down="true",__meta_kubernetes_pod_label_statefulset_kubernetes_io_pod_name="prod-s0-green-vicmetrics-storage-19",__meta_kubernetes_pod_labelpresent_apps_kubernetes_io_pod_index="true",__meta_kubernetes_pod_labelpresent_controller_revision_hash="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_environment="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_name="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_prototype="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_squad="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_stack_id="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_tier="true",__meta_kubernetes_pod_labelpresent_spotinst_io_restrict_scale_down="true",__meta_kubernetes_pod_labelpresent_statefulset_kubernetes_io_pod_name="true",__meta_kubernetes_pod_name="prod-s0-green-vicmetrics-storage-19",__meta_kubernetes_pod_node_name="ip-10-0-233-208.sa-east-1.compute.internal",__meta_kubernetes_pod_phase="Running",__meta_kubernetes_pod_ready="true",__meta_kubernetes_pod_uid="e9981946-ec08-47f6-b813-83cd491ce749",__metrics_path__="/metrics",__scheme__="https",__scrape_interval__="5s",__scrape_timeout__="5s",job="metrics-platform-status"}`,

  true,

  `{address="100.71.188.49:8482",metrics_path="/metrics",scheme="http",job="metrics-platform-status",kubernetes_namespace="monitoring",kubernetes_pod_name="prod-s0-green-vicmetrics-storage-19",service="vicmetrics-storage",environment="prod",name="vicmetrics-storage",prototype="s0",squad="reliability-metrics",stack_id="green",tier="critical"}`)
}