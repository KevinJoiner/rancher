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

type CatalogTemplateVersionHandler func(string, *v3.CatalogTemplateVersion) (*v3.CatalogTemplateVersion, error)

type CatalogTemplateVersionController interface {
	generic.ControllerMeta
	CatalogTemplateVersionClient

	OnChange(ctx context.Context, name string, sync CatalogTemplateVersionHandler)
	OnRemove(ctx context.Context, name string, sync CatalogTemplateVersionHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() CatalogTemplateVersionCache
}

type CatalogTemplateVersionClient interface {
	Create(*v3.CatalogTemplateVersion) (*v3.CatalogTemplateVersion, error)
	Update(*v3.CatalogTemplateVersion) (*v3.CatalogTemplateVersion, error)
	UpdateStatus(*v3.CatalogTemplateVersion) (*v3.CatalogTemplateVersion, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.CatalogTemplateVersion, error)
	List(namespace string, opts metav1.ListOptions) (*v3.CatalogTemplateVersionList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.CatalogTemplateVersion, err error)
}

type CatalogTemplateVersionCache interface {
	Get(namespace, name string) (*v3.CatalogTemplateVersion, error)
	List(namespace string, selector labels.Selector) ([]*v3.CatalogTemplateVersion, error)

	AddIndexer(indexName string, indexer CatalogTemplateVersionIndexer)
	GetByIndex(indexName, key string) ([]*v3.CatalogTemplateVersion, error)
}

type CatalogTemplateVersionIndexer func(obj *v3.CatalogTemplateVersion) ([]string, error)

type CatalogTemplateVersionGenericController struct {
	generic.ControllerInterface[*v3.CatalogTemplateVersion, *v3.CatalogTemplateVersionList]
}

func (c *CatalogTemplateVersionGenericController) OnChange(ctx context.Context, name string, sync CatalogTemplateVersionHandler) {
	c.ControllerInterface.OnChange(ctx, name, generic.ObjectHandler[*v3.CatalogTemplateVersion](sync))
}

func (c *CatalogTemplateVersionGenericController) OnRemove(ctx context.Context, name string, sync CatalogTemplateVersionHandler) {
	c.ControllerInterface.OnRemove(ctx, name, generic.ObjectHandler[*v3.CatalogTemplateVersion](sync))
}

func (c *CatalogTemplateVersionGenericController) Cache() CatalogTemplateVersionCache {
	return &CatalogTemplateVersionGenericCache{
		c.ControllerInterface.Cache(),
	}
}

type CatalogTemplateVersionGenericCache struct {
	generic.CacheInterface[*v3.CatalogTemplateVersion]
}

func (c CatalogTemplateVersionGenericCache) AddIndexer(indexName string, indexer CatalogTemplateVersionIndexer) {
	c.CacheInterface.AddIndexer(indexName, generic.Indexer[*v3.CatalogTemplateVersion](indexer))
}

type CatalogTemplateVersionStatusHandler func(obj *v3.CatalogTemplateVersion, status v3.TemplateVersionStatus) (v3.TemplateVersionStatus, error)

type CatalogTemplateVersionGeneratingHandler func(obj *v3.CatalogTemplateVersion, status v3.TemplateVersionStatus) ([]runtime.Object, v3.TemplateVersionStatus, error)

func FromCatalogTemplateVersionHandlerToHandler(sync CatalogTemplateVersionHandler) generic.Handler {
	return generic.FromObjectHandlerToHandler(generic.ObjectHandler[*v3.CatalogTemplateVersion](sync))
}

func RegisterCatalogTemplateVersionStatusHandler(ctx context.Context, controller CatalogTemplateVersionController, condition condition.Cond, name string, handler CatalogTemplateVersionStatusHandler) {
	statusHandler := &catalogTemplateVersionStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromCatalogTemplateVersionHandlerToHandler(statusHandler.sync))
}

func RegisterCatalogTemplateVersionGeneratingHandler(ctx context.Context, controller CatalogTemplateVersionController, apply apply.Apply,
	condition condition.Cond, name string, handler CatalogTemplateVersionGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &catalogTemplateVersionGeneratingHandler{
		CatalogTemplateVersionGeneratingHandler: handler,
		apply:                                   apply,
		name:                                    name,
		gvk:                                     controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterCatalogTemplateVersionStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type catalogTemplateVersionStatusHandler struct {
	client    CatalogTemplateVersionClient
	condition condition.Cond
	handler   CatalogTemplateVersionStatusHandler
}

func (a *catalogTemplateVersionStatusHandler) sync(key string, obj *v3.CatalogTemplateVersion) (*v3.CatalogTemplateVersion, error) {
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

type catalogTemplateVersionGeneratingHandler struct {
	CatalogTemplateVersionGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *catalogTemplateVersionGeneratingHandler) Remove(key string, obj *v3.CatalogTemplateVersion) (*v3.CatalogTemplateVersion, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v3.CatalogTemplateVersion{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *catalogTemplateVersionGeneratingHandler) Handle(obj *v3.CatalogTemplateVersion, status v3.TemplateVersionStatus) (v3.TemplateVersionStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.CatalogTemplateVersionGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
