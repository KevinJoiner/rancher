/*
Copyright 2023 Rancher Labs, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by main. DO NOT EDIT.

package v3

import (
	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/rancher/wrangler/pkg/generic"
)

// GoogleOAuthProviderController interface for managing GoogleOAuthProvider resources.
type GoogleOAuthProviderController interface {
	generic.NoNsControllerInterface[*v3.GoogleOAuthProvider, *v3.GoogleOAuthProviderList]
}

// GoogleOAuthProviderClient interface for managing GoogleOAuthProvider resources in Kubernetes.
type GoogleOAuthProviderClient interface {
	generic.NoNsClientInterface[*v3.GoogleOAuthProvider, *v3.GoogleOAuthProviderList]
}

// GoogleOAuthProviderCache interface for retrieving GoogleOAuthProvider resources in memory.
type GoogleOAuthProviderCache interface {
	generic.NoNsCacheInterface[*v3.GoogleOAuthProvider]
}
