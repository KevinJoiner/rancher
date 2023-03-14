package gke

import (
	"embed"
	stderrors "errors"
	"fmt"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/ghodss/yaml"
	"github.com/golang/mock/gomock"
	v1 "github.com/rancher/gke-operator/pkg/apis/gke.cattle.io/v1"
	apisv3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/rancher/rancher/pkg/controllers/management/clusteroperator"
	mgmtv3 "github.com/rancher/rancher/pkg/generated/norman/management.cattle.io/v3"
	"github.com/rancher/wrangler/pkg/generic/fake"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	secretv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/json"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

//go:embed test/*
var testFs embed.FS

type mockGkeOperatorController struct {
	gkeOperatorController
	mock.Mock
}

func getMockGkeOperatorController(t *testing.T, clusterState string) mockGkeOperatorController {
	t.Helper()
	ctrl := gomock.NewController(t)
	clusterMock := fake.NewMockNoNsClientInterface[*apisv3.Cluster, *apisv3.ClusterList](ctrl)
	clusterMock.EXPECT().Update(gomock.Any()).DoAndReturn(
		func(c *apisv3.Cluster) (*apisv3.Cluster, error) {
			return c, nil
		},
	).AnyTimes()

	var dynamicClient dynamic.NamespaceableResourceInterface

	switch clusterState {
	case "default":
		dynamicClient = MockNamespaceableResourceInterfaceDefault{}
	case "create":
		dynamicClient = MockNamespaceableResourceInterfaceCreate{}
	case "active":
		dynamicClient = MockNamespaceableResourceInterfaceActive{}
	case "update":
		dynamicClient = MockNamespaceableResourceInterfaceUpdate{}
	case "gkecc":
		dynamicClient = MockNamespaceableResourceInterfaceGkeCC{}
	default:
		dynamicClient = nil
	}

	return mockGkeOperatorController{
		gkeOperatorController: gkeOperatorController{
			OperatorController: clusteroperator.OperatorController{
				ClusterEnqueueAfter:  func(name string, duration time.Duration) {},
				SecretsCache:         nil,
				Secrets:              nil,
				TemplateCache:        nil,
				ProjectCache:         nil,
				AppLister:            nil,
				AppClient:            nil,
				NsClient:             nil,
				ClusterClient:        clusterMock,
				CatalogManager:       nil,
				SystemAccountManager: nil,
				DynamicClient:        dynamicClient,
				ClientDialer:         MockFactory{},
				Discovery:            MockDiscovery{},
			},
		},
		Mock: mock.Mock{},
	}
}

// test setInitialUpstreamSpec

func (m *mockGkeOperatorController) setInitialUpstreamSpec(cluster *mgmtv3.Cluster) (*mgmtv3.Cluster, error) {
	logrus.Infof("setting initial upstreamSpec on cluster [%s]", cluster.Name)

	// mock
	upstreamSpec := &v1.GKEClusterConfigSpec{}

	cluster = cluster.DeepCopy()
	cluster.Status.GKEStatus.UpstreamSpec = upstreamSpec
	return m.ClusterClient.Update(cluster)
}

// test generateAndSetServiceAccount with mock sibling func (getRestConfig)

func (m *mockGkeOperatorController) generateAndSetServiceAccount(cluster *mgmtv3.Cluster) (*mgmtv3.Cluster, error) {
	// mock
	m.Mock.On("getRestConfig", cluster).Return(&rest.Config{}, nil)

	restConfig, err := m.getRestConfig(cluster)
	if err != nil {
		return cluster, fmt.Errorf("error getting kube config: %v", err)
	}

	clusterDialer, err := m.ClientDialer.ClusterDialer(cluster.Name)
	if err != nil {
		return cluster, err
	}

	restConfig.Dial = clusterDialer
	cluster = cluster.DeepCopy()

	// mock
	secret := secretv1.Secret{}
	secret.Name = "cluster-serviceaccounttoken-sl7wm"

	if err != nil {
		return nil, err
	}
	cluster.Status.ServiceAccountTokenSecret = secret.Name
	cluster.Status.ServiceAccountToken = ""
	return m.ClusterClient.Update(cluster)
}

// test generateSATokenWithPublicAPI with mock sibling func (getRestConfig)

func (m *mockGkeOperatorController) generateSATokenWithPublicAPI(cluster *mgmtv3.Cluster) (string, *bool, error) {
	// mock
	m.Mock.On("getRestConfig", cluster).Return(&rest.Config{}, nil)

	restConfig, err := m.getRestConfig(cluster)
	if err != nil {
		return "", nil, err
	}
	requiresTunnel := new(bool)
	restConfig.Dial = (&net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}).DialContext

	// mock serviceToken
	serviceToken, err := "testtoken12345", nil

	if err != nil {
		*requiresTunnel = true
		var dnsError *net.DNSError
		if stderrors.As(err, &dnsError) && !dnsError.IsTemporary {
			return "", requiresTunnel, nil
		}
		var urlError *url.Error
		if stderrors.As(err, &urlError) && urlError.Timeout() {
			return "", requiresTunnel, nil
		}
		requiresTunnel = nil
	}
	return serviceToken, requiresTunnel, err
}

func (m *mockGkeOperatorController) getRestConfig(cluster *mgmtv3.Cluster) (*rest.Config, error) {

	_mc_ret := m.Called(cluster)

	var _r0 *rest.Config

	if _rfn, ok := _mc_ret.Get(0).(func(*mgmtv3.Cluster) *rest.Config); ok {
		_r0 = _rfn(cluster)
	} else {
		if _mc_ret.Get(0) != nil {
			_r0 = _mc_ret.Get(0).(*rest.Config)
		}
	}

	var _r1 error

	if _rfn, ok := _mc_ret.Get(1).(func(*mgmtv3.Cluster) error); ok {
		_r1 = _rfn(cluster)
	} else {
		_r1 = _mc_ret.Error(1)
	}

	return _r0, _r1

}

// utility

func getMockV3Cluster(filename string) (mgmtv3.Cluster, error) {
	var mockCluster mgmtv3.Cluster

	// Read the embedded file
	cluster, err := testFs.ReadFile(filename); if err != nil {
		return mockCluster, err
	}
	// Unmarshal cluster yaml into a management v3 cluster object
	err = yaml.Unmarshal(cluster, &mockCluster); if err != nil {
		return mockCluster, err
	}

	return mockCluster, nil
}

func getMockGkeClusterConfig(filename string) (*unstructured.Unstructured, error) {
	var gkeClusterConfig *unstructured.Unstructured

	// Read the embedded file
	bytes, err := testFs.ReadFile(filename); if err != nil {
		return gkeClusterConfig, err
	}
	// Unmarshal json into an unstructured cluster config object
	err = json.Unmarshal(bytes, &gkeClusterConfig); if err != nil {
		return gkeClusterConfig, err
	}

	return gkeClusterConfig, nil
}
