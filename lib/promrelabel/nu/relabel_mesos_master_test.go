package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelMesosMaster(t *testing.T) {
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

  `{__address__="10.40.37.222:9105",__meta_ec2_ami="ami-0800ccecb35a809d2",__meta_ec2_architecture="x86_64",__meta_ec2_availability_zone="us-east-1c",__meta_ec2_availability_zone_id="use1-az6",__meta_ec2_instance_id="i-0c5acd488049acf2f",__meta_ec2_instance_state="running",__meta_ec2_instance_type="m5.2xlarge",__meta_ec2_owner_id="877163210394",__meta_ec2_primary_subnet_id="subnet-06fd4566c6a541685",__meta_ec2_private_dns_name="ip-10-40-37-222.ec2.internal",__meta_ec2_private_ip="10.40.37.222",__meta_ec2_region="us-east-1",__meta_ec2_subnet_id=",subnet-06fd4566c6a541685,",__meta_ec2_tag_NU_APP="jetty",__meta_ec2_tag_NU_ENV="prod",__meta_ec2_tag_NU_LAYER="mesos-master",__meta_ec2_tag_NU_NAME="prod-foz-lead-mesos-master",__meta_ec2_tag_NU_OPERATING_COST_CENTER="110005",__meta_ec2_tag_NU_PROTOTYPE="foz",__meta_ec2_tag_NU_ROLE="api",__meta_ec2_tag_NU_SERVICE="mesos-master",__meta_ec2_tag_NU_SQUAD="data-infra",__meta_ec2_tag_NU_STACK="prod-lead",__meta_ec2_tag_Name="prod-foz-lead-mesos-master",__meta_ec2_tag_REMEDIATED_AT="2024-05-13T19:57:04.772",__meta_ec2_tag_aws_autoscaling_groupName="prod-foz-lead-mesos-master-BlueGroup-1CHHUY3Z3U17M",__meta_ec2_tag_aws_cloudformation_logical_id="BlueGroup",__meta_ec2_tag_aws_cloudformation_stack_id="arn:aws:cloudformation:us-east-1:877163210394:stack/prod-foz-lead-mesos-master/8e71e3c0-a2f1-11ed-be66-0e7349a98a8b",__meta_ec2_tag_aws_cloudformation_stack_name="prod-foz-lead-mesos-master",__meta_ec2_vpc_id="vpc-00f2a5d5f6afcb6dc",__metrics_path__="/metrics",__scheme__="http",__scrape_interval__="5m",__scrape_timeout__="4m",job="mesos-master"}`,

  true,

  `{job="mesos-master",environment="prod",instance="10.40.37.222",prototype="foz",service="mesos-master",squad="data-infra",stack_name="prod-lead"}`)
}