package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelKubernetesDatomicPods(t *testing.T) {
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
      target_label: __metrics_path__
    - action: replace
      regex: ^(https?)$
      source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scheme]
      target_label: __scheme__
    - action: replace
      regex: ([^:]+)(?::\d+)?;(\d+)
      replacement: $1:$2
      source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
      target_label: __address__
    - action: keep
      regex: databases
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
    - action: replace
      source_labels: [__meta_kubernetes_pod_label_nubank_com_br_layer]
      target_label: layer
    - {action: labeldrop, regex: __meta_kubernetes_pod_label_nubank_com_br_name}
    - {action: labelmap, regex: __meta_kubernetes_pod_label_nubank_com_br_(.+)}
    - action: replace
      regex: dev-tx-metrics
      replacement: datomic-restore-transactor
      source_labels: [__meta_kubernetes_pod_container_port_name]
      target_label: service
  `,

  `{__address__="100.72.173.67",__meta_kubernetes_namespace="databases",__meta_kubernetes_pod_annotation_iam_amazonaws_com_role="prod/lannister-datomic/pink/prod-pink-lannister-datomic-role",__meta_kubernetes_pod_annotation_prometheus_io_alerts="{:datomic-active-transactor-heartbeat-failing [{:default-instance true}], :datomic-storage-throttling-write-requests [{:default-instance true}], :datomic-transactor-adaptive-high-memory-index [{:default-instance true}], :datomic-transactor-memcached-hit-ratio-below-threshold [{:default-instance true}], :datomic-transactor-restarts-above-threshold [{:default-instance true}], :excessive-cpu-usage-above-requests [{:default-instance true}], :excessive-cpu-overbooking [{:default-instance true}]}",__meta_kubernetes_pod_annotation_prometheus_io_path="/metrics",__meta_kubernetes_pod_annotation_prometheus_io_port="9100",__meta_kubernetes_pod_annotation_prometheus_io_scheme="http",__meta_kubernetes_pod_annotation_prometheus_io_scrape="true",__meta_kubernetes_pod_annotationpresent_iam_amazonaws_com_role="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_alerts="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_path="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_port="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_scheme="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_scrape="true",__meta_kubernetes_pod_container_id="containerd://9ab3e83a0d3013f2e06471959fe7fdd7651e78ab15d914e4f33193e5ee861da7",__meta_kubernetes_pod_container_image="193814090748.dkr.ecr.us-east-1.amazonaws.com/nu-datomic-init:latest",__meta_kubernetes_pod_container_init="true",__meta_kubernetes_pod_container_name="nu-datomic-init",__meta_kubernetes_pod_controller_kind="ReplicaSet",__meta_kubernetes_pod_controller_name="prod-s0-blue-lannister-datomic-deployment-6dcbd86866",__meta_kubernetes_pod_host_ip="10.0.88.142",__meta_kubernetes_pod_ip="100.72.173.67",__meta_kubernetes_pod_label_app_kubernetes_io_name="lannister-datomic",__meta_kubernetes_pod_label_app_kubernetes_io_version="latest",__meta_kubernetes_pod_label_nubank_com_br_business_unit="credit-card",__meta_kubernetes_pod_label_nubank_com_br_environment="prod",__meta_kubernetes_pod_label_nubank_com_br_layer="blue",__meta_kubernetes_pod_label_nubank_com_br_name="lannister-datomic",__meta_kubernetes_pod_label_nubank_com_br_operating_cost_center="130071",__meta_kubernetes_pod_label_nubank_com_br_prototype="s0",__meta_kubernetes_pod_label_nubank_com_br_squad="gccp-musicc-mpb",__meta_kubernetes_pod_label_nubank_com_br_stack_id="blue",__meta_kubernetes_pod_label_pod_template_hash="6dcbd86866",__meta_kubernetes_pod_labelpresent_app_kubernetes_io_name="true",__meta_kubernetes_pod_labelpresent_app_kubernetes_io_version="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_business_unit="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_environment="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_layer="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_name="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_operating_cost_center="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_prototype="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_squad="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_stack_id="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_tier="true",__meta_kubernetes_pod_labelpresent_pod_template_hash="true",__meta_kubernetes_pod_name="prod-s0-blue-lannister-datomic-deployment-6dcbd86866-rv7sr",__meta_kubernetes_pod_node_name="ip-10-0-88-142.sa-east-1.compute.internal",__meta_kubernetes_pod_phase="Running",__meta_kubernetes_pod_ready="true",__meta_kubernetes_pod_uid="be53445a-194c-478a-b3ac-5ddfe9151805",__metrics_path__="/metrics",__scheme__="https",__scrape_interval__="30s",__scrape_timeout__="30s",job="kubernetes-datomic-pods"}`,

  true,

  `{job="kubernetes-datomic-pods",kubernetes_namespace="databases",kubernetes_pod_name="prod-s0-blue-lannister-datomic-deployment-6dcbd86866-rv7sr",service="lannister-datomic",layer="blue",business_unit="credit-card",environment="prod",operating_cost_center="130071",prototype="s0",squad="gccp-musicc-mpb",stack_id="blue"}`)
}