package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelItaipuSparkDrivers(t *testing.T) {
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
    - source_labels: [__meta_ec2_tag_NU_ENV]
      target_label: environment
    - source_labels: [__meta_ec2_tag_NU_SERVICE]
      target_label: service
    - source_labels: [__meta_ec2_tag_NU_ITAIPU_JOB]
      target_label: itaipu_job
    - source_labels: [__meta_ec2_tag_NU_SQUAD]
      target_label: squad
  `,

  `{__address__="10.40.49.225:32003",__meta_ec2_ami="ami-0864d46b555debb38",__meta_ec2_architecture="x86_64",__meta_ec2_availability_zone="us-east-1d",__meta_ec2_availability_zone_id="use1-az1",__meta_ec2_instance_id="i-072f5a85b4281ed0e",__meta_ec2_instance_state="running",__meta_ec2_instance_type="r5d.2xlarge",__meta_ec2_owner_id="877163210394",__meta_ec2_primary_subnet_id="subnet-02eac723490b26c53",__meta_ec2_private_dns_name="ip-10-40-49-225.ec2.internal",__meta_ec2_private_ip="10.40.49.225",__meta_ec2_region="us-east-1",__meta_ec2_subnet_id=",subnet-02eac723490b26c53,",__meta_ec2_tag_NU_CLOTHO="default",__meta_ec2_tag_NU_COMPONENT="driver",__meta_ec2_tag_NU_DATA_INFRA_JOB_CRITICALITY="standard",__meta_ec2_tag_NU_DATA_INFRA_TRANSACTION_ID="fd98df96-0304-5a73-8aa6-1b5ef94d56d9",__meta_ec2_tag_NU_DIOL_NAME="db2e1da3-c81f-4c1a-9145-168cb2345cfa",__meta_ec2_tag_NU_DIOL_NAMESPACE="itaipu-rt",__meta_ec2_tag_NU_ENV="prod",__meta_ec2_tag_NU_ITAIPU_JOB="itaipu-runtime-db2e1da3-c81f-4c1a-9145-168cb2345cfa",__meta_ec2_tag_NU_NAME="prod-foz-liquorice-mesos-fixed",__meta_ec2_tag_NU_OPERATING_COST_CENTER="110005",__meta_ec2_tag_NU_PROTOTYPE="foz",__meta_ec2_tag_NU_PROVIDER="standalone",__meta_ec2_tag_NU_SERVICE="mesos-fixed",__meta_ec2_tag_NU_SPOT_RATIO="100",__meta_ec2_tag_NU_SQUAD="data-infra",__meta_ec2_tag_NU_STACK="prod-foz",__meta_ec2_tag_NU_STACK_ID="prod-foz",__meta_ec2_tag_Name="prod-foz-liquorice-mesos-fixed",__meta_ec2_tag_REMEDIATED_AT="2024-05-18T18:24:28.907",__meta_ec2_vpc_id="vpc-00f2a5d5f6afcb6dc",__metrics_path__="/metrics/prometheus",__scheme__="http",__scrape_interval__="2m5s",__scrape_timeout__="1m",job="itaipu-spark-drivers"}`,

  true,

  `{job="itaipu-spark-drivers",environment="prod",service="mesos-fixed",itaipu_job="itaipu-runtime-db2e1da3-c81f-4c1a-9145-168cb2345cfa",squad="data-infra"}`)
}