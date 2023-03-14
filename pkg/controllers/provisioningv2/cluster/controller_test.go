package cluster

import (
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	v1 "github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
	rkev1 "github.com/rancher/rancher/pkg/apis/rke.cattle.io/v1"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/generic/fake"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"
)

func TestRegexp(t *testing.T) {
	assert.True(t, mgmtNameRegexp.MatchString("local"))
	assert.False(t, mgmtNameRegexp.MatchString("alocal"))
	assert.False(t, mgmtNameRegexp.MatchString("localb"))
	assert.True(t, mgmtNameRegexp.MatchString("c-12345"))
	assert.False(t, mgmtNameRegexp.MatchString("ac-12345"))
	assert.False(t, mgmtNameRegexp.MatchString("c-12345b"))
	assert.False(t, mgmtNameRegexp.MatchString("ac-12345b"))
}

func TestController_generateProvisioningClusterFromLegacyCluster(t *testing.T) {
	tests := []struct {
		name    string
		cluster *v3.Cluster
	}{
		{
			name: "test-default",
			cluster: &v3.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "c-11111",
				},
				Spec: v3.ClusterSpec{
					ClusterSpecBase:    v3.ClusterSpecBase{},
					FleetWorkspaceName: "test-fleet-workspace-name",
				},
			},
		},
		{
			name: "test-cluster-agent-customization",
			cluster: &v3.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					Name: "c-22222",
				},
				Spec: v3.ClusterSpec{
					ClusterSpecBase: v3.ClusterSpecBase{
						ClusterAgentDeploymentCustomization: getTestClusterAgentCustomizationV3(),
						FleetAgentDeploymentCustomization:   getTestFleetAgentCustomizationV3(),
					},
					FleetWorkspaceName: "test-fleet-workspace-name",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler{}

			obj, _, err := h.generateProvisioningClusterFromLegacyCluster(tt.cluster, tt.cluster.Status)

			assert.Nil(t, err)
			assert.NotNil(t, obj, "Expected non-nil prov cluster obj")
			provCluster := obj[0].(*v1.Cluster)

			switch tt.name {
			case "test-default":
				assert.Nil(t, provCluster.Spec.ClusterAgentDeploymentCustomization)
				assert.Nil(t, provCluster.Spec.FleetAgentDeploymentCustomization)
			case "test-cluster-agent-customization":
				assert.Equal(t, getTestClusterAgentCustomizationV1(), provCluster.Spec.ClusterAgentDeploymentCustomization)
				assert.Equal(t, getTestFleetAgentCustomizationV1(), provCluster.Spec.FleetAgentDeploymentCustomization)
			}
		})
	}
}

func TestController_createNewCluster(t *testing.T) {
	tests := []struct {
		name        string
		cluster     *v1.Cluster
		clusterSpec v3.ClusterSpec
	}{
		{
			name: "test-default",
			cluster: &v1.Cluster{
				ObjectMeta: metav1.ObjectMeta{},
				Spec:       v1.ClusterSpec{},
			},
			clusterSpec: v3.ClusterSpec{},
		},
		{
			name: "test-cluster-agent-customization",
			cluster: &v1.Cluster{
				ObjectMeta: metav1.ObjectMeta{},
				Spec: v1.ClusterSpec{
					ClusterAgentDeploymentCustomization: getTestClusterAgentCustomizationV1(),
					FleetAgentDeploymentCustomization:   getTestFleetAgentCustomizationV1(),
				},
			},
			clusterSpec: v3.ClusterSpec{
				ClusterSpecBase: v3.ClusterSpecBase{},
			},
		},
	}
	mockCtrl := gomock.NewController(t)
	clusterCache := fake.NewMockNoNsCacheInterface[*v3.Cluster](mockCtrl)
	clusterCache.EXPECT().Get(gomock.Any()).Return(&v3.Cluster{}, nil).AnyTimes()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := handler{
				mgmtClusterCache: clusterCache,
			}

			obj, _, err := h.createNewCluster(tt.cluster, tt.cluster.Status, tt.clusterSpec)

			assert.Nil(t, err)
			assert.NotNil(t, obj, "Expected non-nil v3 cluster obj")
			jsonData, _ := json.Marshal(obj[0])
			var legacyCluster v3.Cluster
			json.Unmarshal(jsonData, &legacyCluster)

			switch tt.name {
			case "test-default":
				assert.Nil(t, legacyCluster.Spec.ClusterAgentDeploymentCustomization)
				assert.Nil(t, legacyCluster.Spec.FleetAgentDeploymentCustomization)
			case "test-cluster-agent-customization":
				assert.Equal(t, getTestClusterAgentCustomizationV3(), legacyCluster.Spec.ClusterAgentDeploymentCustomization)
				assert.Equal(t, getTestFleetAgentCustomizationV3(), legacyCluster.Spec.FleetAgentDeploymentCustomization)
			}
		})
	}
}

