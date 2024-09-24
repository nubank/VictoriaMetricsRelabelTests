package promrelabel

import (
	"testing"

	"reflect"
	"sort"
	"strings"

	pr "github.com/VictoriaMetrics/VictoriaMetrics/lib/promrelabel"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/promutils"
)


func TestRelabelKubernetesKubelets(t *testing.T) {
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
    - action: replace
      regex: .*
      replacement: staging
      source_labels: [__address__]
      target_label: environment
    - action: replace
      regex: .*
      replacement: global
      source_labels: [__address__]
      target_label: prototype
    - {action: labelmap, regex: __meta_kubernetes_node_label_(.+)}
    - {action: labeldrop, regex: kubernetes_io_os}
    - {action: labeldrop, regex: node_kubernetes_io_exclude_from_external_load_balancers}
    - {action: labeldrop, regex: spotinst_io_DrainingInitiated}
    - {action: labeldrop, regex: k8s_io_cloud_provider_aws}
    - {replacement: kubernetes.default.svc:443, target_label: address}
    - regex: (.+)
      replacement: /api/v1/nodes/$${1}/proxy/metrics
      source_labels: [__meta_kubernetes_node_name]
      target_label: metrics_path
  `,

  `{__address__="10.0.193.138:10250",__meta_kubernetes_node_address_Hostname="ip-10-0-193-138.sa-east-1.compute.internal",__meta_kubernetes_node_address_InternalDNS="ip-10-0-193-138.sa-east-1.compute.internal",__meta_kubernetes_node_address_InternalIP="10.0.193.138",__meta_kubernetes_node_annotation_alpha_kubernetes_io_provided_node_ip="10.0.193.138",__meta_kubernetes_node_annotation_csi_volume_kubernetes_io_nodeid="{\"ebs.csi.aws.com\":\"i-0efaf57eb8c0d205d\",\"efs.csi.aws.com\":\"i-0efaf57eb8c0d205d\"}",__meta_kubernetes_node_annotation_karpenter_k8s_aws_ec2nodeclass_hash="2360970719142568193",__meta_kubernetes_node_annotation_karpenter_k8s_aws_ec2nodeclass_hash_version="v2",__meta_kubernetes_node_annotation_karpenter_sh_managed_by="prod-s0-green-kubernetes",__meta_kubernetes_node_annotation_karpenter_sh_nodepool_hash="1059986104871095930",__meta_kubernetes_node_annotation_karpenter_sh_nodepool_hash_version="v2",__meta_kubernetes_node_annotation_node_alpha_kubernetes_io_ttl="15",__meta_kubernetes_node_annotation_volumes_kubernetes_io_controller_managed_attach_detach="true",__meta_kubernetes_node_annotationpresent_alpha_kubernetes_io_provided_node_ip="true",__meta_kubernetes_node_annotationpresent_csi_volume_kubernetes_io_nodeid="true",__meta_kubernetes_node_annotationpresent_karpenter_k8s_aws_ec2nodeclass_hash="true",__meta_kubernetes_node_annotationpresent_karpenter_k8s_aws_ec2nodeclass_hash_version="true",__meta_kubernetes_node_annotationpresent_karpenter_sh_managed_by="true",__meta_kubernetes_node_annotationpresent_karpenter_sh_nodepool_hash="true",__meta_kubernetes_node_annotationpresent_karpenter_sh_nodepool_hash_version="true",__meta_kubernetes_node_annotationpresent_node_alpha_kubernetes_io_ttl="true",__meta_kubernetes_node_annotationpresent_volumes_kubernetes_io_controller_managed_attach_detach="true",__meta_kubernetes_node_label_beta_kubernetes_io_arch="arm64",__meta_kubernetes_node_label_beta_kubernetes_io_instance_type="m6g.8xlarge",__meta_kubernetes_node_label_beta_kubernetes_io_os="linux",__meta_kubernetes_node_label_failure_domain_beta_kubernetes_io_region="sa-east-1",__meta_kubernetes_node_label_failure_domain_beta_kubernetes_io_zone="sa-east-1c",__meta_kubernetes_node_label_k8s_io_cloud_provider_aws="7c1dba6327c948048786a5a5c3707470",__meta_kubernetes_node_label_karpenter_k8s_aws_instance_category="m",__meta_kubernetes_node_label_karpenter_k8s_aws_instance_cpu="32",__meta_kubernetes_node_label_karpenter_k8s_aws_instance_cpu_manufacturer="aws",__meta_kubernetes_node_label_karpenter_k8s_aws_instance_ebs_bandwidth="9500",__meta_kubernetes_node_label_karpenter_k8s_aws_instance_encryption_in_transit_supported="false",__meta_kubernetes_node_label_karpenter_k8s_aws_instance_family="m6g",__meta_kubernetes_node_label_karpenter_k8s_aws_instance_generation="6",__meta_kubernetes_node_label_karpenter_k8s_aws_instance_hypervisor="nitro",__meta_kubernetes_node_label_karpenter_k8s_aws_instance_memory="131072",__meta_kubernetes_node_label_karpenter_k8s_aws_instance_network_bandwidth="12000",__meta_kubernetes_node_label_karpenter_k8s_aws_instance_size="8xlarge",__meta_kubernetes_node_label_karpenter_sh_capacity_type="on-demand",__meta_kubernetes_node_label_karpenter_sh_initialized="true",__meta_kubernetes_node_label_karpenter_sh_nodepool="prod-s0-green-od-only-arm",__meta_kubernetes_node_label_karpenter_sh_registered="true",__meta_kubernetes_node_label_kubernetes_io_arch="arm64",__meta_kubernetes_node_label_kubernetes_io_hostname="ip-10-0-193-138.sa-east-1.compute.internal",__meta_kubernetes_node_label_kubernetes_io_os="linux",__meta_kubernetes_node_label_node_kubernetes_io_instance_type="m6g.8xlarge",__meta_kubernetes_node_label_nubank_com_br_ami_id="ami-05d7bea580de645a5",__meta_kubernetes_node_label_nubank_com_br_architecture="arm",__meta_kubernetes_node_label_nubank_com_br_efs_shared_000="true",__meta_kubernetes_node_label_nubank_com_br_eni_config="prod-sa-east-1c",__meta_kubernetes_node_label_nubank_com_br_instance_id="i-0efaf57eb8c0d205d",__meta_kubernetes_node_label_nubank_com_br_local_instance_storage="true",__meta_kubernetes_node_label_nubank_com_br_subnet="prod",__meta_kubernetes_node_label_pool_lifecycle="od-only-arm",__meta_kubernetes_node_label_subnet="prod",__meta_kubernetes_node_label_topology_ebs_csi_aws_com_zone="sa-east-1c",__meta_kubernetes_node_label_topology_k8s_aws_zone_id="sae1-az3",__meta_kubernetes_node_label_topology_kubernetes_io_region="sa-east-1",__meta_kubernetes_node_label_topology_kubernetes_io_zone="sa-east-1c",__meta_kubernetes_node_labelpresent_beta_kubernetes_io_arch="true",__meta_kubernetes_node_labelpresent_beta_kubernetes_io_instance_type="true",__meta_kubernetes_node_labelpresent_beta_kubernetes_io_os="true",__meta_kubernetes_node_labelpresent_failure_domain_beta_kubernetes_io_region="true",__meta_kubernetes_node_labelpresent_failure_domain_beta_kubernetes_io_zone="true",__meta_kubernetes_node_labelpresent_k8s_io_cloud_provider_aws="true",__meta_kubernetes_node_labelpresent_karpenter_k8s_aws_instance_category="true",__meta_kubernetes_node_labelpresent_karpenter_k8s_aws_instance_cpu="true",__meta_kubernetes_node_labelpresent_karpenter_k8s_aws_instance_cpu_manufacturer="true",__meta_kubernetes_node_labelpresent_karpenter_k8s_aws_instance_ebs_bandwidth="true",__meta_kubernetes_node_labelpresent_karpenter_k8s_aws_instance_encryption_in_transit_supported="true",__meta_kubernetes_node_labelpresent_karpenter_k8s_aws_instance_family="true",__meta_kubernetes_node_labelpresent_karpenter_k8s_aws_instance_generation="true",__meta_kubernetes_node_labelpresent_karpenter_k8s_aws_instance_hypervisor="true",__meta_kubernetes_node_labelpresent_karpenter_k8s_aws_instance_memory="true",__meta_kubernetes_node_labelpresent_karpenter_k8s_aws_instance_network_bandwidth="true",__meta_kubernetes_node_labelpresent_karpenter_k8s_aws_instance_size="true",__meta_kubernetes_node_labelpresent_karpenter_sh_capacity_type="true",__meta_kubernetes_node_labelpresent_karpenter_sh_initialized="true",__meta_kubernetes_node_labelpresent_karpenter_sh_nodepool="true",__meta_kubernetes_node_labelpresent_karpenter_sh_registered="true",__meta_kubernetes_node_labelpresent_kubernetes_io_arch="true",__meta_kubernetes_node_labelpresent_kubernetes_io_hostname="true",__meta_kubernetes_node_labelpresent_kubernetes_io_os="true",__meta_kubernetes_node_labelpresent_node_kubernetes_io_instance_type="true",__meta_kubernetes_node_labelpresent_nubank_com_br_ami_id="true",__meta_kubernetes_node_labelpresent_nubank_com_br_architecture="true",__meta_kubernetes_node_labelpresent_nubank_com_br_efs_shared_000="true",__meta_kubernetes_node_labelpresent_nubank_com_br_eni_config="true",__meta_kubernetes_node_labelpresent_nubank_com_br_instance_id="true",__meta_kubernetes_node_labelpresent_nubank_com_br_local_instance_storage="true",__meta_kubernetes_node_labelpresent_nubank_com_br_subnet="true",__meta_kubernetes_node_labelpresent_pool_lifecycle="true",__meta_kubernetes_node_labelpresent_subnet="true",__meta_kubernetes_node_labelpresent_topology_ebs_csi_aws_com_zone="true",__meta_kubernetes_node_labelpresent_topology_k8s_aws_zone_id="true",__meta_kubernetes_node_labelpresent_topology_kubernetes_io_region="true",__meta_kubernetes_node_labelpresent_topology_kubernetes_io_zone="true",__meta_kubernetes_node_name="ip-10-0-193-138.sa-east-1.compute.internal",__meta_kubernetes_node_provider_id="aws:///sa-east-1c/i-0efaf57eb8c0d205d",__metrics_path__="/metrics",__scheme__="https",__scrape_interval__="30s",__scrape_timeout__="30s",instance="ip-10-0-193-138.sa-east-1.compute.internal",job="kubernetes-kubelets"}`,

  true,

  `{job="kubernetes-kubelets",environment="staging",prototype="global",beta_kubernetes_io_arch="arm64",beta_kubernetes_io_instance_type="m6g.8xlarge",beta_kubernetes_io_os="linux",failure_domain_beta_kubernetes_io_region="sa-east-1",failure_domain_beta_kubernetes_io_zone="sa-east-1c",karpenter_k8s_aws_instance_category="m",karpenter_k8s_aws_instance_cpu="32",karpenter_k8s_aws_instance_cpu_manufacturer="aws",karpenter_k8s_aws_instance_ebs_bandwidth="9500",karpenter_k8s_aws_instance_encryption_in_transit_supported="false",karpenter_k8s_aws_instance_family="m6g",karpenter_k8s_aws_instance_generation="6",karpenter_k8s_aws_instance_hypervisor="nitro",karpenter_k8s_aws_instance_memory="131072",karpenter_k8s_aws_instance_network_bandwidth="12000",karpenter_k8s_aws_instance_size="8xlarge",karpenter_sh_capacity_type="on-demand",karpenter_sh_initialized="true",karpenter_sh_nodepool="prod-s0-green-od-only-arm",karpenter_sh_registered="true",kubernetes_io_arch="arm64",kubernetes_io_hostname="ip-10-0-193-138.sa-east-1.compute.internal",node_kubernetes_io_instance_type="m6g.8xlarge",nubank_com_br_ami_id="ami-05d7bea580de645a5",nubank_com_br_architecture="arm",nubank_com_br_efs_shared_000="true",nubank_com_br_eni_config="prod-sa-east-1c",nubank_com_br_instance_id="i-0efaf57eb8c0d205d",nubank_com_br_local_instance_storage="true",nubank_com_br_subnet="prod",pool_lifecycle="od-only-arm",subnet="prod",topology_ebs_csi_aws_com_zone="sa-east-1c",topology_k8s_aws_zone_id="sae1-az3",topology_kubernetes_io_region="sa-east-1",topology_kubernetes_io_zone="sa-east-1c",address="kubernetes.default.svc:443",metrics_path="/api/v1/nodes/ip-10-0-193-138.sa-east-1.compute.internal/proxy/metrics",instance="ip-10-0-193-138.sa-east-1.compute.internal"}`)
}