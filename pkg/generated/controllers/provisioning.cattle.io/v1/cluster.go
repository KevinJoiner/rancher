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

package v1

import (
	"context"
	"time"

	v1 "github.com/rancher/rancher/pkg/apis/provisioning.cattle.io/v1"
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

type ClusterHandler func(string, *v1.Cluster) (*v1.Cluster, error)

type ClusterController interface {
	generic.ControllerMeta
	ClusterClient

	OnChange(ctx context.Context, name string, sync ClusterHandler)
	OnRemove(ctx context.Context, name string, sync ClusterHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() ClusterCache
}

type ClusterClient interface {
	Create(*v1.Cluster) (*v1.Cluster, error)
	Update(*v1.Cluster) (*v1.Cluster, error)
	UpdateStatus(*v1.Cluster) (*v1.Cluster, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1.Cluster, error)
	List(namespace string, opts metav1.ListOptions) (*v1.ClusterList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Cluster, err error)
}

type ClusterCache interface {
	Get(namespace, name string) (*v1.Cluster, error)
	List(namespace string, selector labels.Selector) ([]*v1.Cluster, error)

	AddIndexer(indexName string, indexer ClusterIndexer)
	GetByIndex(indexName, key string) ([]*v1.Cluster, error)
}

type ClusterIndexer func(obj *v1.Cluster) ([]string, error)

type ClusterGenericController struct {
	generic.ControllerInterface[*v1.Cluster, *v1.ClusterList]
}

func (c *ClusterGenericController) OnChange(ctx context.Context, name string, sync ClusterHandler) {
	c.ControllerInterface.OnChange(ctx, name, generic.ObjectHandler[*v1.Cluster](sync))
}

func (c *ClusterGenericController) OnRemove(ctx context.Context, name string, sync ClusterHandler) {
	c.ControllerInterface.OnRemove(ctx, name, generic.ObjectHandler[*v1.Cluster](sync))
}

func (c *ClusterGenericController) Cache() ClusterCache {
	return &ClusterGenericCache{
		c.ControllerInterface.Cache(),
	}
}

type ClusterGenericCache struct {
	generic.CacheInterface[*v1.Cluster]
}

func (c ClusterGenericCache) AddIndexer(indexName string, indexer ClusterIndexer) {
	c.CacheInterface.AddIndexer(indexName, generic.Indexer[*v1.Cluster](indexer))
}

type ClusterStatusHandler func(obj *v1.Cluster, status v1.ClusterStatus) (v1.ClusterStatus, error)

type ClusterGeneratingHandler func(obj *v1.Cluster, status v1.ClusterStatus) ([]runtime.Object, v1.ClusterStatus, error)

func FromClusterHandlerToHandler(sync ClusterHandler) generic.Handler {
	return generic.FromObjectHandlerToHandler(generic.ObjectHandler[*v1.Cluster](sync))
}

func RegisterClusterStatusHandler(ctx context.Context, controller ClusterController, condition condition.Cond, name string, handler ClusterStatusHandler) {
	statusHandler := &clusterStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromClusterHandlerToHandler(statusHandler.sync))
}

func RegisterClusterGeneratingHandler(ctx context.Context, controller ClusterController, apply apply.Apply,
	condition condition.Cond, name string, handler ClusterGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &clusterGeneratingHandler{
		ClusterGeneratingHandler: handler,
		apply:                    apply,
		name:                     name,
		gvk:                      controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterClusterStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type clusterStatusHandler struct {
	client    ClusterClient
	condition condition.Cond
	handler   ClusterStatusHandler
}

func (a *clusterStatusHandler) sync(key string, obj *v1.Cluster) (*v1.Cluster, error) {
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

type clusterGeneratingHandler struct {
	ClusterGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *clusterGeneratingHandler) Remove(key string, obj *v1.Cluster) (*v1.Cluster, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1.Cluster{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *clusterGeneratingHandler) Handle(obj *v1.Cluster, status v1.ClusterStatus) (v1.ClusterStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.ClusterGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
