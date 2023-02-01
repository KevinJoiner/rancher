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

type AppRevisionHandler func(string, *v3.AppRevision) (*v3.AppRevision, error)

type AppRevisionController interface {
	generic.ControllerMeta
	AppRevisionClient

	OnChange(ctx context.Context, name string, sync AppRevisionHandler)
	OnRemove(ctx context.Context, name string, sync AppRevisionHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() AppRevisionCache
}

type AppRevisionClient interface {
	Create(*v3.AppRevision) (*v3.AppRevision, error)
	Update(*v3.AppRevision) (*v3.AppRevision, error)
	UpdateStatus(*v3.AppRevision) (*v3.AppRevision, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.AppRevision, error)
	List(namespace string, opts metav1.ListOptions) (*v3.AppRevisionList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.AppRevision, err error)
}

type AppRevisionCache interface {
	Get(namespace, name string) (*v3.AppRevision, error)
	List(namespace string, selector labels.Selector) ([]*v3.AppRevision, error)

	AddIndexer(indexName string, indexer AppRevisionIndexer)
	GetByIndex(indexName, key string) ([]*v3.AppRevision, error)
}

type AppRevisionIndexer func(obj *v3.AppRevision) ([]string, error)

type AppRevisionGenericController struct {
	generic.ControllerInterface[*v3.AppRevision, *v3.AppRevisionList]
}

func (c *AppRevisionGenericController) OnChange(ctx context.Context, name string, sync AppRevisionHandler) {
	c.ControllerInterface.OnChange(ctx, name, generic.ObjectHandler[*v3.AppRevision](sync))
}

func (c *AppRevisionGenericController) OnRemove(ctx context.Context, name string, sync AppRevisionHandler) {
	c.ControllerInterface.OnRemove(ctx, name, generic.ObjectHandler[*v3.AppRevision](sync))
}

func (c *AppRevisionGenericController) Cache() AppRevisionCache {
	return &AppRevisionGenericCache{
		c.ControllerInterface.Cache(),
	}
}

type AppRevisionGenericCache struct {
	generic.CacheInterface[*v3.AppRevision]
}

func (c AppRevisionGenericCache) AddIndexer(indexName string, indexer AppRevisionIndexer) {
	c.CacheInterface.AddIndexer(indexName, generic.Indexer[*v3.AppRevision](indexer))
}

type AppRevisionStatusHandler func(obj *v3.AppRevision, status v3.AppRevisionStatus) (v3.AppRevisionStatus, error)

type AppRevisionGeneratingHandler func(obj *v3.AppRevision, status v3.AppRevisionStatus) ([]runtime.Object, v3.AppRevisionStatus, error)

func FromAppRevisionHandlerToHandler(sync AppRevisionHandler) generic.Handler {
	return generic.FromObjectHandlerToHandler(generic.ObjectHandler[*v3.AppRevision](sync))
}

func RegisterAppRevisionStatusHandler(ctx context.Context, controller AppRevisionController, condition condition.Cond, name string, handler AppRevisionStatusHandler) {
	statusHandler := &appRevisionStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromAppRevisionHandlerToHandler(statusHandler.sync))
}

func RegisterAppRevisionGeneratingHandler(ctx context.Context, controller AppRevisionController, apply apply.Apply,
	condition condition.Cond, name string, handler AppRevisionGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &appRevisionGeneratingHandler{
		AppRevisionGeneratingHandler: handler,
		apply:                        apply,
		name:                         name,
		gvk:                          controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterAppRevisionStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type appRevisionStatusHandler struct {
	client    AppRevisionClient
	condition condition.Cond
	handler   AppRevisionStatusHandler
}

func (a *appRevisionStatusHandler) sync(key string, obj *v3.AppRevision) (*v3.AppRevision, error) {
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

type appRevisionGeneratingHandler struct {
	AppRevisionGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *appRevisionGeneratingHandler) Remove(key string, obj *v3.AppRevision) (*v3.AppRevision, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v3.AppRevision{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *appRevisionGeneratingHandler) Handle(obj *v3.AppRevision, status v3.AppRevisionStatus) (v3.AppRevisionStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.AppRevisionGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
