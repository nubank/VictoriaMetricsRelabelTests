package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelAuroraNodeExporter(t *testing.T) {
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
      regex: aurora-scheduler
      source_labels: [__meta_ec2_tag_NU_SERVICE]
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

  `{__address__="10.40.10.2:9100",__meta_ec2_ami="ami-0800ccecb35a809d2",__meta_ec2_architecture="x86_64",__meta_ec2_availability_zone="us-east-1a",__meta_ec2_availability_zone_id="use1-az2",__meta_ec2_instance_id="i-07732a0ad45ac1a3a",__meta_ec2_instance_state="running",__meta_ec2_instance_type="m5.xlarge",__meta_ec2_owner_id="877163210394",__meta_ec2_primary_subnet_id="subnet-06b1d3acb95b8aa94",__meta_ec2_private_dns_name="ip-10-40-10-2.ec2.internal",__meta_ec2_private_ip="10.40.10.2",__meta_ec2_region="us-east-1",__meta_ec2_subnet_id=",subnet-06b1d3acb95b8aa94,",__meta_ec2_tag_NU_APP="jetty",__meta_ec2_tag_NU_ENV="prod",__meta_ec2_tag_NU_LAYER="aurora-scheduler",__meta_ec2_tag_NU_NAME="prod-foz-liquorice-aurora-scheduler",__meta_ec2_tag_NU_OPERATING_COST_CENTER="110005",__meta_ec2_tag_NU_PROTOTYPE="foz",__meta_ec2_tag_NU_ROLE="api",__meta_ec2_tag_NU_SERVICE="aurora-scheduler",__meta_ec2_tag_NU_SQUAD="data-infra",__meta_ec2_tag_NU_STACK="prod-liquorice",__meta_ec2_tag_Name="prod-foz-liquorice-aurora-scheduler",__meta_ec2_tag_REMEDIATED_AT="2024-05-13T20:01:52.337",__meta_ec2_tag_aws_autoscaling_groupName="prod-foz-liquorice-aurora-scheduler-BlueGroup-10B1IIBBOO9J0",__meta_ec2_tag_aws_cloudformation_logical_id="BlueGroup",__meta_ec2_tag_aws_cloudformation_stack_id="arn:aws:cloudformation:us-east-1:877163210394:stack/prod-foz-liquorice-aurora-scheduler/02646640-76f6-11ed-b7a4-1261d1718cb3",__meta_ec2_tag_aws_cloudformation_stack_name="prod-foz-liquorice-aurora-scheduler",__meta_ec2_vpc_id="vpc-00f2a5d5f6afcb6dc",__metrics_path__="/metrics",__scheme__="http",__scrape_interval__="35s",__scrape_timeout__="30s",job="aurora-node-exporter"}`,

  true,

  `{job="aurora-node-exporter",environment="prod",instance="10.40.10.2",prototype="foz",service="aurora-scheduler",squad="data-infra",stack_name="prod-liquorice"}`)
}