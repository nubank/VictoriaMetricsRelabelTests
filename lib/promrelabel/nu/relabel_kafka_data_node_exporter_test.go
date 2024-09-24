package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelKafkaDataNodeExporter(t *testing.T) {
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
      regex: kafka-data
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

  `{__address__="10.31.233.98:9100",__meta_ec2_ami="ami-0ed095c0b003f2d27",__meta_ec2_architecture="x86_64",__meta_ec2_availability_zone="us-east-2a",__meta_ec2_availability_zone_id="use2-az1",__meta_ec2_instance_id="i-08f707d44a6360024",__meta_ec2_instance_state="running",__meta_ec2_instance_type="m5.4xlarge",__meta_ec2_owner_id="540236187627",__meta_ec2_primary_subnet_id="subnet-097e79614ef5042cb",__meta_ec2_private_dns_name="ip-10-31-233-98.us-east-2.compute.internal",__meta_ec2_private_ip="10.31.233.98",__meta_ec2_region="us-east-2",__meta_ec2_subnet_id=",subnet-097e79614ef5042cb,",__meta_ec2_tag_NU_BUSINESS_UNIT="eng-horizontal",__meta_ec2_tag_NU_ENV="prod",__meta_ec2_tag_NU_OPERATING_COST_CENTER="140018",__meta_ec2_tag_NU_PROTOTYPE="global",__meta_ec2_tag_NU_SERVICE="kafka-data",__meta_ec2_tag_NU_SQUAD="messaging",__meta_ec2_tag_NU_STACK="prod-global-green-kafka-data",__meta_ec2_tag_NU_STACK_ID="green",__meta_ec2_tag_Name="prod-global-green-kafka-data",__meta_ec2_tag_REMEDIATED_AT="2024-07-17T10:07:42.580",__meta_ec2_tag_aws_autoscaling_groupName="prod-global-green-kafka-data-BlueGroup-1V90T43YTXY1F",__meta_ec2_tag_aws_cloudformation_logical_id="BlueGroup",__meta_ec2_tag_aws_cloudformation_stack_id="arn:aws:cloudformation:us-east-2:540236187627:stack/prod-global-green-kafka-data/ba1ae130-7885-11ec-85d6-0aa4824526ea",__meta_ec2_tag_aws_cloudformation_stack_name="prod-global-green-kafka-data",__meta_ec2_vpc_id="vpc-050edee8268e2e8f1",__metrics_path__="/metrics",__scheme__="http",__scrape_interval__="35s",__scrape_timeout__="30s",job="kafka-data-node-exporter"}`,

  true,

  `{job="kafka-data-node-exporter",environment="prod",instance="10.31.233.98",prototype="global",service="kafka-data",squad="messaging",stack_id="green",stack_name="prod-global-green-kafka-data"}`)
}