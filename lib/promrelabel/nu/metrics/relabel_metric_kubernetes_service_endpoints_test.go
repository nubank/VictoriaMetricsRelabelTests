package promrelabel

import (
	"testing"
)


func TestRelabelKubernetesServiceEndpoints(t *testing.T) {
  MetricsTest(t,

  `
    - {action: labelmap, regex: label_nubank_com_br_(.+)}
    - action: replace
      regex: ^$
      replacement: platform
      source_labels: [squad]
      target_label: squad
    - regex: (nu-)?(.+)
      replacement: $2
      source_labels: [container]
      target_label: service
    - action: drop
      regex: ^chaos_daemon_iptables_packets|chaos_daemon_iptables_packet_bytes$
      source_labels: [__name__]
    - action: replace
      regex: nu-datomic;staging-global-[a-zA-Z0-9]+-(.+)-datomic.*
      replacement: $${1}-datomic
      source_labels: [container, pod]
      target_label: service
  `,

  `metric_name{label_nubank_com_br_example="value",squad="",container="nu-foo"}`,

  true,

  `metric_name{label_nubank_com_br_example="value",example="value",squad="platform",service="foo",container="nu-foo"}`)
}

func TestRelabelKubernetesServiceEndpointsDropChaosDaemon(t *testing.T) {
  MetricsTest(t,

  `
    - {action: labelmap, regex: label_nubank_com_br_(.+)}
    - action: replace
      regex: ^$
      replacement: platform
      source_labels: [squad]
      target_label: squad
    - regex: (nu-)?(.+)
      replacement: $2
      source_labels: [container]
      target_label: service
    - action: drop
      regex: ^chaos_daemon_iptables_packets|chaos_daemon_iptables_packet_bytes$
      source_labels: [__name__]
    - action: replace
      regex: nu-datomic;staging-global-[a-zA-Z0-9]+-(.+)-datomic.*
      replacement: $${1}-datomic
      source_labels: [container, pod]
      target_label: service
  `,

  `chaos_daemon_iptables_packets{label_nubank_com_br_example="value",squad="",container="nu-foo"}`,

  true,

  `{}`)
}

func TestRelabelKubernetesServiceEndpointsReplace(t *testing.T) {
  MetricsTest(t,

  `
    - {action: labelmap, regex: label_nubank_com_br_(.+)}
    - action: replace
      regex: ^$
      replacement: platform
      source_labels: [squad]
      target_label: squad
    - regex: (nu-)?(.+)
      replacement: $2
      source_labels: [container]
      target_label: service
    - action: drop
      regex: ^chaos_daemon_iptables_packets|chaos_daemon_iptables_packet_bytes$
      source_labels: [__name__]
    - action: replace
      regex: nu-datomic;staging-global-[a-zA-Z0-9]+-(.+)-datomic.*
      replacement: $${1}-datomic
      source_labels: [container, pod]
      target_label: service
  `,

  `metric_name{label_nubank_com_br_example="value",squad="",container="nu-datomic",pod="staging-global-123-backend-datomic"}`,

  true,

  `metric_name{label_nubank_com_br_example="value",example="value",squad="platform",service="backend-datomic",container="nu-datomic",pod="staging-global-123-backend-datomic"}`)
}