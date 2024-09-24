package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelPseNginxExporter(t *testing.T) {
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
      regex: pse-proxy
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
    - source_labels: [__meta_ec2_tag_NU_STACK_ID]
      target_label: stack_id
  `,

  `{__address__="10.127.120.122:9113",__meta_ec2_ami="ami-03130878b60947df3",__meta_ec2_architecture="x86_64",__meta_ec2_availability_zone="us-west-1c",__meta_ec2_availability_zone_id="usw1-az1",__meta_ec2_instance_id="i-0f442e6a75b83b627",__meta_ec2_instance_state="running",__meta_ec2_instance_type="t3.micro",__meta_ec2_owner_id="986808460310",__meta_ec2_primary_subnet_id="subnet-0c4d8e8950c86ae06",__meta_ec2_private_dns_name="ip-10-127-120-122.us-west-1.compute.internal",__meta_ec2_private_ip="10.127.120.122",__meta_ec2_region="us-west-1",__meta_ec2_subnet_id=",subnet-0c4d8e8950c86ae06,",__meta_ec2_tag_NU_ENV="staging",__meta_ec2_tag_NU_OPERATING_COST_CENTER="140018",__meta_ec2_tag_NU_PROTOTYPE="global",__meta_ec2_tag_NU_SERVICE="pse-proxy",__meta_ec2_tag_NU_SQUAD="orchestration",__meta_ec2_tag_NU_STACK="green",__meta_ec2_tag_NU_STACK_ID="green",__meta_ec2_tag_Name="staging-global-green-pse-proxy",__meta_ec2_tag_REMEDIATED_AT="2024-07-17T08:04:41.496",__meta_ec2_tag_aws_ec2launchtemplate_id="lt-09aeb9df061d517a7",__meta_ec2_tag_aws_ec2launchtemplate_version="8",__meta_ec2_vpc_id="vpc-048342060eaf88885",__metrics_path__="/metrics",__scheme__="http",__scrape_interval__="35s",__scrape_timeout__="30s",job="pse-nginx-exporter"}`,

  true,

  `{job="pse-nginx-exporter",instance="10.127.120.122",environment="staging",prototype="global",service="pse-proxy",stack_name="green",squad="orchestration",stack_id="green"}`)
}