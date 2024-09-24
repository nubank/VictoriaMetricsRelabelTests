package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelBlackboxExporterGrafanaHealth(t *testing.T) {
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
    - source_labels: [__address__]
      target_label: __param_target
    - source_labels: [__param_target]
      target_label: endpoint
    - {replacement: staging-global-green-blackbox-exporter.monitoring:9115, target_label: __address__}
  `,

  `{__address__="https://prod-grafana.ist.nubank.world/api/health",__metrics_path__="/probe",__param_module="grafana_health",__scheme__="http",__scrape_interval__="30s",__scrape_timeout__="30s",job="blackbox-exporter-grafana-health"}`,

  true,

  `{job="blackbox-exporter-grafana-health",endpoint="https://prod-grafana.ist.nubank.world/api/health"}`)
}