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

type ActiveDirectoryProviderHandler func(string, *v3.ActiveDirectoryProvider) (*v3.ActiveDirectoryProvider, error)

type ActiveDirectoryProviderController interface {
	generic.ControllerMeta
	ActiveDirectoryProviderClient

	OnChange(ctx context.Context, name string, sync ActiveDirectoryProviderHandler)
	OnRemove(ctx context.Context, name string, sync ActiveDirectoryProviderHandler)
	Enqueue(name string)
	EnqueueAfter(name string, duration time.Duration)

	Cache() ActiveDirectoryProviderCache
}

type ActiveDirectoryProviderClient interface {
	Create(*v3.ActiveDirectoryProvider) (*v3.ActiveDirectoryProvider, error)
	Update(*v3.ActiveDirectoryProvider) (*v3.ActiveDirectoryProvider, error)

	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*v3.ActiveDirectoryProvider, error)
	List(opts metav1.ListOptions) (*v3.ActiveDirectoryProviderList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.ActiveDirectoryProvider, err error)
}

type ActiveDirectoryProviderCache interface {
	Get(name string) (*v3.ActiveDirectoryProvider, error)
	List(selector labels.Selector) ([]*v3.ActiveDirectoryProvider, error)

	AddIndexer(indexName string, indexer ActiveDirectoryProviderIndexer)
	GetByIndex(indexName, key string) ([]*v3.ActiveDirectoryProvider, error)
}

type ActiveDirectoryProviderIndexer func(obj *v3.ActiveDirectoryProvider) ([]string, error)

type ActiveDirectoryProviderGenericController struct {
	generic.NonNamespacedControllerInterface[*v3.ActiveDirectoryProvider, *v3.ActiveDirectoryProviderList]
}

func (c *ActiveDirectoryProviderGenericController) OnChange(ctx context.Context, name string, sync ActiveDirectoryProviderHandler) {
	c.NonNamespacedControllerInterface.OnChange(ctx, name, generic.ObjectHandler[*v3.ActiveDirectoryProvider](sync))
}

func (c *ActiveDirectoryProviderGenericController) OnRemove(ctx context.Context, name string, sync ActiveDirectoryProviderHandler) {
	c.NonNamespacedControllerInterface.OnRemove(ctx, name, generic.ObjectHandler[*v3.ActiveDirectoryProvider](sync))
}

func (c *ActiveDirectoryProviderGenericController) Cache() ActiveDirectoryProviderCache {
	return &ActiveDirectoryProviderGenericCache{
		c.NonNamespacedControllerInterface.Cache(),
	}
}

type ActiveDirectoryProviderGenericCache struct {
	generic.NonNamespacedCacheInterface[*v3.ActiveDirectoryProvider]
}

func (c ActiveDirectoryProviderGenericCache) AddIndexer(indexName string, indexer ActiveDirectoryProviderIndexer) {
	c.NonNamespacedCacheInterface.AddIndexer(indexName, generic.Indexer[*v3.ActiveDirectoryProvider](indexer))
}
