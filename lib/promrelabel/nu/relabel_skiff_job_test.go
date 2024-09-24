package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelSkiffJob(t *testing.T) {
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
    - action: drop
      regex: databases
      source_labels: [__meta_kubernetes_namespace]
    - action: drop
      regex: kube-system
      source_labels: [__meta_kubernetes_namespace]
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
    - {action: labeldrop, regex: __meta_kubernetes_pod_label_nubank_com_br_name}
    - {action: labelmap, regex: __meta_kubernetes_pod_label_nubank_com_br_(.+)}
  `,

  `{__address__="100.72.150.15",__meta_kubernetes_namespace="skiff",__meta_kubernetes_pod_annotation_prometheus_io_path="/ops/prometheus/metrics",__meta_kubernetes_pod_annotation_prometheus_io_port="4443",__meta_kubernetes_pod_annotation_prometheus_io_scheme="http",__meta_kubernetes_pod_annotation_prometheus_io_scrape="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_path="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_port="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_scheme="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_scrape="true",__meta_kubernetes_pod_container_id="containerd://945e1570ef37f8b54476519f8962de55f3e373fd973853cd4fd2d5fb3c01f845",__meta_kubernetes_pod_container_image="193814090748.dkr.ecr.us-east-1.amazonaws.com/nu-skiff:20240521-3",__meta_kubernetes_pod_container_init="false",__meta_kubernetes_pod_container_name="nu-skiff",__meta_kubernetes_pod_controller_kind="DaemonSet",__meta_kubernetes_pod_controller_name="skiff",__meta_kubernetes_pod_host_ip="10.0.94.80",__meta_kubernetes_pod_ip="100.72.150.15",__meta_kubernetes_pod_label_controller_revision_hash="5884465f98",__meta_kubernetes_pod_label_nubank_com_br_environment="prod",__meta_kubernetes_pod_label_nubank_com_br_name="skiff",__meta_kubernetes_pod_label_nubank_com_br_prototype="s0",__meta_kubernetes_pod_label_nubank_com_br_squad="systems-performance",__meta_kubernetes_pod_label_nubank_com_br_stack_id="green",__meta_kubernetes_pod_label_pod_template_generation="1",__meta_kubernetes_pod_labelpresent_controller_revision_hash="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_environment="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_name="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_prototype="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_squad="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_stack_id="true",__meta_kubernetes_pod_labelpresent_pod_template_generation="true",__meta_kubernetes_pod_name="skiff-r2whx",__meta_kubernetes_pod_node_name="ip-10-0-94-80.sa-east-1.compute.internal",__meta_kubernetes_pod_phase="Running",__meta_kubernetes_pod_ready="true",__meta_kubernetes_pod_uid="1994bb34-6bdc-4e00-ad2a-810d12c476ac",__metrics_path__="/metrics",__scheme__="https",__scrape_interval__="30s",__scrape_timeout__="30s",job="skiff-job"}`,

  true,

  `{job="skiff-job",address="100.72.150.15:4443",metrics_path="/ops/prometheus/metrics",scheme="http",kubernetes_namespace="skiff",kubernetes_pod_name="skiff-r2whx",service="skiff",environment="prod",prototype="s0",squad="systems-performance",stack_id="green"}`)
}