package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelStrimziResources(t *testing.T) {
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
    - {action: labelmap, regex: __meta_kubernetes_pod_label_(strimzi_io_.+), replacement: $1,
      separator: ;}
    - action: replace
      regex: (.*)
      replacement: $1
      separator: ;
      source_labels: [__meta_kubernetes_namespace]
      target_label: namespace
    - action: replace
      regex: (.*)
      replacement: $1
      separator: ;
      source_labels: [__meta_kubernetes_pod_name]
      target_label: kubernetes_pod_name
    - action: replace
      regex: (.*)
      replacement: $1
      separator: ;
      source_labels: [__meta_kubernetes_pod_node_name]
      target_label: node_name
    - action: replace
      regex: (.*)
      replacement: $1
      separator: ;
      source_labels: [__meta_kubernetes_pod_host_ip]
      target_label: node_ip
  `,

  `{__address__="localhost:9404",__metrics_path__="/metrics",__scheme__="http",__scrape_interval__="2m",__scrape_timeout__="2m",job="strimzi-resources"}`,

  true,

  `{job="strimzi-resources"}`)
}