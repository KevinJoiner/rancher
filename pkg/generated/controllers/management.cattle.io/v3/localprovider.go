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

type LocalProviderHandler func(string, *v3.LocalProvider) (*v3.LocalProvider, error)

type LocalProviderController interface {
	generic.ControllerMeta
	LocalProviderClient

	OnChange(ctx context.Context, name string, sync LocalProviderHandler)
	OnRemove(ctx context.Context, name string, sync LocalProviderHandler)
	Enqueue(name string)
	EnqueueAfter(name string, duration time.Duration)

	Cache() LocalProviderCache
}

type LocalProviderClient interface {
	Create(*v3.LocalProvider) (*v3.LocalProvider, error)
	Update(*v3.LocalProvider) (*v3.LocalProvider, error)

	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*v3.LocalProvider, error)
	List(opts metav1.ListOptions) (*v3.LocalProviderList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.LocalProvider, err error)
}

type LocalProviderCache interface {
	Get(name string) (*v3.LocalProvider, error)
	List(selector labels.Selector) ([]*v3.LocalProvider, error)

	AddIndexer(indexName string, indexer LocalProviderIndexer)
	GetByIndex(indexName, key string) ([]*v3.LocalProvider, error)
}

type LocalProviderIndexer func(obj *v3.LocalProvider) ([]string, error)

type LocalProviderGenericController struct {
	generic.NonNamespacedControllerInterface[*v3.LocalProvider, *v3.LocalProviderList]
}

func (c *LocalProviderGenericController) OnChange(ctx context.Context, name string, sync LocalProviderHandler) {
	c.NonNamespacedControllerInterface.OnChange(ctx, name, generic.ObjectHandler[*v3.LocalProvider](sync))
}

func (c *LocalProviderGenericController) OnRemove(ctx context.Context, name string, sync LocalProviderHandler) {
	c.NonNamespacedControllerInterface.OnRemove(ctx, name, generic.ObjectHandler[*v3.LocalProvider](sync))
}

func (c *LocalProviderGenericController) Cache() LocalProviderCache {
	return &LocalProviderGenericCache{
		c.NonNamespacedControllerInterface.Cache(),
	}
}

type LocalProviderGenericCache struct {
	generic.NonNamespacedCacheInterface[*v3.LocalProvider]
}

func (c LocalProviderGenericCache) AddIndexer(indexName string, indexer LocalProviderIndexer) {
	c.NonNamespacedCacheInterface.AddIndexer(indexName, generic.Indexer[*v3.LocalProvider](indexer))
}
