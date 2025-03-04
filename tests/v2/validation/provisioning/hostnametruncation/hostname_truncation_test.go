package hostnametruncation

import (
	"fmt"
	"testing"

	"github.com/rancher/rancher/tests/framework/clients/rancher"
	"github.com/rancher/rancher/tests/framework/extensions/clusters"
	"github.com/rancher/rancher/tests/framework/extensions/machinepools"
	"github.com/rancher/rancher/tests/framework/extensions/provisioning"
	"github.com/rancher/rancher/tests/framework/extensions/provisioninginput"
	"github.com/rancher/rancher/tests/framework/pkg/config"
	namegen "github.com/rancher/rancher/tests/framework/pkg/namegenerator"
	"github.com/rancher/rancher/tests/framework/pkg/session"
	"github.com/rancher/rancher/tests/v2/validation/provisioning/permutations"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HostnameTruncationTestSuite struct {
	suite.Suite
	client         *rancher.Client
	session        *session.Session
	clustersConfig *provisioninginput.Config
}

func (r *HostnameTruncationTestSuite) TearDownSuite() {
	r.session.Cleanup()
}

func (r *HostnameTruncationTestSuite) SetupSuite() {
	testSession := session.NewSession()
	r.session = testSession

	r.clustersConfig = new(provisioninginput.Config)
	config.LoadConfig(provisioninginput.ConfigurationFileKey, r.clustersConfig)

	client, err := rancher.NewClient("", testSession)
	require.NoError(r.T(), err)

	r.client = client
}

func (r *HostnameTruncationTestSuite) TestProvisioningRKE2ClusterTruncation() {
	tests := []struct {
		name                        string
		machinePoolNameLengths      []int
		hostnameLengthLimits        []int
		defaultHostnameLengthLimits []int
	}{
		{
			name:                        "Cluster level truncation",
			machinePoolNameLengths:      []int{10, 31, 63},
			defaultHostnameLengthLimits: []int{10, 31, 63},
		},
		{
			name:                        "Machine pool level truncation - 10 characters",
			machinePoolNameLengths:      []int{10, 10, 10},
			hostnameLengthLimits:        []int{10, 31, 63},
			defaultHostnameLengthLimits: []int{10, 16, 63},
		},
		{
			name:                        "Machine pool level truncation - 31 characters",
			machinePoolNameLengths:      []int{10, 31, 63},
			hostnameLengthLimits:        []int{31, 31, 31},
			defaultHostnameLengthLimits: []int{10, 31, 63},
		},
		{
			name:                        "Machine pool level truncation - 63 characters",
			machinePoolNameLengths:      []int{10, 31, 63},
			hostnameLengthLimits:        []int{63, 63, 63},
			defaultHostnameLengthLimits: []int{10, 31, 63},
		},
		{
			name:                        "Cluster and machine pool level truncation - 31 characters",
			machinePoolNameLengths:      []int{10, 31, 63},
			hostnameLengthLimits:        []int{31, 31},
			defaultHostnameLengthLimits: []int{10, 63, 31},
		},
	}
	for _, tt := range tests {
		for _, defaultLength := range tt.defaultHostnameLengthLimits {
			r.Run(tt.name+fmt.Sprintf("_defaultHostnameLimit:%d", defaultLength), func() {
				var hostnamePools []machinepools.HostnameTruncation
				for i, nameLength := range tt.machinePoolNameLengths {
					currentTruncationPool := machinepools.HostnameTruncation{
						Name:                   namegen.RandStringLower(nameLength),
						ClusterNameLengthLimit: defaultLength,
					}
					if len(tt.hostnameLengthLimits) >= i && len(tt.hostnameLengthLimits) > 0 {
						currentTruncationPool.PoolNameLengthLimit = tt.hostnameLengthLimits[i]
					}
					hostnamePools = append(hostnamePools, currentTruncationPool)
				}
				testConfig := clusters.ConvertConfigToClusterConfig(r.clustersConfig)
				testConfig.KubernetesVersion = r.clustersConfig.RKE2KubernetesVersions[0]
				testConfig.CNI = r.clustersConfig.CNIs[0]

				rke2Provider, _, _, _ := permutations.GetClusterProvider(permutations.RKE2ProvisionCluster, (*testConfig.Providers)[0], r.clustersConfig)

				clusterObject, err := provisioning.CreateProvisioningCluster(r.client, *rke2Provider, testConfig, hostnamePools)
				require.NoError(r.T(), err)

				provisioning.VerifyCluster(r.T(), r.client, testConfig, clusterObject)
				provisioning.VerifyHostnameLength(r.T(), r.client, clusterObject)
			})
		}
	}
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestHostnameTruncationTestSuite(t *testing.T) {
	suite.Run(t, new(HostnameTruncationTestSuite))
}
