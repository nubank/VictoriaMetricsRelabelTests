package promrelabel

import (
	"testing"
)


func TestRelabelIstioIngressCiphers(t *testing.T) {
  MetricsTest(t,

  `
    - regex: envoy_listener_ssl_ciphers_(.*)
      replacement: $${1}
      source_labels: [__name__]
      target_label: cipher
    - regex: envoy_listener_ssl_ciphers_(.*)
      replacement: envoy_listener_ssl_ciphers
      source_labels: [__name__]
      target_label: __name__
    - regex: envoy_listener_ssl_ciphers;.+
      replacement: ''
      source_labels: [__name__, listener_address]
      target_label: listener_address
    - regex: envoy_listener_ssl_versions_(.*)
      replacement: $${1}
      source_labels: [__name__]
      target_label: version
    - regex: envoy_listener_ssl_versions_(.*)
      replacement: envoy_listener_ssl_versions
      source_labels: [__name__]
      target_label: __name__
    - regex: envoy_listener_ssl_versions;.+
      replacement: ''
      source_labels: [__name__, listener_address]
      target_label: listener_address
  `,

  `envoy_listener_ssl_ciphers_ECDHE_RSA_AES128_GCM_SHA256{listener_address="0.0.0.0_443"}`,

  true,

  `envoy_listener_ssl_ciphers{cipher="ECDHE_RSA_AES128_GCM_SHA256"}`)
}

func TestRelabelIstioIngressVersions(t *testing.T) {
  MetricsTest(t,

  `
    - regex: envoy_listener_ssl_ciphers_(.*)
      replacement: $${1}
      source_labels: [__name__]
      target_label: cipher
    - regex: envoy_listener_ssl_ciphers_(.*)
      replacement: envoy_listener_ssl_ciphers
      source_labels: [__name__]
      target_label: __name__
    - regex: envoy_listener_ssl_ciphers;.+
      replacement: ''
      source_labels: [__name__, listener_address]
      target_label: listener_address
    - regex: envoy_listener_ssl_versions_(.*)
      replacement: $${1}
      source_labels: [__name__]
      target_label: version
    - regex: envoy_listener_ssl_versions_(.*)
      replacement: envoy_listener_ssl_versions
      source_labels: [__name__]
      target_label: __name__
    - regex: envoy_listener_ssl_versions;.+
      replacement: ''
      source_labels: [__name__, listener_address]
      target_label: listener_address
  `,

  `envoy_listener_ssl_versions_TLSv1_2{listener_address="0.0.0.0_443"}`,

  true,

  `envoy_listener_ssl_versions{version="TLSv1_2"}`)
}