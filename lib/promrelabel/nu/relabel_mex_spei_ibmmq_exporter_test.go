package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelMexSpeiIbmmqExporter(t *testing.T) {
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
      regex: mex-spei-ibmmq
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

  `{__address__="10.32.254.39:9113",__meta_ec2_ami="ami-0b0c3c984874d1b89",__meta_ec2_architecture="x86_64",__meta_ec2_availability_zone="us-east-2c",__meta_ec2_availability_zone_id="use2-az3",__meta_ec2_instance_id="i-0328e34b25f926f69",__meta_ec2_instance_state="running",__meta_ec2_instance_type="m5.xlarge",__meta_ec2_owner_id="540236187627",__meta_ec2_primary_subnet_id="subnet-0fade33ec6e6cebde",__meta_ec2_private_dns_name="ip-10-32-254-39.us-east-2.compute.internal",__meta_ec2_private_ip="10.32.254.39",__meta_ec2_region="us-east-2",__meta_ec2_subnet_id=",subnet-0fade33ec6e6cebde,",__meta_ec2_tag_NU_BUSINESS_UNIT="mex-payments",__meta_ec2_tag_NU_ENV="prod",__meta_ec2_tag_NU_PROTOTYPE="global",__meta_ec2_tag_NU_SERVICE="mex-spei-ibmmq",__meta_ec2_tag_NU_SQUAD="mex-payments",__meta_ec2_tag_Name="prod-spei-mx-ibm-mq-blue",__meta_ec2_tag_aws_autoscaling_groupName="prod-spei-mx-ibm-mq-blue-AutoScalingGroup-5KEBE450DYWK",__meta_ec2_tag_aws_cloudformation_logical_id="AutoScalingGroup",__meta_ec2_tag_aws_cloudformation_stack_id="arn:aws:cloudformation:us-east-2:540236187627:stack/prod-spei-mx-ibm-mq-blue/4c8e7d30-25bc-11ee-98c3-0a2e4ddd1993",__meta_ec2_tag_aws_cloudformation_stack_name="prod-spei-mx-ibm-mq-blue",__meta_ec2_tag_aws_ec2launchtemplate_id="lt-04145d3be0481e11d",__meta_ec2_tag_aws_ec2launchtemplate_version="1",__meta_ec2_vpc_id="vpc-0a42f02c02f1d2a72",__metrics_path__="/metrics",__scheme__="http",__scrape_interval__="1m",__scrape_timeout__="30s",job="mex-spei-ibmmq-exporter"}`,

  true,

  `{job="mex-spei-ibmmq-exporter",environment="prod",instance="10.32.254.39",prototype="global",service="mex-spei-ibmmq",squad="mex-payments"}`)
}