package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelOpsHealthTrigger(t *testing.T) {
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
    - action: drop
      regex: vicmetrics.*
      source_labels: [__meta_kubernetes_pod_label_nubank_com_br_name]
    - action: keep
      regex: true
      source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scrape]
    - action: replace
      regex: ^(https?)$
      source_labels: [__meta_kubernetes_pod_annotation_prometheus_io_scheme]
      target_label: scheme
    - action: replace
      regex: ([^:]+)(?::\d+)?;(\d+)
      replacement: $1:$2
      source_labels: [__address__, __meta_kubernetes_pod_annotation_prometheus_io_port]
      target_label: address
  `,

  `{__address__="100.72.250.225:4445",__meta_kubernetes_namespace="default",__meta_kubernetes_pod_annotation_iam_amazonaws_com_role="prod/cheshire/pink/prod-pink-cheshire-role",__meta_kubernetes_pod_annotation_nubank_com_br_last_cycle="1700778255",__meta_kubernetes_pod_annotation_nubank_com_br_last_cycle_author="wennder.santos",__meta_kubernetes_pod_annotation_prometheus_io_alerts="{:http-error-ratio-above-threshold [{:severity :critical, :threshold 0.4, :min-reqs-per-second-per-path 10, :default-instance true} {:severity :warning, :threshold 0.15, :min-reqs-per-second-per-path 0, :default-instance true}], :excessive-cpu-usage-above-requests [{:default-instance true}], :deadletter-stuck-above-threshold [{:threshold 0, :for-minutes 60, :severity :warning}], :nauvoo-rejection-ratio-above-threshold [{:default-instance true}], :bdc-errors-ratio-above-threshold [{:default-instance true}], :bdc-provider-not-found-by-service-above-threshold [{:default-instance true}], :service-is-down [{:default-instance true}], :service-is-underprovisioned [{:default-instance true}], :component-ops-health-check-data-not-being-updated [{:default-instance true}], :kafka-lag-above-threshold [{:severity :warning, :threshold 5, :for-minutes 10, :params {:consumer-group \"CHESHIRE\"}}], :too-many-ad-hoc-queries [{:default-instance true}], :deadletter-count-above-threshold [{:severity :warning, :threshold 1, :for-minutes 10, :params {:producer \"cheshire\", :excluded-topic \"NUCOIN.EXCHANGE-FINISHED\"}} {:severity :warning, :threshold 50, :for-minutes 10, :params {:topic \"NUCOIN.EXCHANGE-FINISHED\", :producer \"cheshire\"}}], :service-has-too-many-log-types [{:default-instance true}], :certificate-error-above-threshold [{:default-instance true}], :service-is-deploying [{:default-instance true}], :cpu-throttling-above-threshold [{:default-instance true}], :scaling-limited-due-to-current-replicas-being-max-replicas [{:default-instance true}], :cpu-throttling-above-threshold-custom-channel [{:default-instance true}], :service-canary-is-unhealthy [{:default-instance true}], :nauvoo-rejection-ratio-above-threshold-custom-channel-temp [{:default-instance true}], :excessive-cpu-overbooking [{:default-instance true}], :ops-health-failure [{:severity :warning, :default-instance true}], :catalyst-error-ratio-above-threshold-v2 [{:default-instance true}], :service-container-frequently-oom-killed [{:default-instance true}], :token-error-above-threshold [{:default-instance true}]}",__meta_kubernetes_pod_annotation_prometheus_io_path="/ops/prometheus/metrics",__meta_kubernetes_pod_annotation_prometheus_io_port="4443",__meta_kubernetes_pod_annotation_prometheus_io_scrape="true",__meta_kubernetes_pod_annotationpresent_iam_amazonaws_com_role="true",__meta_kubernetes_pod_annotationpresent_nubank_com_br_last_cycle="true",__meta_kubernetes_pod_annotationpresent_nubank_com_br_last_cycle_author="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_alerts="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_path="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_port="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_scrape="true",__meta_kubernetes_pod_container_id="containerd://deb0bc51cfda51cf1664e21f63c0c66e20fc5109801f740a89ce25c0cdc6e171",__meta_kubernetes_pod_container_image="193814090748.dkr.ecr.us-east-1.amazonaws.com/nu-cheshire:fcaa3193db4e7bcea287ada449ba28f15a17ae46",__meta_kubernetes_pod_container_init="false",__meta_kubernetes_pod_container_name="nu-cheshire",__meta_kubernetes_pod_container_port_name="port4445",__meta_kubernetes_pod_container_port_number="4445",__meta_kubernetes_pod_container_port_protocol="TCP",__meta_kubernetes_pod_controller_kind="ReplicaSet",__meta_kubernetes_pod_controller_name="prod-s0-green-cheshire-deployment-5bf7dcd59b",__meta_kubernetes_pod_host_ip="10.0.86.80",__meta_kubernetes_pod_ip="100.72.250.225",__meta_kubernetes_pod_label_app_kubernetes_io_name="cheshire",__meta_kubernetes_pod_label_app_kubernetes_io_version="fcaa3193db4e7bcea287ada449ba28f15a17ae46",__meta_kubernetes_pod_label_nubank_com_br_business_unit="nucoin",__meta_kubernetes_pod_label_nubank_com_br_environment="prod",__meta_kubernetes_pod_label_nubank_com_br_infosec_filter="false",__meta_kubernetes_pod_label_nubank_com_br_mtls_enabled="true",__meta_kubernetes_pod_label_nubank_com_br_name="cheshire",__meta_kubernetes_pod_label_nubank_com_br_operating_cost_center="000006",__meta_kubernetes_pod_label_nubank_com_br_prototype="s0",__meta_kubernetes_pod_label_nubank_com_br_squad="nucoin",__meta_kubernetes_pod_label_nubank_com_br_stack_id="green",__meta_kubernetes_pod_label_nubank_com_br_tier="critical",__meta_kubernetes_pod_label_pod_template_hash="5bf7dcd59b",__meta_kubernetes_pod_labelpresent_app_kubernetes_io_name="true",__meta_kubernetes_pod_labelpresent_app_kubernetes_io_version="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_business_unit="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_environment="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_infosec_filter="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_mtls_enabled="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_name="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_operating_cost_center="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_prototype="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_squad="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_stack_id="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_tier="true",__meta_kubernetes_pod_labelpresent_pod_template_hash="true",__meta_kubernetes_pod_name="prod-s0-green-cheshire-deployment-5bf7dcd59b-84wh9",__meta_kubernetes_pod_node_name="ip-10-0-86-80.sa-east-1.compute.internal",__meta_kubernetes_pod_phase="Running",__meta_kubernetes_pod_ready="true",__meta_kubernetes_pod_uid="832be124-096f-44ec-bbe2-885c02e22bd8",__metrics_path__="/ops/health",__scheme__="https",__scrape_interval__="30s",__scrape_timeout__="30s",job="ops-health-trigger"}`,

  true,

  `{job="ops-health-trigger",address="100.72.250.225:4443"}`)
}