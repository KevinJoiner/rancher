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

	v3 "github.com/rancher/rancher/pkg/apis/project.cattle.io/v3"
	"github.com/rancher/wrangler/pkg/generic"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

type NamespacedServiceAccountTokenHandler func(string, *v3.NamespacedServiceAccountToken) (*v3.NamespacedServiceAccountToken, error)

type NamespacedServiceAccountTokenController interface {
	generic.ControllerMeta
	NamespacedServiceAccountTokenClient

	OnChange(ctx context.Context, name string, sync NamespacedServiceAccountTokenHandler)
	OnRemove(ctx context.Context, name string, sync NamespacedServiceAccountTokenHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() NamespacedServiceAccountTokenCache
}

type NamespacedServiceAccountTokenClient interface {
	Create(*v3.NamespacedServiceAccountToken) (*v3.NamespacedServiceAccountToken, error)
	Update(*v3.NamespacedServiceAccountToken) (*v3.NamespacedServiceAccountToken, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.NamespacedServiceAccountToken, error)
	List(namespace string, opts metav1.ListOptions) (*v3.NamespacedServiceAccountTokenList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.NamespacedServiceAccountToken, err error)
}

type NamespacedServiceAccountTokenCache interface {
	Get(namespace, name string) (*v3.NamespacedServiceAccountToken, error)
	List(namespace string, selector labels.Selector) ([]*v3.NamespacedServiceAccountToken, error)

	AddIndexer(indexName string, indexer NamespacedServiceAccountTokenIndexer)
	GetByIndex(indexName, key string) ([]*v3.NamespacedServiceAccountToken, error)
}

type NamespacedServiceAccountTokenIndexer func(obj *v3.NamespacedServiceAccountToken) ([]string, error)

type NamespacedServiceAccountTokenGenericController struct {
	generic.ControllerInterface[*v3.NamespacedServiceAccountToken, *v3.NamespacedServiceAccountTokenList]
}

func (c *NamespacedServiceAccountTokenGenericController) OnChange(ctx context.Context, name string, sync NamespacedServiceAccountTokenHandler) {
	c.ControllerInterface.OnChange(ctx, name, generic.ObjectHandler[*v3.NamespacedServiceAccountToken](sync))
}

func (c *NamespacedServiceAccountTokenGenericController) OnRemove(ctx context.Context, name string, sync NamespacedServiceAccountTokenHandler) {
	c.ControllerInterface.OnRemove(ctx, name, generic.ObjectHandler[*v3.NamespacedServiceAccountToken](sync))
}

func (c *NamespacedServiceAccountTokenGenericController) Cache() NamespacedServiceAccountTokenCache {
	return &NamespacedServiceAccountTokenGenericCache{
		c.ControllerInterface.Cache(),
	}
}

type NamespacedServiceAccountTokenGenericCache struct {
	generic.CacheInterface[*v3.NamespacedServiceAccountToken]
}

func (c NamespacedServiceAccountTokenGenericCache) AddIndexer(indexName string, indexer NamespacedServiceAccountTokenIndexer) {
	c.CacheInterface.AddIndexer(indexName, generic.Indexer[*v3.NamespacedServiceAccountToken](indexer))
}
