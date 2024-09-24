package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelZookeeper(t *testing.T) {
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
      regex: zookeeper
      source_labels: [__meta_ec2_tag_NU_SERVICE]
    - action: keep
      regex: ^s0$
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

  `{__address__="10.0.202.253:9141",__meta_ec2_ami="ami-0e3f9a6ab03111832",__meta_ec2_architecture="x86_64",__meta_ec2_availability_zone="sa-east-1c",__meta_ec2_availability_zone_id="sae1-az3",__meta_ec2_instance_id="i-072335df82f6ca40c",__meta_ec2_instance_state="running",__meta_ec2_instance_type="m5.xlarge",__meta_ec2_owner_id="552767473918",__meta_ec2_primary_subnet_id="subnet-64b7a53d",__meta_ec2_private_dns_name="ip-10-0-202-253.sa-east-1.compute.internal",__meta_ec2_private_ip="10.0.202.253",__meta_ec2_region="sa-east-1",__meta_ec2_subnet_id=",subnet-64b7a53d,",__meta_ec2_tag_NU_ENV="prod",__meta_ec2_tag_NU_OPERATING_COST_CENTER="140018",__meta_ec2_tag_NU_PROTOTYPE="s0",__meta_ec2_tag_NU_SERVICE="zookeeper",__meta_ec2_tag_NU_SQUAD="traffic-management",__meta_ec2_tag_NU_STACK="prod-s0-blue-zookeeper-c",__meta_ec2_tag_NU_STACK_ID="blue",__meta_ec2_tag_Name="prod-s0-blue-zookeeper",__meta_ec2_tag_REMEDIATED_AT="2024-05-03T08:06:17.887",__meta_ec2_tag_aws_autoscaling_groupName="prod-s0-blue-zookeeper-AutoScalingGroupC-1DHF6WS57GV5X",__meta_ec2_tag_aws_cloudformation_logical_id="AutoScalingGroupC",__meta_ec2_tag_aws_cloudformation_stack_id="arn:aws:cloudformation:sa-east-1:552767473918:stack/prod-s0-blue-zookeeper/b092e230-106c-11ea-837b-0a997a645308",__meta_ec2_tag_aws_cloudformation_stack_name="prod-s0-blue-zookeeper",__meta_ec2_vpc_id="vpc-ca9778af",__metrics_path__="/metrics",__scheme__="http",__scrape_interval__="35s",__scrape_timeout__="35s",job="zookeeper"}`,

  true,

  `{job="zookeeper",environment="prod",instance="10.0.202.253",prototype="s0",service="zookeeper",squad="traffic-management",stack_id="blue",stack_name="prod-s0-blue-zookeeper-c"}`)
}