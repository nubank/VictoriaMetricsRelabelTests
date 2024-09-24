package promrelabel

import (
	"testing"
)


func TestRelabelKubernetesPodsJetty(t *testing.T) {
  MetricsTest(t,

  `
    - action: keep
      regex: ^(nu_crl|nauvoo_|jetty_).+
      source_labels: [__name__]
  `,

  `nu_crl_metric{foo="bar"}`,

  true,

  `nu_crl_metric{foo="bar"}`)
}

func TestRelabelKubernetesPodsJettyNauvoo(t *testing.T) {
  MetricsTest(t,

  `
    - action: keep
      regex: ^(nu_crl|nauvoo_|jetty_).+
      source_labels: [__name__]
  `,

  `nauvoo_metric{foo="bar"}`,

  true,

  `nauvoo_metric{foo="bar"}`)
}

func TestRelabelKubernetesPodsJettyDrop(t *testing.T) {
  MetricsTest(t,

  `
    - action: keep
      regex: ^(nu_crl|nauvoo_|jetty_).+
      source_labels: [__name__]
  `,

  `foo{foo="bar"}`,

  true,

  `{}`)
}