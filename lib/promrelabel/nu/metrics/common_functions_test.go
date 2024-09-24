package promrelabel

import (
	"reflect"
	"strings"
	"testing"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)

func parseMetricString(metric string) (string, map[string]string) {
	// Separar o nome da métrica dos labels
	parts := strings.SplitN(metric, "{", 2)
	if len(parts) != 2 {
		return metric, nil
	}
	metricName := parts[0]
	labelsPart := strings.TrimSuffix(parts[1], "}")

	// Separar os labels em chave/valor e armazenar em um mapa
	labelsMap := make(map[string]string)
	labels := strings.Split(labelsPart, ",")
	for _, label := range labels {
		kv := strings.SplitN(label, "=", 2)
		if len(kv) == 2 {
			key := kv[0]
			value := strings.Trim(kv[1], `"`)
			labelsMap[key] = value
		}
	}
	return metricName, labelsMap
}

func CompareMetrics(metric1, metric2 string) bool {
	name1, labels1 := parseMetricString(metric1)
	name2, labels2 := parseMetricString(metric2)

	// Comparar o nome da métrica
	if name1 != name2 {
		return false
	}

	// Comparar os mapas de labels
	return reflect.DeepEqual(labels1, labels2)
}

func MetricsTest(t *testing.T, config, metric string, isFinalize bool, resultExpected string) {
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

	// Usar a função de comparação de métricas
	if !CompareMetrics(result, resultExpected) {
		t.Fatalf("unexpected result; got\n%s\nwant\n%s", result, resultExpected)
	}
}