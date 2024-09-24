package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelKubernetesPodsHighLatency(t *testing.T) {
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
    - action: keep
      regex: true
      source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_high_latency]
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
    - {action: labeldrop, regex: __meta_kubernetes_pod_label_nubank_com_br_name}
    - {action: labelmap, regex: __meta_kubernetes_pod_label_nubank_com_br_(.+)}
    - {action: labelmap}
    - action: replace
      regex: (.*)
      replacement: global
      source_labels: [__address__]
      target_label: prototype
    - action: replace
      regex: (.*)
      replacement: staging
      source_labels: [__address__]
      target_label: environment
  `,

  `{__address__="100.72.83.38:9106",__meta_kubernetes_namespace="monitoring",__meta_kubernetes_pod_annotation_iam_amazonaws_com_role="prod/cloudwatch-exporter/pink/prod-pink-cloudwatch-exporter-role",__meta_kubernetes_pod_annotation_prometheus_io_alerts="{:http-error-ratio-above-threshold [{:severity :critical, :threshold 0.4, :min-reqs-per-second-per-path 10, :default-instance true} {:severity :warning, :threshold 0.15, :min-reqs-per-second-per-path 0, :default-instance true}], :excessive-cpu-usage-above-requests [{:default-instance true}], :nauvoo-rejection-ratio-above-threshold [{:default-instance true}], :bdc-errors-ratio-above-threshold [{:default-instance true}], :bdc-provider-not-found-by-service-above-threshold [{:default-instance true}], :service-is-down [{:default-instance true} {:for-minutes 15}], :service-is-underprovisioned [{:default-instance true}], :component-ops-health-check-data-not-being-updated [{:default-instance true}], :too-many-ad-hoc-queries [{:default-instance true}], :service-has-too-many-log-types [{:default-instance true}], :certificate-error-above-threshold [{:default-instance true}], :service-is-deploying [{:default-instance true}], :cpu-throttling-above-threshold [{:default-instance true}], :scaling-limited-due-to-current-replicas-being-max-replicas [{:default-instance true}], :cpu-throttling-above-threshold-custom-channel [{:default-instance true}], :service-canary-is-unhealthy [{:default-instance true}], :nauvoo-rejection-ratio-above-threshold-custom-channel-temp [{:default-instance true}], :excessive-cpu-overbooking [{:default-instance true}], :ops-health-failure [{:severity :warning, :default-instance true}], :cloudwatch-exporter-request-rate-above-threshold [{:threshold 270}], :catalyst-error-ratio-above-threshold-v2 [{:default-instance true}], :service-container-frequently-oom-killed [{:default-instance true}], :token-error-above-threshold [{:default-instance true}]}",__meta_kubernetes_pod_annotation_prometheus_io_high_latency="true",__meta_kubernetes_pod_annotation_prometheus_io_job="kubernetes-pods-high-latency",__meta_kubernetes_pod_annotation_prometheus_io_path="/metrics",__meta_kubernetes_pod_annotation_prometheus_io_port="9106",__meta_kubernetes_pod_annotation_prometheus_io_scheme="http",__meta_kubernetes_pod_annotation_prometheus_io_scrape="true",__meta_kubernetes_pod_annotationpresent_iam_amazonaws_com_role="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_alerts="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_high_latency="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_job="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_path="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_port="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_scheme="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_scrape="true",__meta_kubernetes_pod_container_id="containerd://d537f4a1830a031c57c9c68bfeba62ed6e5677b6bf4756c73550df9e8d7a987e",__meta_kubernetes_pod_container_image="193814090748.dkr.ecr.us-east-1.amazonaws.com/nu-cloudwatch-exporter:b61059311db251d4cf056160572f5ae3b63d18af",__meta_kubernetes_pod_container_init="false",__meta_kubernetes_pod_container_name="nu-cloudwatch-exporter",__meta_kubernetes_pod_container_port_name="port9106",__meta_kubernetes_pod_container_port_number="9106",__meta_kubernetes_pod_container_port_protocol="TCP",__meta_kubernetes_pod_controller_kind="ReplicaSet",__meta_kubernetes_pod_controller_name="prod-s0-green-cloudwatch-exporter-deployment-8568864f94",__meta_kubernetes_pod_host_ip="10.0.88.114",__meta_kubernetes_pod_ip="100.72.83.38",__meta_kubernetes_pod_label_app_kubernetes_io_name="cloudwatch-exporter",__meta_kubernetes_pod_label_app_kubernetes_io_version="b61059311db251d4cf056160572f5ae3b63d18af",__meta_kubernetes_pod_label_nubank_com_br_business_unit="ctp",__meta_kubernetes_pod_label_nubank_com_br_environment="prod",__meta_kubernetes_pod_label_nubank_com_br_infosec_filter="false",__meta_kubernetes_pod_label_nubank_com_br_mtls_enabled="false",__meta_kubernetes_pod_label_nubank_com_br_name="cloudwatch-exporter",__meta_kubernetes_pod_label_nubank_com_br_operating_cost_center="120512",__meta_kubernetes_pod_label_nubank_com_br_prototype="s0",__meta_kubernetes_pod_label_nubank_com_br_squad="reliability-metrics",__meta_kubernetes_pod_label_nubank_com_br_stack_id="green",__meta_kubernetes_pod_label_nubank_com_br_tier="important",__meta_kubernetes_pod_label_pod_template_hash="8568864f94",__meta_kubernetes_pod_labelpresent_app_kubernetes_io_name="true",__meta_kubernetes_pod_labelpresent_app_kubernetes_io_version="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_business_unit="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_environment="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_infosec_filter="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_mtls_enabled="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_name="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_operating_cost_center="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_prototype="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_squad="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_stack_id="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_tier="true",__meta_kubernetes_pod_labelpresent_pod_template_hash="true",__meta_kubernetes_pod_name="prod-s0-green-cloudwatch-exporter-deployment-8568864f94-8552h",__meta_kubernetes_pod_node_name="ip-10-0-88-114.sa-east-1.compute.internal",__meta_kubernetes_pod_phase="Running",__meta_kubernetes_pod_ready="true",__meta_kubernetes_pod_uid="89bafd15-173a-44f4-995f-e77671e5d443",__metrics_path__="/metrics",__scheme__="https",__scrape_interval__="2m",__scrape_timeout__="2m",job="kubernetes-pods-high-latency"}`,

  true,

  `{job="kubernetes-pods-high-latency",prototype="global",kubernetes_namespace="monitoring",kubernetes_pod_name="prod-s0-green-cloudwatch-exporter-deployment-8568864f94-8552h",service="cloudwatch-exporter",business_unit="ctp",environment="staging",infosec_filter="false",mtls_enabled="false",operating_cost_center="120512",squad="reliability-metrics",stack_id="green",tier="important"}`)
}