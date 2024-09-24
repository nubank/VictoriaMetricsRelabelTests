package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelAirflow(t *testing.T) {
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
    - source_labels: [__meta_ec2_tag_NU_STACK_ID]
      target_label: stack_id
    - source_labels: [__meta_ec2_tag_NU_STACK_ID]
      target_label: layer
    - source_labels: [__meta_ec2_tag_NU_SERVICE]
      target_label: service
    - source_labels: [__meta_ec2_tag_NU_SQUAD]
      target_label: squad
  `,

  `{__address__="10.40.18.13:8081",__meta_ec2_ami="ami-0800ccecb35a809d2",__meta_ec2_architecture="x86_64",__meta_ec2_availability_zone="us-east-1b",__meta_ec2_availability_zone_id="use1-az4",__meta_ec2_instance_id="i-0a0a44854f22d0862",__meta_ec2_instance_state="running",__meta_ec2_instance_type="m5.2xlarge",__meta_ec2_owner_id="877163210394",__meta_ec2_primary_subnet_id="subnet-0239d42539db7a615",__meta_ec2_private_dns_name="ip-10-40-18-13.ec2.internal",__meta_ec2_private_ip="10.40.18.13",__meta_ec2_region="us-east-1",__meta_ec2_subnet_id=",subnet-0239d42539db7a615,",__meta_ec2_tag_NU_APP="jetty",__meta_ec2_tag_NU_ENV="prod",__meta_ec2_tag_NU_LAYER="airflow",__meta_ec2_tag_NU_NAME="prod-foz-green-airflow",__meta_ec2_tag_NU_PROTOTYPE="foz",__meta_ec2_tag_NU_ROLE="api",__meta_ec2_tag_NU_SERVICE="airflow",__meta_ec2_tag_NU_SQUAD="data-infra",__meta_ec2_tag_NU_STACK="prod-green",__meta_ec2_tag_Name="prod-foz-green-airflow",__meta_ec2_tag_aws_autoscaling_groupName="prod-foz-green-airflow-BlueGroup-13CS6GHROSSYD",__meta_ec2_tag_aws_cloudformation_logical_id="BlueGroup",__meta_ec2_tag_aws_cloudformation_stack_id="arn:aws:cloudformation:us-east-1:877163210394:stack/prod-foz-green-airflow/0aec4d30-adf3-11ed-8141-0a5f1c044ead",__meta_ec2_tag_aws_cloudformation_stack_name="prod-foz-green-airflow",__meta_ec2_vpc_id="vpc-00f2a5d5f6afcb6dc",__metrics_path__="/admin/metrics/",__scheme__="http",__scrape_interval__="2m",__scrape_timeout__="1m",job="airflow"}`,

  true,

  `{environment="prod",instance="10.40.18.13",job="airflow",prototype="foz",service="airflow",squad="data-infra"}`)
}