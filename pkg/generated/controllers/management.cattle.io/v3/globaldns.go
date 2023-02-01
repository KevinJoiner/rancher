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
	"github.com/rancher/wrangler/pkg/apply"
	"github.com/rancher/wrangler/pkg/condition"
	"github.com/rancher/wrangler/pkg/generic"
	"github.com/rancher/wrangler/pkg/kv"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

type GlobalDnsHandler func(string, *v3.GlobalDns) (*v3.GlobalDns, error)

type GlobalDnsController interface {
	generic.ControllerMeta
	GlobalDnsClient

	OnChange(ctx context.Context, name string, sync GlobalDnsHandler)
	OnRemove(ctx context.Context, name string, sync GlobalDnsHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() GlobalDnsCache
}

type GlobalDnsClient interface {
	Create(*v3.GlobalDns) (*v3.GlobalDns, error)
	Update(*v3.GlobalDns) (*v3.GlobalDns, error)
	UpdateStatus(*v3.GlobalDns) (*v3.GlobalDns, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.GlobalDns, error)
	List(namespace string, opts metav1.ListOptions) (*v3.GlobalDnsList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.GlobalDns, err error)
}

type GlobalDnsCache interface {
	Get(namespace, name string) (*v3.GlobalDns, error)
	List(namespace string, selector labels.Selector) ([]*v3.GlobalDns, error)

	AddIndexer(indexName string, indexer GlobalDnsIndexer)
	GetByIndex(indexName, key string) ([]*v3.GlobalDns, error)
}

type GlobalDnsIndexer func(obj *v3.GlobalDns) ([]string, error)

type GlobalDnsGenericController struct {
	generic.ControllerInterface[*v3.GlobalDns, *v3.GlobalDnsList]
}

func (c *GlobalDnsGenericController) OnChange(ctx context.Context, name string, sync GlobalDnsHandler) {
	c.ControllerInterface.OnChange(ctx, name, generic.ObjectHandler[*v3.GlobalDns](sync))
}

func (c *GlobalDnsGenericController) OnRemove(ctx context.Context, name string, sync GlobalDnsHandler) {
	c.ControllerInterface.OnRemove(ctx, name, generic.ObjectHandler[*v3.GlobalDns](sync))
}

func (c *GlobalDnsGenericController) Cache() GlobalDnsCache {
	return &GlobalDnsGenericCache{
		c.ControllerInterface.Cache(),
	}
}

type GlobalDnsGenericCache struct {
	generic.CacheInterface[*v3.GlobalDns]
}

func (c GlobalDnsGenericCache) AddIndexer(indexName string, indexer GlobalDnsIndexer) {
	c.CacheInterface.AddIndexer(indexName, generic.Indexer[*v3.GlobalDns](indexer))
}

type GlobalDnsStatusHandler func(obj *v3.GlobalDns, status v3.GlobalDNSStatus) (v3.GlobalDNSStatus, error)

type GlobalDnsGeneratingHandler func(obj *v3.GlobalDns, status v3.GlobalDNSStatus) ([]runtime.Object, v3.GlobalDNSStatus, error)

func FromGlobalDnsHandlerToHandler(sync GlobalDnsHandler) generic.Handler {
	return generic.FromObjectHandlerToHandler(generic.ObjectHandler[*v3.GlobalDns](sync))
}

func RegisterGlobalDnsStatusHandler(ctx context.Context, controller GlobalDnsController, condition condition.Cond, name string, handler GlobalDnsStatusHandler) {
	statusHandler := &globalDnsStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromGlobalDnsHandlerToHandler(statusHandler.sync))
}

func RegisterGlobalDnsGeneratingHandler(ctx context.Context, controller GlobalDnsController, apply apply.Apply,
	condition condition.Cond, name string, handler GlobalDnsGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &globalDnsGeneratingHandler{
		GlobalDnsGeneratingHandler: handler,
		apply:                      apply,
		name:                       name,
		gvk:                        controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterGlobalDnsStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type globalDnsStatusHandler struct {
	client    GlobalDnsClient
	condition condition.Cond
	handler   GlobalDnsStatusHandler
}

func (a *globalDnsStatusHandler) sync(key string, obj *v3.GlobalDns) (*v3.GlobalDns, error) {
	if obj == nil {
		return obj, nil
	}

	origStatus := obj.Status.DeepCopy()
	obj = obj.DeepCopy()
	newStatus, err := a.handler(obj, obj.Status)
	if err != nil {
		// Revert to old status on error
		newStatus = *origStatus.DeepCopy()
	}

	if a.condition != "" {
		if errors.IsConflict(err) {
			a.condition.SetError(&newStatus, "", nil)
		} else {
			a.condition.SetError(&newStatus, "", err)
		}
	}
	if !equality.Semantic.DeepEqual(origStatus, &newStatus) {
		if a.condition != "" {
			// Since status has changed, update the lastUpdatedTime
			a.condition.LastUpdated(&newStatus, time.Now().UTC().Format(time.RFC3339))
		}

		var newErr error
		obj.Status = newStatus
		newObj, newErr := a.client.UpdateStatus(obj)
		if err == nil {
			err = newErr
		}
		if newErr == nil {
			obj = newObj
		}
	}
	return obj, err
}

type globalDnsGeneratingHandler struct {
	GlobalDnsGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *globalDnsGeneratingHandler) Remove(key string, obj *v3.GlobalDns) (*v3.GlobalDns, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v3.GlobalDns{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *globalDnsGeneratingHandler) Handle(obj *v3.GlobalDns, status v3.GlobalDNSStatus) (v3.GlobalDNSStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.GlobalDnsGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
