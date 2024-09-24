package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelFireblocksCosigner(t *testing.T) {
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
      regex: fireblocks-cosigner-infra
      source_labels: [__meta_ec2_tag_NU_SERVICE]
    - action: keep
      regex: ^global$
      source_labels: [__meta_ec2_tag_NU_PROTOTYPE]
    - source_labels: [__meta_ec2_tag_NU_PROTOTYPE]
      target_label: prototype
    - source_labels: [__meta_ec2_private_ip]
      target_label: instance
    - source_labels: [__meta_ec2_tag_NU_ENV]
      target_label: environment
    - source_labels: [__meta_ec2_tag_NU_STACK]
      target_label: stack_name
    - source_labels: [__meta_ec2_tag_NU_SERVICE]
      target_label: service
    - source_labels: [__meta_ec2_tag_NU_SQUAD]
      target_label: squad

  `,

  `{__address__="10.139.206.188:9100",__meta_ec2_ami="ami-09538990a0c4fe9be",__meta_ec2_architecture="x86_64",__meta_ec2_availability_zone="us-east-1d",__meta_ec2_availability_zone_id="use1-az4",__meta_ec2_instance_id="i-0c732a3433c749e15",__meta_ec2_instance_state="running",__meta_ec2_instance_type="c5a.xlarge",__meta_ec2_owner_id="679868755743",__meta_ec2_primary_subnet_id="subnet-0231766f80908e701",__meta_ec2_private_dns_name="ip-10-139-206-188.ec2.internal",__meta_ec2_private_ip="10.139.206.188",__meta_ec2_region="us-east-1",__meta_ec2_subnet_id=",subnet-0231766f80908e701,",__meta_ec2_tag_NU_ENV="staging",__meta_ec2_tag_NU_OPERATING_COST_CENTER="120170",__meta_ec2_tag_NU_PROTOTYPE="global",__meta_ec2_tag_NU_SERVICE="fireblocks-cosigner-infra",__meta_ec2_tag_NU_SQUAD="crypto",__meta_ec2_tag_NU_STACK="br-staging-staging-green-fireblocks-cosigner-infra",__meta_ec2_tag_Name="br-staging-staging-green-fireblocks-cosigner-infra-ec2-mari",__meta_ec2_vpc_id="vpc-0646504720826d54b",__metrics_path__="/metrics",__scheme__="http",__scrape_interval__="30s",__scrape_timeout__="30s",job="fireblocks-cosigner"}`,

  true,

  `{job="fireblocks-cosigner",instance="10.139.206.188",environment="staging",prototype="global",service="fireblocks-cosigner-infra",squad="crypto",stack_name="br-staging-staging-green-fireblocks-cosigner-infra"}`)
}