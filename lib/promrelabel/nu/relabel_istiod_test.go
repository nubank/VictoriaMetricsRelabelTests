package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelIstiod(t *testing.T) {
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
      regex: istiod;http-monitoring
      source_labels: [__meta_kubernetes_service_name, __meta_kubernetes_endpoint_port_name]
    - action: replace
      source_labels: [__meta_kubernetes_pod_name]
      target_label: kubernetes_pod_name
    - action: replace
      source_labels: [__meta_kubernetes_pod_label_istio]
      target_label: service
    - replacement: critical
      source_labels: [__address__]
      target_label: tier
    - {action: labeldrop, regex: __meta_kubernetes_pod_label_nubank_com_br_name}
    - {action: labelmap, regex: __meta_kubernetes_pod_label_nubank_com_br_(.+)}
  `,

  `{__address__="100.71.129.25:15014",__meta_kubernetes_endpoint_address_target_kind="Pod",__meta_kubernetes_endpoint_address_target_name="istiod-5fd579679b-v2kck",__meta_kubernetes_endpoint_node_name="ip-10-0-239-95.sa-east-1.compute.internal",__meta_kubernetes_endpoint_port_name="http-monitoring",__meta_kubernetes_endpoint_port_protocol="TCP",__meta_kubernetes_endpoint_ready="true",__meta_kubernetes_endpoints_label_app="istiod",__meta_kubernetes_endpoints_label_app_kubernetes_io_instance="istio",__meta_kubernetes_endpoints_label_install_operator_istio_io_owning_resource="unknown",__meta_kubernetes_endpoints_label_istio="pilot",__meta_kubernetes_endpoints_label_istio_io_rev="default",__meta_kubernetes_endpoints_label_operator_istio_io_component="Pilot",__meta_kubernetes_endpoints_label_release="istiod",__meta_kubernetes_endpoints_labelpresent_app="true",__meta_kubernetes_endpoints_labelpresent_app_kubernetes_io_instance="true",__meta_kubernetes_endpoints_labelpresent_install_operator_istio_io_owning_resource="true",__meta_kubernetes_endpoints_labelpresent_istio="true",__meta_kubernetes_endpoints_labelpresent_istio_io_rev="true",__meta_kubernetes_endpoints_labelpresent_operator_istio_io_component="true",__meta_kubernetes_endpoints_labelpresent_release="true",__meta_kubernetes_endpoints_name="istiod",__meta_kubernetes_namespace="istio-system",__meta_kubernetes_pod_annotation_ambient_istio_io_redirection="disabled",__meta_kubernetes_pod_annotation_prometheus_io_port="15014",__meta_kubernetes_pod_annotation_prometheus_io_scrape="true",__meta_kubernetes_pod_annotation_sidecar_istio_io_inject="false",__meta_kubernetes_pod_annotationpresent_ambient_istio_io_redirection="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_port="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_scrape="true",__meta_kubernetes_pod_annotationpresent_sidecar_istio_io_inject="true",__meta_kubernetes_pod_controller_kind="ReplicaSet",__meta_kubernetes_pod_controller_name="istiod-5fd579679b",__meta_kubernetes_pod_host_ip="10.0.239.95",__meta_kubernetes_pod_ip="100.71.129.25",__meta_kubernetes_pod_label_app="istiod",__meta_kubernetes_pod_label_install_operator_istio_io_owning_resource="unknown",__meta_kubernetes_pod_label_istio="pilot",__meta_kubernetes_pod_label_istio_io_rev="default",__meta_kubernetes_pod_label_nubank_com_br_business_unit="ctp",__meta_kubernetes_pod_label_nubank_com_br_environment="prod",__meta_kubernetes_pod_label_nubank_com_br_name="istiod",__meta_kubernetes_pod_label_nubank_com_br_operating_cost_center="140018",__meta_kubernetes_pod_label_nubank_com_br_prototype="s0",__meta_kubernetes_pod_label_nubank_com_br_squad="traffic-management",__meta_kubernetes_pod_label_nubank_com_br_stack_id="green",__meta_kubernetes_pod_label_nubank_com_br_tier="critical",__meta_kubernetes_pod_label_operator_istio_io_component="Pilot",__meta_kubernetes_pod_label_pod_template_hash="5fd579679b",__meta_kubernetes_pod_label_sidecar_istio_io_inject="false",__meta_kubernetes_pod_labelpresent_app="true",__meta_kubernetes_pod_labelpresent_install_operator_istio_io_owning_resource="true",__meta_kubernetes_pod_labelpresent_istio="true",__meta_kubernetes_pod_labelpresent_istio_io_rev="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_business_unit="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_environment="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_name="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_operating_cost_center="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_prototype="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_squad="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_stack_id="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_tier="true",__meta_kubernetes_pod_labelpresent_operator_istio_io_component="true",__meta_kubernetes_pod_labelpresent_pod_template_hash="true",__meta_kubernetes_pod_labelpresent_sidecar_istio_io_inject="true",__meta_kubernetes_pod_name="istiod-5fd579679b-v2kck",__meta_kubernetes_pod_node_name="ip-10-0-239-95.sa-east-1.compute.internal",__meta_kubernetes_pod_phase="Running",__meta_kubernetes_pod_ready="true",__meta_kubernetes_pod_uid="79fc90ea-65a4-4a6b-b4bf-5f241dc69413",__meta_kubernetes_service_annotation_kubectl_kubernetes_io_last_applied_configuration="{\"apiVersion\":\"v1\",\"kind\":\"Service\",\"metadata\":{\"annotations\":{},\"labels\":{\"app\":\"istiod\",\"app.kubernetes.io/instance\":\"istio\",\"install.operator.istio.io/owning-resource\":\"unknown\",\"istio\":\"pilot\",\"istio.io/rev\":\"default\",\"operator.istio.io/component\":\"Pilot\",\"release\":\"istiod\"},\"name\":\"istiod\",\"namespace\":\"istio-system\"},\"spec\":{\"ports\":[{\"name\":\"grpc-xds\",\"port\":15010,\"protocol\":\"TCP\"},{\"name\":\"https-dns\",\"port\":15012,\"protocol\":\"TCP\"},{\"name\":\"https-webhook\",\"port\":443,\"protocol\":\"TCP\",\"targetPort\":15017},{\"name\":\"http-monitoring\",\"port\":15014,\"protocol\":\"TCP\"}],\"selector\":{\"app\":\"istiod\",\"istio\":\"pilot\"}}}",__meta_kubernetes_service_annotationpresent_kubectl_kubernetes_io_last_applied_configuration="true",__meta_kubernetes_service_label_app="istiod",__meta_kubernetes_service_label_app_kubernetes_io_instance="istio",__meta_kubernetes_service_label_install_operator_istio_io_owning_resource="unknown",__meta_kubernetes_service_label_istio="pilot",__meta_kubernetes_service_label_istio_io_rev="default",__meta_kubernetes_service_label_operator_istio_io_component="Pilot",__meta_kubernetes_service_label_release="istiod",__meta_kubernetes_service_labelpresent_app="true",__meta_kubernetes_service_labelpresent_app_kubernetes_io_instance="true",__meta_kubernetes_service_labelpresent_install_operator_istio_io_owning_resource="true",__meta_kubernetes_service_labelpresent_istio="true",__meta_kubernetes_service_labelpresent_istio_io_rev="true",__meta_kubernetes_service_labelpresent_operator_istio_io_component="true",__meta_kubernetes_service_labelpresent_release="true",__meta_kubernetes_service_name="istiod",__metrics_path__="/metrics",__scheme__="http",__scrape_interval__="30s",__scrape_timeout__="30s",job="istiod"}`,

  true,

  `{job="istiod",kubernetes_pod_name="istiod-5fd579679b-v2kck",service="pilot",tier="critical",business_unit="ctp",environment="prod",operating_cost_center="140018",prototype="s0",squad="traffic-management",stack_id="green"}`)
}