func getTestClusterAgentCustomizationV1() *v1.AgentDeploymentCustomization {
	return &v1.AgentDeploymentCustomization{
		AppendTolerations:            getTestClusterAgentToleration(),
		OverrideAffinity:             getTestClusterAgentAffinity(),
		OverrideResourceRequirements: getTestClusterAgentResourceReq(),
	}
}

func getTestClusterAgentCustomizationV3() *v3.AgentDeploymentCustomization {
	return &v3.AgentDeploymentCustomization{
		AppendTolerations:            getTestClusterAgentToleration(),
		OverrideAffinity:             getTestClusterAgentAffinity(),
		OverrideResourceRequirements: getTestClusterAgentResourceReq(),
	}
}

func getTestFleetAgentCustomizationV1() *v1.AgentDeploymentCustomization {
	return &v1.AgentDeploymentCustomization{
		AppendTolerations:            getTestFleetAgentToleration(),
		OverrideAffinity:             getTestFleetAgentAffinity(),
		OverrideResourceRequirements: getTestFleetAgentResourceReq(),
	}
}

func getTestFleetAgentCustomizationV3() *v3.AgentDeploymentCustomization {
	return &v3.AgentDeploymentCustomization{
		AppendTolerations:            getTestFleetAgentToleration(),
		OverrideAffinity:             getTestFleetAgentAffinity(),
		OverrideResourceRequirements: getTestFleetAgentResourceReq(),
	}
}

func getTestClusterAgentToleration() []corev1.Toleration {
	return []corev1.Toleration{{
		Effect: "NoSchedule",
		Key:    "node-role.kubernetes.io/controlplane-test",
		Value:  "true",
	},
	}
}

func getTestClusterAgentAffinity() *corev1.Affinity {
	return &corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
				{
					Weight: 1,
					Preference: corev1.NodeSelectorTerm{
						MatchExpressions: []corev1.NodeSelectorRequirement{
							{
								Key:      "cattle.io/cluster-agent-test",
								Operator: corev1.NodeSelectorOpIn,
								Values:   []string{"true"},
							},
						},
					},
				},
			},
		},
	}
}

func getTestClusterAgentResourceReq() *corev1.ResourceRequirements {
	return &corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU: resource.MustParse("500m"),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse("250Mi"),
		},
	}
}

func getTestFleetAgentToleration() []corev1.Toleration {
	return []corev1.Toleration{
		{
			Key:      "key",
			Operator: corev1.TolerationOpEqual,
			Value:    "value",
		},
	}
}

func getTestFleetAgentAffinity() *corev1.Affinity {
	return &corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
				{
					Weight: 1,
					Preference: corev1.NodeSelectorTerm{
						MatchExpressions: []corev1.NodeSelectorRequirement{
							{
								Key:      "fleet.cattle.io/agent",
								Operator: corev1.NodeSelectorOpIn,
								Values:   []string{"true"},
							},
						},
					},
				},
			},
		},
	}
}

