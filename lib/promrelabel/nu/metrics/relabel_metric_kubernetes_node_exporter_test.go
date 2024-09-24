package promrelabel

import (
	"testing"
)


func TestRelabelKubernetesNodeExporter(t *testing.T) {
  MetricsTest(t,

  `
    - action: replace
      regex: ^node_cpu.*_seconds_total;.*
      replacement: ''
      separator: ;
      source_labels: [__name__, failure_domain_beta_kubernetes_io_region]
      target_label: failure_domain_beta_kubernetes_io_region
    - action: replace
      regex: ^node_cpu.*_seconds_total;.*
      replacement: ''
      separator: ;
      source_labels: [__name__, topology_kubernetes_io_region]
      target_label: topology_kubernetes_io_region
    - action: replace
      regex: ^node_cpu.*_seconds_total;.*
      replacement: ''
      separator: ;
      source_labels: [__name__, failure_domain_beta_kubernetes_io_zone]
      target_label: failure_domain_beta_kubernetes_io_zone
    - action: replace
      regex: ^node_cpu.*_seconds_total;.*
      replacement: ''
      separator: ;
      source_labels: [__name__, topology_ebs_csi_aws_com_zone]
      target_label: topology_ebs_csi_aws_com_zone
    - action: replace
      regex: ^node_cpu.*_seconds_total;.*
      replacement: ''
      separator: ;
      source_labels: [__name__, vpc_amazonaws_com_eniConfig]
      target_label: vpc_amazonaws_com_eniConfig
    - action: replace
      regex: ^node_cpu.*_seconds_total;.*
      replacement: ''
      separator: ;
      source_labels: [__name__, nubank_com_br_eni_config]
      target_label: nubank_com_br_eni_config
    - action: replace
      regex: ^node_cpu.*_seconds_total;.*
      replacement: ''
      separator: ;
      source_labels: [__name__, beta_kubernetes_io_instance_type]
      target_label: beta_kubernetes_io_instance_type
    - action: replace
      regex: ^node_cpu.*_seconds_total;.*
      replacement: ''
      separator: ;
      source_labels: [__name__, beta_kubernetes_io_arch]
      target_label: beta_kubernetes_io_arch
    - action: replace
      regex: ^node_cpu.*_seconds_total;.*
      replacement: ''
      separator: ;
      source_labels: [__name__, kubernetes_io_hostname]
      target_label: kubernetes_io_hostname
    - action: drop
      regex: ^go_.+
      source_labels: [__name__]
  `,

  `node_cpu_guest_seconds_total{cpu="18",mode="user",failure_domain_beta_kubernetes_io_region="a",topology_kubernetes_io_region="b",failure_domain_beta_kubernetes_io_zone="c",topology_ebs_csi_aws_com_zone="d",vpc_amazonaws_com_eniConfig="e",nubank_com_br_eni_config= "f",beta_kubernetes_io_instance_type="g",beta_kubernetes_io_arch="h",kubernetes_io_hostname="i"}`,

  true,

  `node_cpu_guest_seconds_total{cpu="18",mode="user"}`)
}

func TestRelabelKubernetesNodeExporterDropGo(t *testing.T) {
  MetricsTest(t,

  `
    - action: replace
      regex: ^node_cpu.*_seconds_total;.*
      replacement: ''
      separator: ;
      source_labels: [__name__, failure_domain_beta_kubernetes_io_region]
      target_label: failure_domain_beta_kubernetes_io_region
    - action: replace
      regex: ^node_cpu.*_seconds_total;.*
      replacement: ''
      separator: ;
      source_labels: [__name__, topology_kubernetes_io_region]
      target_label: topology_kubernetes_io_region
    - action: replace
      regex: ^node_cpu.*_seconds_total;.*
      replacement: ''
      separator: ;
      source_labels: [__name__, failure_domain_beta_kubernetes_io_zone]
      target_label: failure_domain_beta_kubernetes_io_zone
    - action: replace
      regex: ^node_cpu.*_seconds_total;.*
      replacement: ''
      separator: ;
      source_labels: [__name__, topology_ebs_csi_aws_com_zone]
      target_label: topology_ebs_csi_aws_com_zone
    - action: replace
      regex: ^node_cpu.*_seconds_total;.*
      replacement: ''
      separator: ;
      source_labels: [__name__, vpc_amazonaws_com_eniConfig]
      target_label: vpc_amazonaws_com_eniConfig
    - action: replace
      regex: ^node_cpu.*_seconds_total;.*
      replacement: ''
      separator: ;
      source_labels: [__name__, nubank_com_br_eni_config]
      target_label: nubank_com_br_eni_config
    - action: replace
      regex: ^node_cpu.*_seconds_total;.*
      replacement: ''
      separator: ;
      source_labels: [__name__, beta_kubernetes_io_instance_type]
      target_label: beta_kubernetes_io_instance_type
    - action: replace
      regex: ^node_cpu.*_seconds_total;.*
      replacement: ''
      separator: ;
      source_labels: [__name__, beta_kubernetes_io_arch]
      target_label: beta_kubernetes_io_arch
    - action: replace
      regex: ^node_cpu.*_seconds_total;.*
      replacement: ''
      separator: ;
      source_labels: [__name__, kubernetes_io_hostname]
      target_label: kubernetes_io_hostname
    - action: drop
      regex: ^go_.+
      source_labels: [__name__]
  `,

  `go_gc_duration_seconds{quantile="0"}`,

  true,

  `{}`)
}