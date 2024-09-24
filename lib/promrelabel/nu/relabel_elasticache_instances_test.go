package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelElasticacheInstances(t *testing.T) {
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
    - replacement: $1
      separator: ':'
      source_labels: [__meta_elasticache_endpoint_address, __meta_elasticache_endpoint_port]
      target_label: instance
    - source_labels: [instance]
      target_label: __param_target
    - {action: labelmap, regex: __meta_elasticache_(.+)}
    - action: keep
      regex: staging:global
      separator: ':'
      source_labels: [tag_NU_ENV, tag_NU_PROTOTYPE]
    - source_labels: [tag_NU_PROTOTYPE]
      target_label: prototype
    - source_labels: [tag_NU_ENV]
      target_label: environment
    - source_labels: [tag_NU_STACK_ID]
      target_label: stack_id
    - source_labels: [tag_NU_SERVICE]
      target_label: service
    - source_labels: [tag_NU_SQUAD]
      target_label: squad
    - regex: redis
      replacement: staging-global-green-prometheus-redis-exporter:9121
      source_labels: [engine]
      target_label: __address__
    - regex: memcached
      replacement: staging-global-green-prometheus-memcached-exporter:9150
      source_labels: [engine]
      target_label: __address__
  `,

  `{__address__="undefined",__meta_elasticache_cache_cluster_id="staging-1-cryptofees-002",__meta_elasticache_cache_cluster_status="available",__meta_elasticache_cache_node_id="0001",__meta_elasticache_cache_node_status="available",__meta_elasticache_cache_node_type="cache.t4g.micro",__meta_elasticache_cache_parameter_group_name="default.redis6.x",__meta_elasticache_cache_subnet_group_name="priv-staging",__meta_elasticache_customer_availability_zone="us-east-1c",__meta_elasticache_endpoint_address="staging-1-cryptofees-002.wipw2x.0001.use1.cache.amazonaws.com",__meta_elasticache_endpoint_port="6379",__meta_elasticache_engine="redis",__meta_elasticache_engine_version="6.2.6",__meta_elasticache_preferred_availability_zone="us-east-1c",__meta_elasticache_replication_group_id="staging-1-cryptofees",__meta_elasticache_tag_NU_BUSINESS_UNIT="crypto",__meta_elasticache_tag_NU_ENV="staging",__meta_elasticache_tag_NU_OPERATING_COST_CENTER="120170",__meta_elasticache_tag_NU_PROTOTYPE="global",__meta_elasticache_tag_NU_SERVICE="crypto-fees-redis",__meta_elasticache_tag_NU_SQUAD="crypto",__meta_elasticache_tag_NU_STACK="staging-global-1-crypto-fees-redis",__meta_elasticache_tag_NU_STACK_ID="1",__meta_elasticache_tag_Name="staging-global-1-crypto-fees-redis",__meta_elasticache_tag_REMEDIATED_AT="2024-08-10T09:02:58.258",__meta_url="http://staging-global-green-prometheus-elasticache-sd:8080/elasticache.json",__metrics_path__="/scrape",__scheme__="http",__scrape_interval__="30s",__scrape_timeout__="30s",job="elasticache-instances"}`,

  true,

  `{instance="staging-1-cryptofees-002.wipw2x.0001.use1.cache.amazonaws.com:6379",cache_cluster_id="staging-1-cryptofees-002",cache_cluster_status="available",cache_node_id="0001",cache_node_status="available",cache_node_type="cache.t4g.micro",cache_parameter_group_name="default.redis6.x",cache_subnet_group_name="priv-staging",customer_availability_zone="us-east-1c",endpoint_address="staging-1-cryptofees-002.wipw2x.0001.use1.cache.amazonaws.com",endpoint_port="6379",engine="redis",engine_version="6.2.6",preferred_availability_zone="us-east-1c",replication_group_id="staging-1-cryptofees",tag_NU_BUSINESS_UNIT="crypto",tag_NU_ENV="staging",tag_NU_OPERATING_COST_CENTER="120170",tag_NU_PROTOTYPE="global",tag_NU_SERVICE="crypto-fees-redis",tag_NU_SQUAD="crypto",tag_NU_STACK="staging-global-1-crypto-fees-redis",tag_NU_STACK_ID="1",tag_Name="staging-global-1-crypto-fees-redis",tag_REMEDIATED_AT="2024-08-10T09:02:58.258",environment="staging",prototype="global",service="crypto-fees-redis",squad="crypto",stack_id="1",job="elasticache-instances"}`)
}