package rke1

import (
	"strconv"
	"time"

	"github.com/rancher/norman/types"
	"github.com/rancher/rancher/tests/framework/clients/rancher"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
	nodestat "github.com/rancher/rancher/tests/framework/extensions/nodes"
	"github.com/sirupsen/logrus"
	kwait "k8s.io/apimachinery/pkg/util/wait"
)

const (
	active = "active"
)

type NodeRoles struct {
	ControlPlane bool  `json:"controlplane,omitempty" yaml:"controlplane,omitempty"`
	Etcd         bool  `json:"etcd,omitempty" yaml:"etcd,omitempty"`
	Worker       bool  `json:"worker,omitempty" yaml:"worker,omitempty"`
	Quantity     int64 `json:"quantity" yaml:"quantity"`
}

// NodePoolSetup is a helper method that will loop and setup muliple node pools with the defined node roles from the `nodeRoles` parameter
// `nodeRoles` would be in this format
//
//	  []map[string]bool{
//	  {
//		   ControlPlane: true,
//		   Etcd:         false,
//		   Worker:       false,
//		   Quantity:     1,
//	  },
//	  {
//		   ControlPlane: false,
//		   Etcd:         true,
//		   Worker:       false,
//		   Quantity:     1,
//	  },
//	 }
func NodePoolSetup(client *rancher.Client, nodeRoles []NodeRoles, ClusterID, NodeTemplateID string) (*management.NodePool, error) {
	nodePoolConfig := management.NodePool{
		ClusterID:               ClusterID,
		DeleteNotReadyAfterSecs: 0,
		NodeTemplateID:          NodeTemplateID,
	}

	for index, roles := range nodeRoles {
		nodePoolConfig.ControlPlane = roles.ControlPlane
		nodePoolConfig.Etcd = roles.Etcd
		nodePoolConfig.Worker = roles.Worker
		nodePoolConfig.Quantity = roles.Quantity
		nodePoolConfig.HostnamePrefix = "auto-rke1-" + strconv.Itoa(index) + ClusterID

		_, err := client.Management.NodePool.Create(&nodePoolConfig)

		if err != nil {
			return nil, err
		}
	}

	return &nodePoolConfig, nil
}

// MatchRKE1NodeRoles is a helper method that will return the desired node in the cluster, based on the node role.
func MatchRKE1NodeRoles(client *rancher.Client, cluster *management.Cluster, nodeRoles NodeRoles) (*management.Node, error) {
	nodes, err := client.Management.Node.ListAll(&types.ListOpts{
		Filters: map[string]interface{}{
			"clusterId": cluster.ID,
		},
	})
	if err != nil {
		return nil, err
	}

	for _, node := range nodes.Data {
		if nodeRoles.ControlPlane != node.ControlPlane {
			continue
		}
		if nodeRoles.Etcd != node.Etcd {
			continue
		}
		if nodeRoles.Worker != node.Worker {
			continue
		}

		return &node, nil
	}

	return nil, nil
}

// updateNodePoolQuantity is a helper method that will update the node pool with the desired quantity.
func updateNodePoolQuantity(client *rancher.Client, cluster *management.Cluster, node *management.Node, nodeRoles NodeRoles) (*management.NodePool, error) {
	updatedNodePool, err := client.Management.NodePool.ByID(node.NodePoolID)
	if err != nil {
		return nil, err
	}

	updatedNodePool.Quantity += nodeRoles.Quantity

	logrus.Infof("Scaling the machine pool to %v total nodes", updatedNodePool.Quantity)
	_, err = client.Management.NodePool.Update(updatedNodePool, &updatedNodePool)
	if err != nil {
		return nil, err
	}

	err = kwait.Poll(500*time.Millisecond, 10*time.Minute, func() (done bool, err error) {
		clusterResp, err := client.Management.Cluster.ByID(cluster.ID)
		if err != nil {
			return false, err
		}

		if clusterResp.State == active && nodestat.AllManagementNodeReady(client, clusterResp.ID) == nil {
			logrus.Infof("Node pool is scaled!")
			return true, nil
		} else {
			return false, nil
		}
	})
	if err != nil {
		return nil, err
	}

	return updatedNodePool, nil
}

// ScaleNodePoolNodes is a helper method that will add a new node pool to the existing RKE1 cluster, based on the nodeRoles
// configuration. Once done, it will scale the node pool, scale it back down and then delete the node pool.
func ScaleNodePoolNodes(client *rancher.Client, cluster *management.Cluster, node *management.Node, nodeRoles NodeRoles) (*management.NodePool, error) {
	updatedNodePool, err := updateNodePoolQuantity(client, cluster, node, nodeRoles)
	if err != nil {
		return nil, err
	}

	logrus.Infof("Node pool has been scaled!")

	return updatedNodePool, nil
}
