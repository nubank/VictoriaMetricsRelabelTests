package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelAkalaSpeiMonitor(t *testing.T) {
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
      regex: BANXICO_STATUS
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

  `{__address__="10.31.237.162:8080",__meta_ec2_ami="ami-089fe97bc00bff7cc",__meta_ec2_architecture="x86_64",__meta_ec2_availability_zone="us-east-2a",__meta_ec2_availability_zone_id="use2-az1",__meta_ec2_instance_id="i-0725c2365078228cb",__meta_ec2_instance_state="running",__meta_ec2_instance_type="t2.small",__meta_ec2_owner_id="540236187627",__meta_ec2_primary_subnet_id="subnet-097e79614ef5042cb",__meta_ec2_private_dns_name="ip-10-31-237-162.us-east-2.compute.internal",__meta_ec2_private_ip="10.31.237.162",__meta_ec2_region="us-east-2",__meta_ec2_subnet_id=",subnet-097e79614ef5042cb,",__meta_ec2_tag_NU_ENV="prod",__meta_ec2_tag_NU_PROTOTYPE="global",__meta_ec2_tag_NU_SERVICE="BANXICO_STATUS",__meta_ec2_tag_NU_SQUAD="spei-mx",__meta_ec2_tag_NU_STACK_ID="blue",__meta_ec2_tag_Name="prod-global-banxico-scraper",__meta_ec2_tag_aws_autoscaling_groupName="BANXICO-scraper",__meta_ec2_tag_aws_ec2_fleet_id="fleet-e4df8b55-7e88-ecf6-8c32-0308f18cbee8",__meta_ec2_tag_aws_ec2launchtemplate_id="lt-073b874079ac3cecd",__meta_ec2_tag_aws_ec2launchtemplate_version="3",__meta_ec2_vpc_id="vpc-050edee8268e2e8f1",__metrics_path__="/metrics",__scheme__="http",__scrape_interval__="35s",__scrape_timeout__="30s",job="akala-spei-monitor"}`,

  true,

  `{job="akala-spei-monitor",environment="prod",instance="10.31.237.162",prototype="global",service="BANXICO_STATUS",squad="spei-mx",stack_id="blue"}`)
}