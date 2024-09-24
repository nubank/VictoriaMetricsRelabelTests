package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelKubernetesPods(t *testing.T) {
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
    - action: drop
      regex: (burrow|burrow-data|magnitude-consumers|otel-collector-metrics|skiff)
      source_labels: [__meta_kubernetes_pod_label_nubank_com_br_name]
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

  `{__address__="100.72.34.197:4445",__meta_kubernetes_namespace="default",__meta_kubernetes_pod_annotation_iam_amazonaws_com_role="prod/talk-box/pink/prod-pink-talk-box-role",__meta_kubernetes_pod_annotation_prometheus_io_alerts="{:http-error-ratio-above-threshold [{:severity :critical, :threshold 0.4, :min-reqs-per-second-per-path 10, :default-instance true} {:severity :warning, :threshold 0.15, :min-reqs-per-second-per-path 0, :default-instance true}], :excessive-cpu-usage-above-requests [{:default-instance true}], :nauvoo-rejection-ratio-above-threshold [{:default-instance true}], :bdc-errors-ratio-above-threshold [{:default-instance true}], :bdc-provider-not-found-by-service-above-threshold [{:default-instance true}], :service-is-down [{:default-instance true} {:severity :critical, :for-minutes 10}], :service-is-underprovisioned [{:default-instance true} {:for-minutes 10}], :component-ops-health-check-data-not-being-updated [{:default-instance true}], :kafka-lag-above-threshold [{:params {:topic \"DCP.TECBAN.NEW-WITHDRAWAL-AUTHORIZATION\"}, :severity :critical, :threshold 5, :for-minutes 1, :annotations {:playbook \"https://nubank.atlassian.net/wiki/spaces/DE/pages/263048760708\", :dashboard \"https://prod-grafana.ist.nubank.world/d/-WI6tgpVk/dcp-life-cycle-monitoring?var-PROMETHEUS=prod-metrics-br&from=now-3h&to=now\"}} {:params {:topic \"DCP.TECBAN.NEW-WITHDRAWAL-CONFIRMATION\"}, :severity :critical, :threshold 5, :for-minutes 1, :annotations {:playbook \"https://nubank.atlassian.net/wiki/spaces/DE/pages/263048858323\", :dashboard \"https://prod-grafana.ist.nubank.world/d/-WI6tgpVk/dcp-life-cycle-monitoring?var-PROMETHEUS=prod-metrics-br&from=now-3h&to=now\"}} {:params {:topic \"DCP.TECBAN.NEW-WITHDRAWAL-RECONCILIATION\"}, :severity :critical, :threshold 5, :for-minutes 1, :annotations {:playbook \"https://nubank.atlassian.net/wiki/spaces/DE/pages/263048760715\", :dashboard \"https://prod-grafana.ist.nubank.world/d/-WI6tgpVk/dcp-life-cycle-monitoring?var-PROMETHEUS=prod-metrics-br&from=now-3h&to=now\"}} {:params {:topic \"DCP.TRANSACTION-REQUESTED\"}, :severity :critical, :threshold 5, :for-minutes 1, :annotations {:playbook \"https://nubank.atlassian.net/wiki/spaces/DE/pages/263047714544\", :dashboard \"https://prod-grafana.ist.nubank.world/d/-WI6tgpVk/dcp-life-cycle-monitoring?var-PROMETHEUS=prod-metrics-br&from=now-3h&to=now\"}}], :too-many-ad-hoc-queries [{:default-instance true}], :service-has-too-many-log-types [{:default-instance true}], :certificate-error-above-threshold [{:default-instance true}], :service-is-deploying [{:default-instance true}], :cpu-throttling-above-threshold [{:default-instance true}], :scaling-limited-due-to-current-replicas-being-max-replicas [{:default-instance true}], :cpu-usage-above-threshold [{:threshold 0.8}], :cpu-throttling-above-threshold-custom-channel [{:default-instance true}], :service-canary-is-unhealthy [{:default-instance true}], :nauvoo-rejection-ratio-above-threshold-custom-channel-temp [{:default-instance true}], :excessive-cpu-overbooking [{:default-instance true}], :ops-health-failure [{:severity :warning, :default-instance true} {:severity :critical, :for-minutes 10}], :catalyst-error-ratio-above-threshold-v2 [{:default-instance true}], :service-container-frequently-oom-killed [{:default-instance true}], :token-error-above-threshold [{:default-instance true}]}",__meta_kubernetes_pod_annotation_prometheus_io_path="/ops/prometheus/metrics",__meta_kubernetes_pod_annotation_prometheus_io_port="4443",__meta_kubernetes_pod_annotation_prometheus_io_scrape="true",__meta_kubernetes_pod_annotationpresent_iam_amazonaws_com_role="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_alerts="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_path="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_port="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_scrape="true",__meta_kubernetes_pod_container_id="containerd://b95ace16330bf70833735a9fb95a3e7b699523c02aa645b142648ae38f18ebc2",__meta_kubernetes_pod_container_image="193814090748.dkr.ecr.us-east-1.amazonaws.com/nu-talk-box:5439dfd1a46a308d328b9f24e26879fa1304b2ea",__meta_kubernetes_pod_container_init="false",__meta_kubernetes_pod_container_name="nu-talk-box",__meta_kubernetes_pod_container_port_name="port4445",__meta_kubernetes_pod_container_port_number="4445",__meta_kubernetes_pod_container_port_protocol="TCP",__meta_kubernetes_pod_controller_kind="ReplicaSet",__meta_kubernetes_pod_controller_name="prod-s0-green-talk-box-deployment-567cb98485",__meta_kubernetes_pod_host_ip="10.0.84.94",__meta_kubernetes_pod_ip="100.72.34.197",__meta_kubernetes_pod_label_app_kubernetes_io_name="talk-box",__meta_kubernetes_pod_label_app_kubernetes_io_version="5439dfd1a46a308d328b9f24e26879fa1304b2ea",__meta_kubernetes_pod_label_nubank_com_br_business_unit="global-bank-account",__meta_kubernetes_pod_label_nubank_com_br_environment="prod",__meta_kubernetes_pod_label_nubank_com_br_infosec_filter="false",__meta_kubernetes_pod_label_nubank_com_br_mtls_enabled="true",__meta_kubernetes_pod_label_nubank_com_br_name="talk-box",__meta_kubernetes_pod_label_nubank_com_br_operating_cost_center="120027",__meta_kubernetes_pod_label_nubank_com_br_prototype="s0",__meta_kubernetes_pod_label_nubank_com_br_squad="gba-debit-foundation",__meta_kubernetes_pod_label_nubank_com_br_stack_id="green",__meta_kubernetes_pod_label_nubank_com_br_tier="important",__meta_kubernetes_pod_label_pod_template_hash="567cb98485",__meta_kubernetes_pod_labelpresent_app_kubernetes_io_name="true",__meta_kubernetes_pod_labelpresent_app_kubernetes_io_version="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_business_unit="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_environment="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_infosec_filter="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_mtls_enabled="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_name="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_operating_cost_center="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_prototype="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_squad="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_stack_id="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_tier="true",__meta_kubernetes_pod_labelpresent_pod_template_hash="true",__meta_kubernetes_pod_name="prod-s0-green-talk-box-deployment-567cb98485-5gp99",__meta_kubernetes_pod_node_name="ip-10-0-84-94.sa-east-1.compute.internal",__meta_kubernetes_pod_phase="Running",__meta_kubernetes_pod_ready="true",__meta_kubernetes_pod_uid="4a894bfd-4641-4da4-9a32-54e8f3b469c0",__metrics_path__="/metrics",__scheme__="https",__scrape_interval__="30s",__scrape_timeout__="30s",job="kubernetes-pods"}`,

  true,

  `{job="kubernetes-pods",metrics_path="/ops/prometheus/metrics",address="100.72.34.197:4443",kubernetes_namespace="default",kubernetes_pod_name="prod-s0-green-talk-box-deployment-567cb98485-5gp99",service="talk-box",business_unit="global-bank-account",environment="prod",infosec_filter="false",mtls_enabled="true",operating_cost_center="120027",prototype="s0",squad="gba-debit-foundation",stack_id="green",tier="important"}`)
}