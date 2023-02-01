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

type TokenHandler func(string, *v3.Token) (*v3.Token, error)

type TokenController interface {
	generic.ControllerMeta
	TokenClient

	OnChange(ctx context.Context, name string, sync TokenHandler)
	OnRemove(ctx context.Context, name string, sync TokenHandler)
	Enqueue(name string)
	EnqueueAfter(name string, duration time.Duration)

	Cache() TokenCache
}

type TokenClient interface {
	Create(*v3.Token) (*v3.Token, error)
	Update(*v3.Token) (*v3.Token, error)

	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*v3.Token, error)
	List(opts metav1.ListOptions) (*v3.TokenList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.Token, err error)
}

type TokenCache interface {
	Get(name string) (*v3.Token, error)
	List(selector labels.Selector) ([]*v3.Token, error)

	AddIndexer(indexName string, indexer TokenIndexer)
	GetByIndex(indexName, key string) ([]*v3.Token, error)
}

type TokenIndexer func(obj *v3.Token) ([]string, error)

type TokenGenericController struct {
	generic.NonNamespacedControllerInterface[*v3.Token, *v3.TokenList]
}

func (c *TokenGenericController) OnChange(ctx context.Context, name string, sync TokenHandler) {
	c.NonNamespacedControllerInterface.OnChange(ctx, name, generic.ObjectHandler[*v3.Token](sync))
}

func (c *TokenGenericController) OnRemove(ctx context.Context, name string, sync TokenHandler) {
	c.NonNamespacedControllerInterface.OnRemove(ctx, name, generic.ObjectHandler[*v3.Token](sync))
}

func (c *TokenGenericController) Cache() TokenCache {
	return &TokenGenericCache{
		c.NonNamespacedControllerInterface.Cache(),
	}
}

type TokenGenericCache struct {
	generic.NonNamespacedCacheInterface[*v3.Token]
}

func (c TokenGenericCache) AddIndexer(indexName string, indexer TokenIndexer) {
	c.NonNamespacedCacheInterface.AddIndexer(indexName, generic.Indexer[*v3.Token](indexer))
}
