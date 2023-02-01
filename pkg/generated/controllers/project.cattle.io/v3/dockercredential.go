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

type DockerCredentialHandler func(string, *v3.DockerCredential) (*v3.DockerCredential, error)

type DockerCredentialController interface {
	generic.ControllerMeta
	DockerCredentialClient

	OnChange(ctx context.Context, name string, sync DockerCredentialHandler)
	OnRemove(ctx context.Context, name string, sync DockerCredentialHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() DockerCredentialCache
}

type DockerCredentialClient interface {
	Create(*v3.DockerCredential) (*v3.DockerCredential, error)
	Update(*v3.DockerCredential) (*v3.DockerCredential, error)

	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.DockerCredential, error)
	List(namespace string, opts metav1.ListOptions) (*v3.DockerCredentialList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.DockerCredential, err error)
}

type DockerCredentialCache interface {
	Get(namespace, name string) (*v3.DockerCredential, error)
	List(namespace string, selector labels.Selector) ([]*v3.DockerCredential, error)

	AddIndexer(indexName string, indexer DockerCredentialIndexer)
	GetByIndex(indexName, key string) ([]*v3.DockerCredential, error)
}

type DockerCredentialIndexer func(obj *v3.DockerCredential) ([]string, error)

type DockerCredentialGenericController struct {
	generic.ControllerInterface[*v3.DockerCredential, *v3.DockerCredentialList]
}

func (c *DockerCredentialGenericController) OnChange(ctx context.Context, name string, sync DockerCredentialHandler) {
	c.ControllerInterface.OnChange(ctx, name, generic.ObjectHandler[*v3.DockerCredential](sync))
}

func (c *DockerCredentialGenericController) OnRemove(ctx context.Context, name string, sync DockerCredentialHandler) {
	c.ControllerInterface.OnRemove(ctx, name, generic.ObjectHandler[*v3.DockerCredential](sync))
}

func (c *DockerCredentialGenericController) Cache() DockerCredentialCache {
	return &DockerCredentialGenericCache{
		c.ControllerInterface.Cache(),
	}
}

type DockerCredentialGenericCache struct {
	generic.CacheInterface[*v3.DockerCredential]
}

func (c DockerCredentialGenericCache) AddIndexer(indexName string, indexer DockerCredentialIndexer) {
	c.CacheInterface.AddIndexer(indexName, generic.Indexer[*v3.DockerCredential](indexer))
}
