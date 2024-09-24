package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelKubernetesPodsHighMetrics(t *testing.T) {
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

  `{__address__="100.72.7.76:4445",__meta_kubernetes_namespace="default",__meta_kubernetes_pod_annotation_iam_amazonaws_com_role="prod/magnitude-consumers/pink/prod-pink-magnitude-consumers-role",__meta_kubernetes_pod_annotation_nubank_com_br_last_cycle="1725486289",__meta_kubernetes_pod_annotation_nubank_com_br_last_cycle_author="gustavo.fernandes2",__meta_kubernetes_pod_annotation_prometheus_io_alerts="{:http-error-ratio-above-threshold [{:severity :critical, :threshold 0.4, :min-reqs-per-second-per-path 10, :default-instance true} {:severity :warning, :threshold 0.15, :min-reqs-per-second-per-path 0, :default-instance true}], :excessive-cpu-usage-above-requests [{:default-instance true}], :nauvoo-rejection-ratio-above-threshold [{:default-instance true}], :bdc-errors-ratio-above-threshold [{:default-instance true}], :bdc-provider-not-found-by-service-above-threshold [{:default-instance true}], :service-is-down [{:default-instance true} {:severity :critical, :for-minutes 30, :params {:job \"kubernetes-pods-high-metrics\"}, :annotations {:channel \"#app-analytics-alarms\"}, :alertmanager-labels {:custom-channel \"true\"}}], :service-is-underprovisioned [{:default-instance true}], :component-ops-health-check-data-not-being-updated [{:default-instance true}], :too-many-ad-hoc-queries [{:default-instance true}], :deadletter-count-above-threshold [{:for-minutes 60, :severity :warning, :threshold 1000, :params {:producer \"magnitude-consumers\"}, :annotations {:channel \"#app-analytics-alarms\"}, :alertmanager-labels {:custom-channel \"true\"}}], :service-has-too-many-log-types [{:default-instance true}], :xp-magnitude-statistics-by-destination-threshold [{:severity :warning, :params {:destination \"prometheus\"}, :annotations {:channel \"#app-analytics-alarms\"}, :alertmanager-labels {:custom-channel \"true\"}} {:severity :warning, :threshold 0.5, :params {:destination \"kafka\"}, :annotations {:channel \"#app-analytics-alarms\"}, :alertmanager-labels {:custom-channel \"true\"}} {:severity :warning, :params {:destination \"etl\"}, :annotations {:channel \"#app-analytics-alarms\"}, :alertmanager-labels {:custom-channel \"true\"}} {:severity :warning, :params {:destination \"splunk\"}, :annotations {:channel \"#app-analytics-alarms\"}, :alertmanager-labels {:custom-channel \"true\"}}], :certificate-error-above-threshold [{:default-instance true}], :service-is-deploying [{:default-instance true}], :cpu-throttling-above-threshold [{:default-instance true}], :scaling-limited-due-to-current-replicas-being-max-replicas [{:default-instance true} :disabled], :cpu-throttling-above-threshold-custom-channel [{:default-instance true} :disabled], :shard-aware-kafka-lag-above-threshold [{:severity :warning, :params {:consumer-group \"MAGNITUDE\"}, :threshold 5000000, :for-minutes 60, :annotations {:channel \"#app-analytics-alarms\"}, :alertmanager-labels {:custom-channel \"true\"}}], :service-canary-is-unhealthy [{:default-instance true}], :nauvoo-rejection-ratio-above-threshold-custom-channel-temp [{:default-instance true}], :excessive-cpu-overbooking [{:default-instance true}], :ops-health-failure [{:severity :warning, :default-instance true} {:severity :warning, :for-minutes 30, :annotations {:channel \"#app-analytics-alarms\"}, :alertmanager-labels {:custom-channel \"true\"}}], :xp-magnitude-statistics-by-source-threshold [{:severity :warning, :params {:source \"Android\"}, :annotations {:channel \"#app-analytics-alarms\"}, :alertmanager-labels {:custom-channel \"true\"}} {:severity :warning, :params {:source \"iOS\"}, :annotations {:channel \"#app-analytics-alarms\"}, :alertmanager-labels {:custom-channel \"true\"}}], :catalyst-error-ratio-above-threshold-v2 [{:default-instance true}], :service-container-frequently-oom-killed [{:default-instance true}], :token-error-above-threshold [{:default-instance true}]}",__meta_kubernetes_pod_annotation_prometheus_io_path="/ops/prometheus/metrics",__meta_kubernetes_pod_annotation_prometheus_io_port="4443",__meta_kubernetes_pod_annotation_prometheus_io_scrape="true",__meta_kubernetes_pod_annotationpresent_iam_amazonaws_com_role="true",__meta_kubernetes_pod_annotationpresent_nubank_com_br_last_cycle="true",__meta_kubernetes_pod_annotationpresent_nubank_com_br_last_cycle_author="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_alerts="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_path="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_port="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_scrape="true",__meta_kubernetes_pod_container_id="containerd://909806e3aebe0b575f433ceb5c741be20417dd29c7fd707e040a32aa3917234f",__meta_kubernetes_pod_container_image="193814090748.dkr.ecr.us-east-1.amazonaws.com/nu-magnitude-consumers:4153dda88b306221340850dc7349c8eeeeef2d32",__meta_kubernetes_pod_container_init="false",__meta_kubernetes_pod_container_name="nu-magnitude-consumers",__meta_kubernetes_pod_container_port_name="port4445",__meta_kubernetes_pod_container_port_number="4445",__meta_kubernetes_pod_container_port_protocol="TCP",__meta_kubernetes_pod_controller_kind="ReplicaSet",__meta_kubernetes_pod_controller_name="prod-s0-green-magnitude-consumers-deployment-564db96bdd",__meta_kubernetes_pod_host_ip="10.0.82.174",__meta_kubernetes_pod_ip="100.72.7.76",__meta_kubernetes_pod_label_app_kubernetes_io_name="magnitude-consumers",__meta_kubernetes_pod_label_app_kubernetes_io_version="4153dda88b306221340850dc7349c8eeeeef2d32",__meta_kubernetes_pod_label_nubank_com_br_business_unit="foundation",__meta_kubernetes_pod_label_nubank_com_br_environment="prod",__meta_kubernetes_pod_label_nubank_com_br_infosec_filter="false",__meta_kubernetes_pod_label_nubank_com_br_mtls_enabled="true",__meta_kubernetes_pod_label_nubank_com_br_name="magnitude-consumers",__meta_kubernetes_pod_label_nubank_com_br_operating_cost_center="130026",__meta_kubernetes_pod_label_nubank_com_br_prototype="s0",__meta_kubernetes_pod_label_nubank_com_br_squad="app-analytics",__meta_kubernetes_pod_label_nubank_com_br_stack_id="green",__meta_kubernetes_pod_label_nubank_com_br_tier="important",__meta_kubernetes_pod_label_pod_template_hash="564db96bdd",__meta_kubernetes_pod_labelpresent_app_kubernetes_io_name="true",__meta_kubernetes_pod_labelpresent_app_kubernetes_io_version="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_business_unit="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_environment="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_infosec_filter="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_mtls_enabled="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_name="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_operating_cost_center="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_prototype="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_squad="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_stack_id="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_tier="true",__meta_kubernetes_pod_labelpresent_pod_template_hash="true",__meta_kubernetes_pod_name="prod-s0-green-magnitude-consumers-deployment-564db96bdd-bv4g6",__meta_kubernetes_pod_node_name="ip-10-0-82-174.sa-east-1.compute.internal",__meta_kubernetes_pod_phase="Running",__meta_kubernetes_pod_ready="true",__meta_kubernetes_pod_uid="e6df503c-a82a-416b-9a23-737c8087e18f",__metrics_path__="/metrics",__scheme__="https",__scrape_interval__="1m",__scrape_timeout__="30s",job="kubernetes-pods-high-metrics"}`,

  true,

  `{job="kubernetes-pods-high-metrics",kubernetes_namespace="default",kubernetes_pod_name="prod-s0-green-magnitude-consumers-deployment-564db96bdd-bv4g6",service="magnitude-consumers",tier="important",business_unit="foundation",environment="prod",infosec_filter="false",mtls_enabled="true",operating_cost_center="130026",prototype="s0",squad="app-analytics",stack_id="green"}`)
}