package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelStpNginxExporter(t *testing.T) {
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
      regex: stp-proxy
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

  `{__address__="10.1.240.213:9113",__meta_ec2_ami="ami-0528712befcd5d885",__meta_ec2_architecture="x86_64",__meta_ec2_availability_zone="us-west-1c",__meta_ec2_availability_zone_id="usw1-az1",__meta_ec2_instance_id="i-054dd65ae3b4c435e",__meta_ec2_instance_state="running",__meta_ec2_instance_type="t3.large",__meta_ec2_owner_id="540236187627",__meta_ec2_primary_subnet_id="subnet-0de2d0a3ba3dc402a",__meta_ec2_private_dns_name="ip-10-1-240-213.us-west-1.compute.internal",__meta_ec2_private_ip="10.1.240.213",__meta_ec2_region="us-west-1",__meta_ec2_subnet_id=",subnet-0de2d0a3ba3dc402a,",__meta_ec2_tag_NU_ENV="staging",__meta_ec2_tag_NU_PROTOTYPE="global",__meta_ec2_tag_NU_SERVICE="stp-proxy",__meta_ec2_tag_NU_SQUAD="payments-mx",__meta_ec2_tag_NU_STACK_ID="blue",__meta_ec2_tag_Name="staging-global-blue-stp-proxy",__meta_ec2_tag_aws_autoscaling_groupName="staging-global-blue-stp-proxy",__meta_ec2_tag_aws_ec2launchtemplate_id="lt-0d40f7cda92f80936",__meta_ec2_tag_aws_ec2launchtemplate_version="2",__meta_ec2_vpc_id="vpc-02a227defe2d97522",__metrics_path__="/metrics",__scheme__="http",__scrape_interval__="35s",__scrape_timeout__="30s",job="stp-nginx-exporter"}`,

  true,

  `{job="stp-nginx-exporter",environment="staging",instance="10.1.240.213",prototype="global",service="stp-proxy",squad="payments-mx",stack_id="blue"}`)
}