package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelKubernetesServiceEndpoints(t *testing.T) {
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
      replacement: staging
      source_labels: [__address__]
      target_label: environment
    - action: replace
      regex: .*
      replacement: global
      source_labels: [__address__]
      target_label: prototype
    - action: keep
      regex: true
      source_labels: [__meta_kubernetes_service_annotation_prometheus_io_scrape]
    - action: replace
      regex: (https?)
      source_labels: [__meta_kubernetes_service_annotation_prometheus_io_scheme]
      target_label: __scheme__
    - action: replace
      regex: (.+)
      source_labels: [__meta_kubernetes_service_annotation_prometheus_io_path]
      target_label: __metrics_path__
    - action: replace
      regex: ([^:]+)(?::\d+)?;(\d+)
      replacement: $1:$2
      source_labels: [__address__, __meta_kubernetes_service_annotation_prometheus_io_port]
      target_label: address
    - {action: labelmap, regex: __meta_kubernetes_service_label_(.+)}
    - action: replace
      source_labels: [__meta_kubernetes_namespace]
      target_label: kubernetes_namespace
    - action: replace
      source_labels: [__meta_kubernetes_service_name]
      target_label: kubernetes_name
    - {action: labeldrop, regex: container_id}
    - {action: labeldrop, regex: label_pod_template_hash}
  `,

  `{__address__="10.0.239.5:8181",__meta_kubernetes_endpoints_name="kube2iam-metrics",__meta_kubernetes_namespace="kube-system",__meta_kubernetes_pod_annotationpresent_scheduler_alpha_kubernetes_io_critical_pod="true",__meta_kubernetes_pod_container_image="193814090748.dkr.ecr.us-east-1.amazonaws.com/nu-kube2iam:16a6145f46da649ed3e7debb3ef464f86798c946@sha256:552d804cb259268d6799a4e7c555187043cbfebee07a9b2b91a3429b59162ab1",__meta_kubernetes_pod_container_name="kube2iam",__meta_kubernetes_pod_container_port_name="http",__meta_kubernetes_pod_container_port_number="8181",__meta_kubernetes_pod_container_port_protocol="TCP",__meta_kubernetes_pod_controller_kind="DaemonSet",__meta_kubernetes_pod_controller_name="kube2iam",__meta_kubernetes_pod_host_ip="10.0.239.5",__meta_kubernetes_pod_ip="10.0.239.5",__meta_kubernetes_pod_label_controller_revision_hash="6bd9796ddb",__meta_kubernetes_pod_label_name="kube2iam",__meta_kubernetes_pod_label_pod_template_generation="5",__meta_kubernetes_pod_labelpresent_controller_revision_hash="true",__meta_kubernetes_pod_labelpresent_name="true",__meta_kubernetes_pod_labelpresent_pod_template_generation="true",__meta_kubernetes_pod_name="kube2iam-hg72x",__meta_kubernetes_pod_node_name="ip-10-0-239-5.sa-east-1.compute.internal",__meta_kubernetes_pod_phase="Running",__meta_kubernetes_pod_ready="true",__meta_kubernetes_pod_uid="42ad7b98-e3f2-4767-a497-1019b07e03ff",__meta_kubernetes_service_annotation_kubectl_kubernetes_io_last_applied_configuration="{\"apiVersion\":\"v1\",\"kind\":\"Service\",\"metadata\":{\"annotations\":{\"prometheus.io/port\":\"8282\",\"prometheus.io/scrape\":\"true\"},\"name\":\"kube2iam-metrics\",\"namespace\":\"kube-system\"},\"spec\":{\"ports\":[{\"name\":\"http-metrics\",\"port\":8282,\"targetPort\":8282}],\"selector\":{\"name\":\"kube2iam\"}}}",__meta_kubernetes_service_annotation_prometheus_io_port="8282",__meta_kubernetes_service_annotation_prometheus_io_scrape="true",__meta_kubernetes_service_annotationpresent_kubectl_kubernetes_io_last_applied_configuration="true",__meta_kubernetes_service_annotationpresent_prometheus_io_port="true",__meta_kubernetes_service_annotationpresent_prometheus_io_scrape="true",__meta_kubernetes_service_name="kube2iam-metrics",__metrics_path__="/metrics",__scheme__="http",__scrape_interval__="30s",__scrape_timeout__="30s",job="kubernetes-service-endpoints"}`,

  true,

  `{job="kubernetes-service-endpoints",environment="staging",prototype="global",address="10.0.239.5:8282",kubernetes_namespace="kube-system",kubernetes_name="kube2iam-metrics"}`)
}