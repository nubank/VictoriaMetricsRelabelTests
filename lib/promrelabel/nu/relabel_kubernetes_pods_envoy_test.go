package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelKubernetesPodsEnvoy(t *testing.T) {
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
    - action: replace
      regex: ([^:]+)
      replacement: $${1}:9901
      source_labels: [__address__]
      target_label: __address__
    - action: replace
      source_labels: [__meta_kubernetes_namespace]
      target_label: kubernetes_namespace
    - action: replace
      source_labels: [__meta_kubernetes_pod_name]
      target_label: kubernetes_pod_name
    - action: keep
      regex: scatterbrain
      source_labels: [__meta_kubernetes_pod_label_nubank_com_br_name]
    - action: replace
      source_labels: [__meta_kubernetes_pod_label_nubank_com_br_name]
      target_label: service
    - action: replace
      source_labels: [__meta_kubernetes_pod_label_nubank_com_br_tier]
      target_label: tier
    - {action: labeldrop, regex: __meta_kubernetes_pod_label_nubank_com_br_name}
    - {action: labelmap, regex: __meta_kubernetes_pod_label_nubank_com_br_(.+)}
  `,

  `{__address__="100.73.175.188:4445",__meta_kubernetes_namespace="default",__meta_kubernetes_pod_annotation_iam_amazonaws_com_role="prod/scatterbrain/pink/prod-pink-scatterbrain-role",__meta_kubernetes_pod_annotation_nubank_com_br_last_cycle="1695149417",__meta_kubernetes_pod_annotation_nubank_com_br_last_cycle_author="felipe.brigalante",__meta_kubernetes_pod_annotation_prometheus_io_alerts="{:http-error-ratio-above-threshold [{:severity :critical, :threshold 0.4, :min-reqs-per-second-per-path 10, :default-instance true} {:severity :warning, :threshold 0.15, :min-reqs-per-second-per-path 0, :default-instance true}], :excessive-cpu-usage-above-requests [{:default-instance true}], :nauvoo-rejection-ratio-above-threshold [{:default-instance true}], :bdc-errors-ratio-above-threshold [{:default-instance true}], :bdc-provider-not-found-by-service-above-threshold [{:default-instance true}], :service-is-down [{:default-instance true}], :service-is-underprovisioned [{:default-instance true}], :component-ops-health-check-data-not-being-updated [{:default-instance true}], :too-many-ad-hoc-queries [{:default-instance true}], :service-has-too-many-log-types [{:default-instance true}], :certificate-error-above-threshold [{:default-instance true}], :service-is-deploying [{:default-instance true}], :cpu-throttling-above-threshold [{:default-instance true}], :scaling-limited-due-to-current-replicas-being-max-replicas [{:default-instance true}], :cpu-throttling-above-threshold-custom-channel [{:default-instance true}], :service-canary-is-unhealthy [{:default-instance true}], :nauvoo-rejection-ratio-above-threshold-custom-channel-temp [{:default-instance true}], :excessive-cpu-overbooking [{:default-instance true}], :ops-health-failure [{:severity :warning, :default-instance true} :disabled], :catalyst-error-ratio-above-threshold-v2 [{:default-instance true}], :service-container-frequently-oom-killed [{:default-instance true}], :token-error-above-threshold [{:default-instance true}], :ops-health-failure-excluding-component [{:params {:component \"http.*\"}}]}",__meta_kubernetes_pod_annotation_prometheus_io_path="/ops/prometheus/metrics",__meta_kubernetes_pod_annotation_prometheus_io_port="4443",__meta_kubernetes_pod_annotation_prometheus_io_scrape="true",__meta_kubernetes_pod_annotationpresent_iam_amazonaws_com_role="true",__meta_kubernetes_pod_annotationpresent_nubank_com_br_last_cycle="true",__meta_kubernetes_pod_annotationpresent_nubank_com_br_last_cycle_author="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_alerts="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_path="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_port="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_scrape="true",__meta_kubernetes_pod_container_id="containerd://23c279301fc1b11810aeec5eea5dbc5185e9e9881651e5f26320a4904b0a11c2",__meta_kubernetes_pod_container_image="193814090748.dkr.ecr.us-east-1.amazonaws.com/nu-scatterbrain:17d9b09dba44a9956d6ec95816cf31542e1cb5bf",__meta_kubernetes_pod_container_init="false",__meta_kubernetes_pod_container_name="nu-scatterbrain",__meta_kubernetes_pod_container_port_name="port4445",__meta_kubernetes_pod_container_port_number="4445",__meta_kubernetes_pod_container_port_protocol="TCP",__meta_kubernetes_pod_controller_kind="ReplicaSet",__meta_kubernetes_pod_controller_name="prod-s0-green-scatterbrain-deployment-697f59f84b",__meta_kubernetes_pod_host_ip="10.0.196.51",__meta_kubernetes_pod_ip="100.73.175.188",__meta_kubernetes_pod_label_app_kubernetes_io_name="scatterbrain",__meta_kubernetes_pod_label_app_kubernetes_io_version="17d9b09dba44a9956d6ec95816cf31542e1cb5bf",__meta_kubernetes_pod_label_nubank_com_br_aware_of_shards="sharded",__meta_kubernetes_pod_label_nubank_com_br_business_unit="ctp",__meta_kubernetes_pod_label_nubank_com_br_environment="prod",__meta_kubernetes_pod_label_nubank_com_br_infosec_filter="false",__meta_kubernetes_pod_label_nubank_com_br_mtls_enabled="true",__meta_kubernetes_pod_label_nubank_com_br_name="scatterbrain",__meta_kubernetes_pod_label_nubank_com_br_operating_cost_center="140018",__meta_kubernetes_pod_label_nubank_com_br_prototype="s0",__meta_kubernetes_pod_label_nubank_com_br_shard_aware="true",__meta_kubernetes_pod_label_nubank_com_br_squad="traffic-management",__meta_kubernetes_pod_label_nubank_com_br_stack_id="green",__meta_kubernetes_pod_label_nubank_com_br_tier="critical",__meta_kubernetes_pod_label_pod_template_hash="697f59f84b",__meta_kubernetes_pod_labelpresent_app_kubernetes_io_name="true",__meta_kubernetes_pod_labelpresent_app_kubernetes_io_version="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_aware_of_shards="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_business_unit="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_environment="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_infosec_filter="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_mtls_enabled="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_name="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_operating_cost_center="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_prototype="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_shard_aware="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_squad="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_stack_id="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_tier="true",__meta_kubernetes_pod_labelpresent_pod_template_hash="true",__meta_kubernetes_pod_name="prod-s0-green-scatterbrain-deployment-697f59f84b-9rvx2",__meta_kubernetes_pod_node_name="ip-10-0-196-51.sa-east-1.compute.internal",__meta_kubernetes_pod_phase="Running",__meta_kubernetes_pod_ready="true",__meta_kubernetes_pod_uid="4e1ec2f3-fdfb-478f-95f1-c5b9713b5bf4",__metrics_path__="/stats/prometheus",__scheme__="http",__scrape_interval__="30s",__scrape_timeout__="30s",job="kubernetes-pods-envoy"}`,

  true,

  `{job="kubernetes-pods-envoy",__address__="100.73.175.188:9901",kubernetes_namespace="default",kubernetes_pod_name="prod-s0-green-scatterbrain-deployment-697f59f84b-9rvx2",service="scatterbrain",aware_of_shards="sharded",business_unit="ctp",environment="prod",infosec_filter="false",mtls_enabled="true",operating_cost_center="140018",prototype="s0",shard_aware="true",squad="traffic-management",stack_id="green",tier="critical"}`)
}