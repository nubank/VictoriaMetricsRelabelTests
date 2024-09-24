package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelIstioIngress(t *testing.T) {
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
      regex: .*-envoy-prom
      source_labels: [__meta_kubernetes_pod_container_port_name]
    - action: replace
      source_labels: [__meta_kubernetes_pod_name]
      target_label: kubernetes_pod_name
    - action: replace
      source_labels: [__meta_kubernetes_pod_label_istio]
      target_label: istio
    - action: replace
      source_labels: [__meta_kubernetes_pod_label_app]
      target_label: service
    - replacement: critical
      source_labels: [__address__]
      target_label: tier
    - {action: labeldrop, regex: __meta_kubernetes_pod_label_nubank_com_br_name}
    - {action: labelmap, regex: __meta_kubernetes_pod_label_nubank_com_br_(.+)}
  `,

  `{__address__="100.73.122.106:15090",__meta_kubernetes_namespace="istio-system",__meta_kubernetes_pod_annotation_inject_istio_io_templates="gateway",__meta_kubernetes_pod_annotation_istio_io_rev="default",__meta_kubernetes_pod_annotation_prometheus_io_path="/stats/prometheus",__meta_kubernetes_pod_annotation_prometheus_io_port="15020",__meta_kubernetes_pod_annotation_prometheus_io_scheme="http",__meta_kubernetes_pod_annotation_prometheus_io_scrape="true",__meta_kubernetes_pod_annotation_proxy_istio_io_config="proxyMetadata:   ISTIO_META_OTEL_RESOURCE_ATTRIBUTES: service.name=$(NU_SERVICE),service.env=$(NU_ENV),service.instance.id=$(POD_NAME),k8s.pod.uid=$(POD_UID),nu.prototype=$(NU_PROTOTYPE),nu.aws_account_alias=$(NU_AWS_ACCOUNT_ALIAS),nu.aws_region=$(AWS_REGION),nu.country=$(NU_ORG_NAME),k8s.pod.private_ip=$(NU_PRIVATE_IP),nu.stack_id=$(NU_STACK_ID),nu.business_unit=$(NU_BUSINESS_UNIT),nu.operating_cost_center=$(NU_OPERATING_COST_CENTER),nu.squad=$(NU_SQUAD),nu.tier=$(NU_TIER),k8s.namespace=$(NU_K8S_NAMESPACE)proxyStatsMatcher:  inclusionRegexps:    - \".*ssl.*\"",__meta_kubernetes_pod_annotation_proxy_istio_io_overrides="{\"containers\":[{\"name\":\"istio-proxy\",\"image\":\"193814090748.dkr.ecr.us-east-1.amazonaws.com/mirror/istio/proxyv2:1.22.0\",\"ports\":[{\"name\":\"http-envoy-prom\",\"containerPort\":15090,\"protocol\":\"TCP\"}],\"resources\":{\"limits\":{\"cpu\":\"2\",\"memory\":\"1536Mi\"},\"requests\":{\"cpu\":\"100m\",\"memory\":\"1Gi\"}},\"volumeMounts\":[{\"name\":\"kube-api-access-6gw7t\",\"readOnly\":true,\"mountPath\":\"/var/run/secrets/kubernetes.io/serviceaccount\"}],\"terminationMessagePath\":\"/dev/termination-log\",\"terminationMessagePolicy\":\"File\",\"imagePullPolicy\":\"IfNotPresent\",\"securityContext\":{\"capabilities\":{\"drop\":[\"ALL\"]},\"privileged\":false,\"runAsUser\":1337,\"runAsGroup\":1337,\"runAsNonRoot\":true,\"readOnlyRootFilesystem\":true,\"allowPrivilegeEscalation\":false}}]}",__meta_kubernetes_pod_annotation_sidecar_istio_io_inject="true",__meta_kubernetes_pod_annotation_sidecar_istio_io_status="{\"initContainers\":null,\"containers\":[\"istio-proxy\"],\"volumes\":[\"workload-socket\",\"credential-socket\",\"workload-certs\",\"istio-envoy\",\"istio-data\",\"istio-podinfo\",\"istio-token\",\"istiod-ca-cert\"],\"imagePullSecrets\":null,\"revision\":\"default\"}",__meta_kubernetes_pod_annotationpresent_inject_istio_io_templates="true",__meta_kubernetes_pod_annotationpresent_istio_io_rev="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_path="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_port="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_scheme="true",__meta_kubernetes_pod_annotationpresent_prometheus_io_scrape="true",__meta_kubernetes_pod_annotationpresent_proxy_istio_io_config="true",__meta_kubernetes_pod_annotationpresent_proxy_istio_io_overrides="true",__meta_kubernetes_pod_annotationpresent_sidecar_istio_io_inject="true",__meta_kubernetes_pod_annotationpresent_sidecar_istio_io_status="true",__meta_kubernetes_pod_container_id="containerd://b0d81c08fd1e3e0744e5ec5c6469658fa192ac924f5525e3bdc2d3992ef3990e",__meta_kubernetes_pod_container_image="193814090748.dkr.ecr.us-east-1.amazonaws.com/mirror/istio/proxyv2:1.22.0",__meta_kubernetes_pod_container_init="false",__meta_kubernetes_pod_container_name="istio-proxy",__meta_kubernetes_pod_container_port_name="http-envoy-prom",__meta_kubernetes_pod_container_port_number="15090",__meta_kubernetes_pod_container_port_protocol="TCP",__meta_kubernetes_pod_controller_kind="ReplicaSet",__meta_kubernetes_pod_controller_name="public-ingress-d4d94dfdb",__meta_kubernetes_pod_host_ip="10.0.203.109",__meta_kubernetes_pod_ip="100.73.122.106",__meta_kubernetes_pod_label_app="public-ingress",__meta_kubernetes_pod_label_app_kubernetes_io_name="public-ingress",__meta_kubernetes_pod_label_app_kubernetes_io_version="1.22.0",__meta_kubernetes_pod_label_istio="ingress",__meta_kubernetes_pod_label_nubank_com_br_business_unit="ctp",__meta_kubernetes_pod_label_nubank_com_br_environment="prod",__meta_kubernetes_pod_label_nubank_com_br_name="public-ingress",__meta_kubernetes_pod_label_nubank_com_br_operating_cost_center="140018",__meta_kubernetes_pod_label_nubank_com_br_prototype="s0",__meta_kubernetes_pod_label_nubank_com_br_squad="traffic-management",__meta_kubernetes_pod_label_nubank_com_br_stack_id="green",__meta_kubernetes_pod_label_nubank_com_br_tier="critical",__meta_kubernetes_pod_label_pod_template_hash="d4d94dfdb",__meta_kubernetes_pod_label_service_istio_io_canonical_name="public-ingress",__meta_kubernetes_pod_label_service_istio_io_canonical_revision="1.22.0",__meta_kubernetes_pod_label_sidecar_istio_io_inject="true",__meta_kubernetes_pod_labelpresent_app="true",__meta_kubernetes_pod_labelpresent_app_kubernetes_io_name="true",__meta_kubernetes_pod_labelpresent_app_kubernetes_io_version="true",__meta_kubernetes_pod_labelpresent_istio="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_business_unit="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_environment="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_name="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_operating_cost_center="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_prototype="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_squad="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_stack_id="true",__meta_kubernetes_pod_labelpresent_nubank_com_br_tier="true",__meta_kubernetes_pod_labelpresent_pod_template_hash="true",__meta_kubernetes_pod_labelpresent_service_istio_io_canonical_name="true",__meta_kubernetes_pod_labelpresent_service_istio_io_canonical_revision="true",__meta_kubernetes_pod_labelpresent_sidecar_istio_io_inject="true",__meta_kubernetes_pod_name="public-ingress-d4d94dfdb-nfrtn",__meta_kubernetes_pod_node_name="ip-10-0-203-109.sa-east-1.compute.internal",__meta_kubernetes_pod_phase="Running",__meta_kubernetes_pod_ready="true",__meta_kubernetes_pod_uid="705b6cd1-c013-48df-99f3-48f66983ed3e",__metrics_path__="/stats/prometheus",__scheme__="http",__scrape_interval__="30s",__scrape_timeout__="30s",job="istio-ingress"}`,

  true,

  `{job="istio-ingress",kubernetes_pod_name="public-ingress-d4d94dfdb-nfrtn",istio="ingress",service="public-ingress",tier="critical",business_unit="ctp",environment="prod",operating_cost_center="140018",prototype="s0",squad="traffic-management",stack_id="green"}`)
}