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
	"context"
	"time"

	v3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/rancher/wrangler/pkg/generic"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

type PrincipalHandler func(string, *v3.Principal) (*v3.Principal, error)

type PrincipalController interface {
	generic.ControllerMeta
	PrincipalClient

	OnChange(ctx context.Context, name string, sync PrincipalHandler)
	OnRemove(ctx context.Context, name string, sync PrincipalHandler)
	Enqueue(name string)
	EnqueueAfter(name string, duration time.Duration)

	Cache() PrincipalCache
}

type PrincipalClient interface {
	Create(*v3.Principal) (*v3.Principal, error)
	Update(*v3.Principal) (*v3.Principal, error)

	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*v3.Principal, error)
	List(opts metav1.ListOptions) (*v3.PrincipalList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.Principal, err error)
}

type PrincipalCache interface {
	Get(name string) (*v3.Principal, error)
	List(selector labels.Selector) ([]*v3.Principal, error)

	AddIndexer(indexName string, indexer PrincipalIndexer)
	GetByIndex(indexName, key string) ([]*v3.Principal, error)
}

type PrincipalIndexer func(obj *v3.Principal) ([]string, error)

type PrincipalGenericController struct {
	generic.NonNamespacedControllerInterface[*v3.Principal, *v3.PrincipalList]
}

func (c *PrincipalGenericController) OnChange(ctx context.Context, name string, sync PrincipalHandler) {
	c.NonNamespacedControllerInterface.OnChange(ctx, name, generic.ObjectHandler[*v3.Principal](sync))
}

func (c *PrincipalGenericController) OnRemove(ctx context.Context, name string, sync PrincipalHandler) {
	c.NonNamespacedControllerInterface.OnRemove(ctx, name, generic.ObjectHandler[*v3.Principal](sync))
}

func (c *PrincipalGenericController) Cache() PrincipalCache {
	return &PrincipalGenericCache{
		c.NonNamespacedControllerInterface.Cache(),
	}
}

type PrincipalGenericCache struct {
	generic.NonNamespacedCacheInterface[*v3.Principal]
}

func (c PrincipalGenericCache) AddIndexer(indexName string, indexer PrincipalIndexer) {
	c.NonNamespacedCacheInterface.AddIndexer(indexName, generic.Indexer[*v3.Principal](indexer))
}
