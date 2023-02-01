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

type ComposeConfigHandler func(string, *v3.ComposeConfig) (*v3.ComposeConfig, error)

type ComposeConfigController interface {
	generic.ControllerMeta
	ComposeConfigClient

	OnChange(ctx context.Context, name string, sync ComposeConfigHandler)
	OnRemove(ctx context.Context, name string, sync ComposeConfigHandler)
	Enqueue(name string)
	EnqueueAfter(name string, duration time.Duration)

	Cache() ComposeConfigCache
}

type ComposeConfigClient interface {
	Create(*v3.ComposeConfig) (*v3.ComposeConfig, error)
	Update(*v3.ComposeConfig) (*v3.ComposeConfig, error)
	UpdateStatus(*v3.ComposeConfig) (*v3.ComposeConfig, error)
	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*v3.ComposeConfig, error)
	List(opts metav1.ListOptions) (*v3.ComposeConfigList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.ComposeConfig, err error)
}

type ComposeConfigCache interface {
	Get(name string) (*v3.ComposeConfig, error)
	List(selector labels.Selector) ([]*v3.ComposeConfig, error)

	AddIndexer(indexName string, indexer ComposeConfigIndexer)
	GetByIndex(indexName, key string) ([]*v3.ComposeConfig, error)
}

type ComposeConfigIndexer func(obj *v3.ComposeConfig) ([]string, error)

type ComposeConfigGenericController struct {
	generic.NonNamespacedControllerInterface[*v3.ComposeConfig, *v3.ComposeConfigList]
}

func (c *ComposeConfigGenericController) OnChange(ctx context.Context, name string, sync ComposeConfigHandler) {
	c.NonNamespacedControllerInterface.OnChange(ctx, name, generic.ObjectHandler[*v3.ComposeConfig](sync))
}

func (c *ComposeConfigGenericController) OnRemove(ctx context.Context, name string, sync ComposeConfigHandler) {
	c.NonNamespacedControllerInterface.OnRemove(ctx, name, generic.ObjectHandler[*v3.ComposeConfig](sync))
}

func (c *ComposeConfigGenericController) Cache() ComposeConfigCache {
	return &ComposeConfigGenericCache{
		c.NonNamespacedControllerInterface.Cache(),
	}
}

type ComposeConfigGenericCache struct {
	generic.NonNamespacedCacheInterface[*v3.ComposeConfig]
}

func (c ComposeConfigGenericCache) AddIndexer(indexName string, indexer ComposeConfigIndexer) {
	c.NonNamespacedCacheInterface.AddIndexer(indexName, generic.Indexer[*v3.ComposeConfig](indexer))
}

type ComposeConfigStatusHandler func(obj *v3.ComposeConfig, status v3.ComposeStatus) (v3.ComposeStatus, error)

type ComposeConfigGeneratingHandler func(obj *v3.ComposeConfig, status v3.ComposeStatus) ([]runtime.Object, v3.ComposeStatus, error)

func FromComposeConfigHandlerToHandler(sync ComposeConfigHandler) generic.Handler {
	return generic.FromObjectHandlerToHandler(generic.ObjectHandler[*v3.ComposeConfig](sync))
}

func RegisterComposeConfigStatusHandler(ctx context.Context, controller ComposeConfigController, condition condition.Cond, name string, handler ComposeConfigStatusHandler) {
	statusHandler := &composeConfigStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromComposeConfigHandlerToHandler(statusHandler.sync))
}

func RegisterComposeConfigGeneratingHandler(ctx context.Context, controller ComposeConfigController, apply apply.Apply,
	condition condition.Cond, name string, handler ComposeConfigGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &composeConfigGeneratingHandler{
		ComposeConfigGeneratingHandler: handler,
		apply:                          apply,
		name:                           name,
		gvk:                            controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterComposeConfigStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type composeConfigStatusHandler struct {
	client    ComposeConfigClient
	condition condition.Cond
	handler   ComposeConfigStatusHandler
}

func (a *composeConfigStatusHandler) sync(key string, obj *v3.ComposeConfig) (*v3.ComposeConfig, error) {
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

type composeConfigGeneratingHandler struct {
	ComposeConfigGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *composeConfigGeneratingHandler) Remove(key string, obj *v3.ComposeConfig) (*v3.ComposeConfig, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v3.ComposeConfig{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *composeConfigGeneratingHandler) Handle(obj *v3.ComposeConfig, status v3.ComposeStatus) (v3.ComposeStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.ComposeConfigGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
