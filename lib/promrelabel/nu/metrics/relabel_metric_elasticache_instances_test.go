package promrelabel

import (
	"testing"
)


func TestRelabelElasticacheInstancesKeep(t *testing.T) {
  MetricsTest(t,

  `
    - action: labelkeep
      regex: (__.+|cache_cluster_id|cache_node_id|cmd|command|country|db|environment|err|instance|job|lru|maxmemory_policy|os|process_id|prototype|redis_build_id|redis_mode|redis_version|replication_group_id|role|run_id|service|slab|slave_ip|slave_port|slave_state|squad|stack_id|status|tcp_port|version)
  `,

  `redis_commands_rejected_calls_total{cmd="acl"}`,

  false,

  `redis_commands_rejected_calls_total{cmd="acl"}`)
}

func TestRelabelElasticacheInstancesDash(t *testing.T) {
  MetricsTest(t,

  `
    - action: labelkeep
      regex: (__.+|cache_cluster_id|cache_node_id|cmd|command|country|db|environment|err|instance|job|lru|maxmemory_policy|os|process_id|prototype|redis_build_id|redis_mode|redis_version|replication_group_id|role|run_id|service|slab|slave_ip|slave_port|slave_state|squad|stack_id|status|tcp_port|version)
  `,

  `foo{__bar__="bar"}`,

  false,

  `foo{__bar__="bar"}`)
}

func TestRelabelElasticacheInstancesDrop(t *testing.T) {
  MetricsTest(t,

  `
    - action: labelkeep
      regex: (__.+|cache_cluster_id|cache_node_id|cmd|command|country|db|environment|err|instance|job|lru|maxmemory_policy|os|process_id|prototype|redis_build_id|redis_mode|redis_version|replication_group_id|role|run_id|service|slab|slave_ip|slave_port|slave_state|squad|stack_id|status|tcp_port|version)
  `,

  `redis_commands_rejected_calls_total{cmd="acl",remover="remover"}`,

  false,

  `redis_commands_rejected_calls_total{cmd="acl"}`)
}