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

type ClusterLoggingHandler func(string, *v3.ClusterLogging) (*v3.ClusterLogging, error)

type ClusterLoggingController interface {
	generic.ControllerMeta
	ClusterLoggingClient

	OnChange(ctx context.Context, name string, sync ClusterLoggingHandler)
	OnRemove(ctx context.Context, name string, sync ClusterLoggingHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() ClusterLoggingCache
}

type ClusterLoggingClient interface {
	Create(*v3.ClusterLogging) (*v3.ClusterLogging, error)
	Update(*v3.ClusterLogging) (*v3.ClusterLogging, error)
	UpdateStatus(*v3.ClusterLogging) (*v3.ClusterLogging, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v3.ClusterLogging, error)
	List(namespace string, opts metav1.ListOptions) (*v3.ClusterLoggingList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.ClusterLogging, err error)
}

type ClusterLoggingCache interface {
	Get(namespace, name string) (*v3.ClusterLogging, error)
	List(namespace string, selector labels.Selector) ([]*v3.ClusterLogging, error)

	AddIndexer(indexName string, indexer ClusterLoggingIndexer)
	GetByIndex(indexName, key string) ([]*v3.ClusterLogging, error)
}

type ClusterLoggingIndexer func(obj *v3.ClusterLogging) ([]string, error)

type ClusterLoggingGenericController struct {
	generic.ControllerInterface[*v3.ClusterLogging, *v3.ClusterLoggingList]
}

func (c *ClusterLoggingGenericController) OnChange(ctx context.Context, name string, sync ClusterLoggingHandler) {
	c.ControllerInterface.OnChange(ctx, name, generic.ObjectHandler[*v3.ClusterLogging](sync))
}

func (c *ClusterLoggingGenericController) OnRemove(ctx context.Context, name string, sync ClusterLoggingHandler) {
	c.ControllerInterface.OnRemove(ctx, name, generic.ObjectHandler[*v3.ClusterLogging](sync))
}

func (c *ClusterLoggingGenericController) Cache() ClusterLoggingCache {
	return &ClusterLoggingGenericCache{
		c.ControllerInterface.Cache(),
	}
}

type ClusterLoggingGenericCache struct {
	generic.CacheInterface[*v3.ClusterLogging]
}

func (c ClusterLoggingGenericCache) AddIndexer(indexName string, indexer ClusterLoggingIndexer) {
	c.CacheInterface.AddIndexer(indexName, generic.Indexer[*v3.ClusterLogging](indexer))
}

type ClusterLoggingStatusHandler func(obj *v3.ClusterLogging, status v3.ClusterLoggingStatus) (v3.ClusterLoggingStatus, error)

type ClusterLoggingGeneratingHandler func(obj *v3.ClusterLogging, status v3.ClusterLoggingStatus) ([]runtime.Object, v3.ClusterLoggingStatus, error)

func FromClusterLoggingHandlerToHandler(sync ClusterLoggingHandler) generic.Handler {
	return generic.FromObjectHandlerToHandler(generic.ObjectHandler[*v3.ClusterLogging](sync))
}

func RegisterClusterLoggingStatusHandler(ctx context.Context, controller ClusterLoggingController, condition condition.Cond, name string, handler ClusterLoggingStatusHandler) {
	statusHandler := &clusterLoggingStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromClusterLoggingHandlerToHandler(statusHandler.sync))
}

func RegisterClusterLoggingGeneratingHandler(ctx context.Context, controller ClusterLoggingController, apply apply.Apply,
	condition condition.Cond, name string, handler ClusterLoggingGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &clusterLoggingGeneratingHandler{
		ClusterLoggingGeneratingHandler: handler,
		apply:                           apply,
		name:                            name,
		gvk:                             controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterClusterLoggingStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type clusterLoggingStatusHandler struct {
	client    ClusterLoggingClient
	condition condition.Cond
	handler   ClusterLoggingStatusHandler
}

func (a *clusterLoggingStatusHandler) sync(key string, obj *v3.ClusterLogging) (*v3.ClusterLogging, error) {
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

type clusterLoggingGeneratingHandler struct {
	ClusterLoggingGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *clusterLoggingGeneratingHandler) Remove(key string, obj *v3.ClusterLogging) (*v3.ClusterLogging, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v3.ClusterLogging{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *clusterLoggingGeneratingHandler) Handle(obj *v3.ClusterLogging, status v3.ClusterLoggingStatus) (v3.ClusterLoggingStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.ClusterLoggingGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
