package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelKubernetesApiservers(t *testing.T) {
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
      regex: default;kubernetes;https
      source_labels: [__meta_kubernetes_namespace, __meta_kubernetes_service_name, __meta_kubernetes_endpoint_port_name]
  `,

  `{__address__="10.0.229.134:443",__meta_kubernetes_endpoint_port_name="https",__meta_kubernetes_endpoint_port_protocol="TCP",__meta_kubernetes_endpoint_ready="true",__meta_kubernetes_endpoints_label_endpointslice_kubernetes_io_skip_mirror="true",__meta_kubernetes_endpoints_labelpresent_endpointslice_kubernetes_io_skip_mirror="true",__meta_kubernetes_endpoints_name="kubernetes",__meta_kubernetes_namespace="default",__meta_kubernetes_service_label_component="apiserver",__meta_kubernetes_service_label_provider="kubernetes",__meta_kubernetes_service_labelpresent_component="true",__meta_kubernetes_service_labelpresent_provider="true",__meta_kubernetes_service_name="kubernetes",__metrics_path__="/metrics",__scheme__="https",__scrape_interval__="30s",__scrape_timeout__="30s",job="kubernetes-apiservers"}`,

  true,

  `{job="kubernetes-apiservers",environment="staging",prototype="global"}`)
}