func getTestFleetAgentResourceReq() *corev1.ResourceRequirements {
	return &corev1.ResourceRequirements{
		Limits: corev1.ResourceList{
			corev1.ResourceCPU: resource.MustParse("1"),
		},
		Requests: corev1.ResourceList{
			corev1.ResourceMemory: resource.MustParse("1Gi"),
		},
	}
}

func TestOnClusterRemove_CAPI_WithOwned(t *testing.T) {
	name := "test"
	namespace := "default"
	rancherCluster := createRancherCluster(name, namespace)
	capiCluster := createCAPICluster(name, namespace, rancherCluster)

	mockCtrl := gomock.NewController(t)

	clusterControllerMock := fake.NewMockControllerInterface[*v1.Cluster, *v1.ClusterList](mockCtrl)
	capiClusterCache := fake.NewMockCacheInterface[*capi.Cluster](mockCtrl)
	capiClusterClient := fake.NewMockClientInterface[*capi.Cluster, *capi.ClusterList](mockCtrl)
	rkeControlPlaneCache := fake.NewMockCacheInterface[*rkev1.RKEControlPlane](mockCtrl)

	capiClusterCache.EXPECT().Get(rancherCluster.Namespace, rancherCluster.Name).Return(capiCluster, nil)
	clusterControllerMock.EXPECT().UpdateStatus(gomock.AssignableToTypeOf(rancherCluster)).Return(rancherCluster, nil)
	capiClusterClient.EXPECT().Delete(rancherCluster.Namespace, rancherCluster.Name, &metav1.DeleteOptions{}).Return(nil)
	rkeControlPlaneCache.EXPECT().Get(rancherCluster.Namespace, rancherCluster.Name).Return(nil, nil)

	h := &handler{
		clusters:              clusterControllerMock,
		capiClustersCache:     capiClusterCache,
		capiClusters:          capiClusterClient,
		rkeControlPlanesCache: rkeControlPlaneCache,
	}
	_, err := h.OnClusterRemove("", rancherCluster)
	assert.Equal(t, generic.ErrSkip, err)
}

func TestOnClusterRemove_CAPI_NotOwned(t *testing.T) {
	name := "test"
	namespace := "default"

	rancherCluster := createRancherCluster(name, namespace)
	capiCluster := createCAPICluster(name, namespace, nil)

	mockCtrl := gomock.NewController(t)

	clusterControllerMock := fake.NewMockControllerInterface[*v1.Cluster, *v1.ClusterList](mockCtrl)
	capiClusterCache := fake.NewMockCacheInterface[*capi.Cluster](mockCtrl)

	capiClusterCache.EXPECT().Get(rancherCluster.Namespace, rancherCluster.Name).Return(capiCluster, nil)
	clusterControllerMock.EXPECT().UpdateStatus(gomock.AssignableToTypeOf(rancherCluster)).Return(rancherCluster, nil)

	h := &handler{
		clusters:          clusterControllerMock,
		capiClustersCache: capiClusterCache,
	}
	_, err := h.OnClusterRemove("", rancherCluster)
	assert.Nil(t, err)
}

func createRancherCluster(name, namespace string) *v1.Cluster {
	return &v1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec:   v1.ClusterSpec{},
		Status: v1.ClusterStatus{},
	}
}

func createCAPICluster(name, namespace string, ownedBy *v1.Cluster) *capi.Cluster {
	cluster := &capi.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec:   capi.ClusterSpec{},
		Status: capi.ClusterStatus{},
	}

	if ownedBy != nil {
		cluster.OwnerReferences = []metav1.OwnerReference{
			{
				APIVersion: v1.SchemeGroupVersion.Identifier(),
				Kind:       "Cluster",
				Name:       ownedBy.Name,
			},
		}
	}

	return cluster
}
