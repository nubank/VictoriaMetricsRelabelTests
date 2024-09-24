package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelJmxSparkExecutor(t *testing.T) {
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
    - source_labels: [__meta_ec2_instance_type]
      target_label: instance_type
    - source_labels: [__meta_ec2_tag_NU_STACK_ID]
      target_label: stack_id
    - source_labels: [__meta_ec2_tag_NU_SERVICE]
      target_label: service
    - source_labels: [__meta_ec2_tag_NU_SQUAD]
      target_label: squad
    - source_labels: [__meta_ec2_tag_SlaveType]
      target_label: itaipu_job
    - source_labels: [__meta_ec2_tag_NU_CLOTHO]
      target_label: queue
    - source_labels: [__meta_ec2_tag_NU_DATA_INFRA_TRANSACTION_ID]
      target_label: transaction_id
  `,

  `{__address__="10.174.236.124:9110",__meta_ec2_ami="ami-0700745fe9e31c13d",__meta_ec2_architecture="arm64",__meta_ec2_availability_zone="us-east-1d",__meta_ec2_availability_zone_id="use1-az1",__meta_ec2_instance_id="i-0e3abb4c36c58e479",__meta_ec2_instance_state="running",__meta_ec2_instance_type="r6gd.4xlarge",__meta_ec2_owner_id="877163210394",__meta_ec2_primary_subnet_id="subnet-0d41a9475e7866797",__meta_ec2_private_dns_name="ip-10-174-236-124.ec2.internal",__meta_ec2_private_ip="10.174.236.124",__meta_ec2_region="us-east-1",__meta_ec2_subnet_id=",subnet-0d41a9475e7866797,",__meta_ec2_tag_NU_CLOTHO="default-generic-spark-mr",__meta_ec2_tag_NU_COMPONENT="executor",__meta_ec2_tag_NU_COUNTRY="platform",__meta_ec2_tag_NU_DATA_CAPABILITY_EXECUTING="batch-compute-runtimes",__meta_ec2_tag_NU_DATA_CAPABILITY_REQUESTING="multi-repo",__meta_ec2_tag_NU_DIOL_NAME="646fe53d-9fa6-44b4-8c0f-83ce4b77e0aa",__meta_ec2_tag_NU_DIOL_NAMESPACE="itaipu-rt",__meta_ec2_tag_NU_ENV="prod",__meta_ec2_tag_NU_ITAIPU_JOB="itaipu-rt-646fe53d-9fa6-44b4-8c0f-83ce4b77e0aa",__meta_ec2_tag_NU_NAME="prod-foz-mesos-on-demand-itaipu-rt-646fe53d-9fa6-44b4-8c0f-83ce4b77e0aa",__meta_ec2_tag_NU_OPERATING_COST_CENTER="110005",__meta_ec2_tag_NU_PROTOTYPE="foz",__meta_ec2_tag_NU_PROVIDER="standalone",__meta_ec2_tag_NU_SERVICE="mesos-on-demand",__meta_ec2_tag_NU_SPOT_RATIO="0",__meta_ec2_tag_NU_SQUAD="data-infra",__meta_ec2_tag_NU_STACK="prod-foz",__meta_ec2_tag_NU_STACK_ID="prod-foz",__meta_ec2_tag_Name="prod-foz-mesos-on-demand-itaipu-rt-646fe53d-9fa6-44b4-8c0f-83ce4b77e0aa",__meta_ec2_tag_REMEDIATED_AT="2024-08-31T09:15:14.921",__meta_ec2_tag_aws_autoscaling_groupName="diol-itaipu-rt-646fe53d-9fa6-44b4-8c0f-83ce4b77e0aa",__meta_ec2_tag_aws_ec2_fleet_id="fleet-459d8f17-9c24-ccad-0618-892204d452e0",__meta_ec2_tag_aws_ec2launchtemplate_id="lt-05a26d341fe3868cf",__meta_ec2_tag_aws_ec2launchtemplate_version="1",__meta_ec2_vpc_id="vpc-00f2a5d5f6afcb6dc",__metrics_path__="/metrics",__scheme__="http",__scrape_interval__="35s",__scrape_timeout__="30s",job="jmx-spark-executor"}`,

  true,

  `{job="jmx-spark-executor",environment="prod",instance="10.174.236.124",instance_type="r6gd.4xlarge",prototype="foz",queue="default-generic-spark-mr",squad="data-infra",service="mesos-on-demand",stack_id="prod-foz"}`)
}