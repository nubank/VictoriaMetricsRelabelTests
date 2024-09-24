package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelJmxSparkDriver(t *testing.T) {
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
  `,

  `{__address__="10.40.31.147:9109",__meta_ec2_ami="ami-0864d46b555debb38",__meta_ec2_architecture="x86_64",__meta_ec2_availability_zone="us-east-1b",__meta_ec2_availability_zone_id="use1-az4",__meta_ec2_instance_id="i-094f87931d1571f01",__meta_ec2_instance_state="running",__meta_ec2_instance_type="r5d.2xlarge",__meta_ec2_owner_id="877163210394",__meta_ec2_primary_subnet_id="subnet-0239d42539db7a615",__meta_ec2_private_dns_name="ip-10-40-31-147.ec2.internal",__meta_ec2_private_ip="10.40.31.147",__meta_ec2_region="us-east-1",__meta_ec2_subnet_id=",subnet-0239d42539db7a615,",__meta_ec2_tag_NU_CLOTHO="default-generic-spark-mr",__meta_ec2_tag_NU_COMPONENT="driver",__meta_ec2_tag_NU_DIOL_NAME="0bec63c3-39a3-4007-a84f-b3455312e571",__meta_ec2_tag_NU_DIOL_NAMESPACE="itaipu-rt",__meta_ec2_tag_NU_ENV="prod",__meta_ec2_tag_NU_ITAIPU_JOB="itaipu-rt-0bec63c3-39a3-4007-a84f-b3455312e571",__meta_ec2_tag_NU_NAME="prod-foz-liquorice-mesos-fixed",__meta_ec2_tag_NU_OPERATING_COST_CENTER="110005",__meta_ec2_tag_NU_PROTOTYPE="foz",__meta_ec2_tag_NU_PROVIDER="standalone",__meta_ec2_tag_NU_SERVICE="mesos-fixed",__meta_ec2_tag_NU_SPOT_RATIO="0",__meta_ec2_tag_NU_SQUAD="data-infra",__meta_ec2_tag_NU_STACK="prod-foz",__meta_ec2_tag_NU_STACK_ID="prod-foz",__meta_ec2_tag_Name="prod-foz-liquorice-mesos-fixed",__meta_ec2_tag_REMEDIATED_AT="2024-05-19T17:52:25.624",__meta_ec2_vpc_id="vpc-00f2a5d5f6afcb6dc",__metrics_path__="/metrics",__scheme__="http",__scrape_interval__="35s",__scrape_timeout__="30s",job="jmx-spark-driver"}`,

  true,

  `{environment="prod",instance="10.40.31.147",instance_type="r5d.2xlarge",job="jmx-spark-driver",prototype="foz",service="mesos-fixed",squad="data-infra",stack_id="prod-foz"}`)
}