package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelBrandProtectionGcpMetrics(t *testing.T) {
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

  `{__address__="10.50.73.209",__meta_kubernetes_namespace="default",__meta_kubernetes_pod_annotation_iam_amazonaws_com_role="prod/brand-protection-gcp-metrics-exporter/pink/prod-pink-brand-protection-gcp-metrics-exporter-role",__meta_kubernetes_pod_annotation_prometheus_io_alerts="{:http-error-ratio-above-threshold [{:severity :critical, :threshold 0.4, :min-reqs-per-second-per-path 10, :default-instance true} {:severity :warning, :threshold 0.15, :min-reqs-per-second-per-path 0, :default-instance true}], :excessive-cpu-usage-above-requests [{:default-instance true}], :nauvoo-rejection-ratio-above-threshold [{:default-instance true}], :service-is-down [{:default-instance true}], :service-is-underprovisioned [{:default-instance true} :disabled], :component-ops-health-check-data-not-being-updated [{:default-instance true}], :too-many-ad-hoc-queries [{:default-instance true}], :service-has-too-many-log-types [{:default-instance true}], :certificate-error-above-threshold [{:default-instance true}], :service-is-deploying [{:default-instance true}], :cpu-throttling-above-threshold [{:default-instance true}], :scaling-limited-due-to-current-replicas-being-max-replicas [{:default-instance true}], :service-canary-is-unhealthy [{:default-instance true}], :nauvoo-rejection-ratio-above-threshold-custom-channel-temp [{:default-instance true}], :excessive-cpu-overbooking [{:default-instance true}], :ops-health-failure [{:severity :warning, :default-instance true}], :service-container-frequently-oom-killed [{:default-instance true}], :token-error-above-threshold [{:default-instance true}]}",__meta_kubernetes_pod_annotation_prometheus_io_job="brand-protection-gcp-metrics",__meta_kubernetes_pod_annotation_prometheus_io_path="/metrics",__meta_kubernetes_pod_annotation_prometheus_io_port="9255",__meta_kubernetes_pod_annotation_prometheus_io_scheme="http",__meta_kubernetes_pod_annotation_prometheus_io_scrape="true",__meta_kubernetes_pod_annotationpresent_iam_amazonaws_com_role="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_alerts="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_job="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_path="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_port="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_scheme="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_scrape="true",__meta_kubernetes_pod_container_id="containerd://759fa3061eaf037a7290551b5f9ee5ad58373ce1ada03d22678fa2c1616ebf01",__meta_kubernetes_pod_container_image="193814090748.dkr.ecr.us-east-1.amazonaws.com/nu-brand-protection-gcp-metrics-exporter:b91b5b65208a00be35cf7265f78682111ed30e18",__meta_kubernetes_pod_container_init="false",__meta_kubernetes_pod_container_name="nu-brand-protection-gcp-metrics-exporter",__meta_kubernetes_pod_controller_kind="ReplicaSet",__meta_kubernetes_pod_controller_name="prod-global-blue-brand-protection-gcp-metrics-exporter-deployment-8469c57b4d",__meta_kubernetes_pod_host_ip="10.50.67.125",__meta_kubernetes_pod_ip="10.50.73.209",__meta_kubernetes_pod_label_nubank_com_br_business_unit="shared",__meta_kubernetes_pod_label_nubank_com_br_environment="prod",__meta_kubernetes_pod_label_nubank_com_br_name="brand-protection-gcp-metrics-exporter",__meta_kubernetes_pod_label_nubank_com_br_prototype="global",__meta_kubernetes_pod_label_nubank_com_br_squad="threat-platforms",__meta_kubernetes_pod_label_nubank_com_br_stack_id="blue",__meta_kubernetes_pod_label_nubank_com_br_tier="experimental",__meta_kubernetes_pod_label_pod_template_hash="8469c57b4d",__meta_kubernetes_pod_labelpresent_nubank_com_br_business_unit="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_environment="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_name="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_prototype="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_squad="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_stack_id="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_tier="true",__meta_kubernetes_pod_labelpresent_pod_template_hash="true",__meta_kubernetes_pod_name="prod-global-blue-brand-protection-gcp-metrics-exporter-deprgm48",__meta_kubernetes_pod_node_name="ip-10-50-67-125.us-east-2.compute.internal",__meta_kubernetes_pod_phase="Running",__meta_kubernetes_pod_ready="true",__meta_kubernetes_pod_uid="a81077d6-06f1-4153-bd83-63a3f5c6ad14",__metrics_path__="/metrics",__scheme__="https",__scrape_interval__="2m",__scrape_timeout__="30s",job="brand-protection-gcp-metrics"}`,

  true,

  `{job="brand-protection-gcp-metrics",prototype="global",kubernetes_namespace="default",kubernetes_pod_name="prod-global-blue-brand-protection-gcp-metrics-exporter-deprgm48",service="brand-protection-gcp-metrics-exporter",tier="experimental",business_unit="shared",environment="prod",name="brand-protection-gcp-metrics-exporter",squad="threat-platforms",stack_id="blue"}`)
}