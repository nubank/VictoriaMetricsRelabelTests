package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelKubernetesAwsVpcCni(t *testing.T) {
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
    - {action: labelmap, regex: __meta_kubernetes_node_label_(.+)}
    - {action: labeldrop, regex: failure_domain_beta_kubernetes_io_.*}
    - {action: labeldrop, regex: beta_kubernetes_io_.*}
    - {action: labeldrop, regex: kubernetes_io_os}
    - {action: labeldrop, regex: node_kubernetes_io_exclude_from_external_load_balancers}
    - {action: labeldrop, regex: spotinst_io_DrainingInitiated}
    - {action: labeldrop, regex: k8s_io_cloud_provider_aws}
    - regex: (.*):10250
      replacement: $${1}:61678
      source_labels: [__address__]
      target_label: __address__
  `,

  `{__address__="10.0.93.179:10250",__meta_kubernetes_node_address_Hostname="ip-10-0-93-179.sa-east-1.compute.internal",__meta_kubernetes_node_address_InternalDNS="ip-10-0-93-179.sa-east-1.compute.internal",__meta_kubernetes_node_address_InternalIP="10.0.93.179",__meta_kubernetes_node_annotation_alpha_kubernetes_io_provided_node_ip="10.0.93.179",__meta_kubernetes_node_annotation_csi_volume_kubernetes_io_nodeid="{\"csi.tigera.io\":\"ip-10-0-93-179.sa-east-1.compute.internal\",\"ebs.csi.aws.com\":\"i-09fcb7137d7a8931a\",\"efs.csi.aws.com\":\"i-09fcb7137d7a8931a\"}",__meta_kubernetes_node_annotation_karpenter_k8s_aws_ec2nodeclass_hash="17371524184729636746",__meta_kubernetes_node_annotation_karpenter_k8s_aws_ec2nodeclass_hash_version="v2",__meta_kubernetes_node_annotation_karpenter_sh_managed_by="prod-s0-green-kubernetes",__meta_kubernetes_node_annotation_karpenter_sh_nodepool_hash="10232932233304101654",__meta_kubernetes_node_annotation_karpenter_sh_nodepool_hash_version="v2",__meta_kubernetes_node_annotation_node_alpha_kubernetes_io_ttl="15",__meta_kubernetes_node_annotation_volumes_kubernetes_io_controller_managed_attach_detach="true",__meta_kubernetes_node_annotationpresent_alpha_kubernetes_io_provided_node_ip="true",__meta_kubernetes_node_annotationpresent_csi_volume_kubernetes_io_nodeid="true",__meta_kubernetes_node_annotationpresent_karpenter_k8s_aws_ec2nodeclass_hash="true",__meta_kubernetes_node_annotationpresent_karpenter_k8s_aws_ec2nodeclass_hash_version="true",__meta_kubernetes_node_annotationpresent_karpenter_sh_managed_by="true",__meta_kubernetes_node_annotationpresent_karpenter_sh_nodepool_hash="true",__meta_kubernetes_node_annotationpresent_karpenter_sh_nodepool_hash_version="true",__meta_kubernetes_node_annotationpresent_node_alpha_kubernetes_io_ttl="true",__meta_kubernetes_node_annotationpresent_volumes_kubernetes_io_controller_managed_attach_detach="true",__meta_kubernetes_node_label_beta_kubernetes_io_arch="amd64",__meta_kubernetes_node_label_beta_kubernetes_io_instance_type="m6i.16xlarge",__meta_kubernetes_node_label_beta_kubernetes_io_os="linux",__meta_kubernetes_node_label_failure_domain_beta_kubernetes_io_region="sa-east-1",__meta_kubernetes_node_label_failure_domain_beta_kubernetes_io_zone="sa-east-1b",__meta_kubernetes_node_label_k8s_io_cloud_provider_aws="7c1dba6327c948048786a5a5c3707470",__meta_kubernetes_node_label_karpenter_k8s_aws_instance_category="m",__meta_kubernetes_node_label_karpenter_k8s_aws_instance_cpu="64",__meta_kubernetes_node_label_karpenter_k8s_aws_instance_cpu_manufacturer="intel",__meta_kubernetes_node_label_karpenter_k8s_aws_instance_ebs_bandwidth="20000",__meta_kubernetes_node_label_karpenter_k8s_aws_instance_encryption_in_transit_supported="true",__meta_kubernetes_node_label_karpenter_k8s_aws_instance_family="m6i",__meta_kubernetes_node_label_karpenter_k8s_aws_instance_generation="6",__meta_kubernetes_node_label_karpenter_k8s_aws_instance_hypervisor="nitro",__meta_kubernetes_node_label_karpenter_k8s_aws_instance_memory="262144",__meta_kubernetes_node_label_karpenter_k8s_aws_instance_network_bandwidth="25000",__meta_kubernetes_node_label_karpenter_k8s_aws_instance_size="16xlarge",__meta_kubernetes_node_label_karpenter_sh_capacity_type="spot",__meta_kubernetes_node_label_karpenter_sh_initialized="true",__meta_kubernetes_node_label_karpenter_sh_nodepool="prod-s0-green-mixed",__meta_kubernetes_node_label_karpenter_sh_registered="true",__meta_kubernetes_node_label_kubernetes_io_arch="amd64",__meta_kubernetes_node_label_kubernetes_io_hostname="ip-10-0-93-179.sa-east-1.compute.internal",__meta_kubernetes_node_label_kubernetes_io_os="linux",__meta_kubernetes_node_label_node_kubernetes_io_instance_type="m6i.16xlarge",__meta_kubernetes_node_label_nubank_com_br_ami_id="ami-0dd84c63b7fe16daf",__meta_kubernetes_node_label_nubank_com_br_efs_shared_000="true",__meta_kubernetes_node_label_nubank_com_br_eni_config="prod-sa-east-1b",__meta_kubernetes_node_label_nubank_com_br_instance_id="i-09fcb7137d7a8931a",__meta_kubernetes_node_label_nubank_com_br_local_instance_storage="true",__meta_kubernetes_node_label_nubank_com_br_subnet="prod",__meta_kubernetes_node_label_pool_lifecycle="mixed",__meta_kubernetes_node_label_subnet="prod",__meta_kubernetes_node_label_topology_ebs_csi_aws_com_zone="sa-east-1b",__meta_kubernetes_node_label_topology_k8s_aws_zone_id="sae1-az2",__meta_kubernetes_node_label_topology_kubernetes_io_region="sa-east-1",__meta_kubernetes_node_label_topology_kubernetes_io_zone="sa-east-1b",__meta_kubernetes_node_labelpresent_beta_kubernetes_io_arch="true",__meta_kubernetes_node_labelpresent_beta_kubernetes_io_instance_type="true",__meta_kubernetes_node_labelpresent_beta_kubernetes_io_os="true",__meta_kubernetes_node_labelpresent_failure_domain_beta_kubernetes_io_region="true",__meta_kubernetes_node_labelpresent_failure_domain_beta_kubernetes_io_zone="true",__meta_kubernetes_node_labelpresent_k8s_io_cloud_provider_aws="true",__meta_kubernetes_node_labelpresent_karpenter_k8s_aws_instance_category="true",__meta_kubernetes_node_labelpresent_karpenter_k8s_aws_instance_cpu="true",__meta_kubernetes_node_labelpresent_karpenter_k8s_aws_instance_cpu_manufacturer="true",__meta_kubernetes_node_labelpresent_karpenter_k8s_aws_instance_ebs_bandwidth="true",__meta_kubernetes_node_labelpresent_karpenter_k8s_aws_instance_encryption_in_transit_supported="true",__meta_kubernetes_node_labelpresent_karpenter_k8s_aws_instance_family="true",__meta_kubernetes_node_labelpresent_karpenter_k8s_aws_instance_generation="true",__meta_kubernetes_node_labelpresent_karpenter_k8s_aws_instance_hypervisor="true",__meta_kubernetes_node_labelpresent_karpenter_k8s_aws_instance_memory="true",__meta_kubernetes_node_labelpresent_karpenter_k8s_aws_instance_network_bandwidth="true",__meta_kubernetes_node_labelpresent_karpenter_k8s_aws_instance_size="true",__meta_kubernetes_node_labelpresent_karpenter_sh_capacity_type="true",__meta_kubernetes_node_labelpresent_karpenter_sh_initialized="true",__meta_kubernetes_node_labelpresent_karpenter_sh_nodepool="true",__meta_kubernetes_node_labelpresent_karpenter_sh_registered="true",__meta_kubernetes_node_labelpresent_kubernetes_io_arch="true",__meta_kubernetes_node_labelpresent_kubernetes_io_hostname="true",__meta_kubernetes_node_labelpresent_kubernetes_io_os="true",__meta_kubernetes_node_labelpresent_node_kubernetes_io_instance_type="true",__meta_kubernetes_node_labelpresent_nubank_com_br_ami_id="true",__meta_kubernetes_node_labelpresent_nubank_com_br_efs_shared_000="true",__meta_kubernetes_node_labelpresent_nubank_com_br_eni_config="true",__meta_kubernetes_node_labelpresent_nubank_com_br_instance_id="true",__meta_kubernetes_node_labelpresent_nubank_com_br_local_instance_storage="true",__meta_kubernetes_node_labelpresent_nubank_com_br_subnet="true",__meta_kubernetes_node_labelpresent_pool_lifecycle="true",__meta_kubernetes_node_labelpresent_subnet="true",__meta_kubernetes_node_labelpresent_topology_ebs_csi_aws_com_zone="true",__meta_kubernetes_node_labelpresent_topology_k8s_aws_zone_id="true",__meta_kubernetes_node_labelpresent_topology_kubernetes_io_region="true",__meta_kubernetes_node_labelpresent_topology_kubernetes_io_zone="true",__meta_kubernetes_node_name="ip-10-0-93-179.sa-east-1.compute.internal",__meta_kubernetes_node_provider_id="aws:///sa-east-1b/i-09fcb7137d7a8931a",__metrics_path__="/metrics",__scheme__="http",__scrape_interval__="30s",__scrape_timeout__="30s",instance="ip-10-0-93-179.sa-east-1.compute.internal",job="kubernetes-aws-vpc-cni"}`,

  true,

  `{address="10.0.93.179:61678",instance="ip-10-0-93-179.sa-east-1.compute.internal",job="kubernetes-aws-vpc-cni",karpenter_k8s_aws_instance_category="m",karpenter_k8s_aws_instance_cpu="64",karpenter_k8s_aws_instance_cpu_manufacturer="intel",karpenter_k8s_aws_instance_ebs_bandwidth="20000",karpenter_k8s_aws_instance_encryption_in_transit_supported="true",karpenter_k8s_aws_instance_family="m6i",karpenter_k8s_aws_instance_generation="6",karpenter_k8s_aws_instance_hypervisor="nitro",karpenter_k8s_aws_instance_memory="262144",karpenter_k8s_aws_instance_network_bandwidth="25000",karpenter_k8s_aws_instance_size="16xlarge",karpenter_sh_capacity_type="spot",karpenter_sh_initialized="true",karpenter_sh_nodepool="prod-s0-green-mixed",karpenter_sh_registered="true",kubernetes_io_arch="amd64",kubernetes_io_hostname="ip-10-0-93-179.sa-east-1.compute.internal",node_kubernetes_io_instance_type="m6i.16xlarge",nubank_com_br_ami_id="ami-0dd84c63b7fe16daf",nubank_com_br_efs_shared_000="true",nubank_com_br_eni_config="prod-sa-east-1b",nubank_com_br_instance_id="i-09fcb7137d7a8931a",nubank_com_br_local_instance_storage="true",nubank_com_br_subnet="prod",pool_lifecycle="mixed",subnet="prod",topology_ebs_csi_aws_com_zone="sa-east-1b",topology_k8s_aws_zone_id="sae1-az2",topology_kubernetes_io_region="sa-east-1",topology_kubernetes_io_zone="sa-east-1b"}`)
}