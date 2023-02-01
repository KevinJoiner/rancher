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

package v1beta1

import (
	"context"
	"time"

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
	v1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

type MachineSetHandler func(string, *v1beta1.MachineSet) (*v1beta1.MachineSet, error)

type MachineSetController interface {
	generic.ControllerMeta
	MachineSetClient

	OnChange(ctx context.Context, name string, sync MachineSetHandler)
	OnRemove(ctx context.Context, name string, sync MachineSetHandler)
	Enqueue(namespace, name string)
	EnqueueAfter(namespace, name string, duration time.Duration)

	Cache() MachineSetCache
}

type MachineSetClient interface {
	Create(*v1beta1.MachineSet) (*v1beta1.MachineSet, error)
	Update(*v1beta1.MachineSet) (*v1beta1.MachineSet, error)
	UpdateStatus(*v1beta1.MachineSet) (*v1beta1.MachineSet, error)
	Delete(namespace, name string, options *metav1.DeleteOptions) error
	Get(namespace, name string, options metav1.GetOptions) (*v1beta1.MachineSet, error)
	List(namespace string, opts metav1.ListOptions) (*v1beta1.MachineSetList, error)
	Watch(namespace string, opts metav1.ListOptions) (watch.Interface, error)
	Patch(namespace, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.MachineSet, err error)
}

type MachineSetCache interface {
	Get(namespace, name string) (*v1beta1.MachineSet, error)
	List(namespace string, selector labels.Selector) ([]*v1beta1.MachineSet, error)

	AddIndexer(indexName string, indexer MachineSetIndexer)
	GetByIndex(indexName, key string) ([]*v1beta1.MachineSet, error)
}

type MachineSetIndexer func(obj *v1beta1.MachineSet) ([]string, error)

type MachineSetGenericController struct {
	generic.ControllerInterface[*v1beta1.MachineSet, *v1beta1.MachineSetList]
}

func (c *MachineSetGenericController) OnChange(ctx context.Context, name string, sync MachineSetHandler) {
	c.ControllerInterface.OnChange(ctx, name, generic.ObjectHandler[*v1beta1.MachineSet](sync))
}

func (c *MachineSetGenericController) OnRemove(ctx context.Context, name string, sync MachineSetHandler) {
	c.ControllerInterface.OnRemove(ctx, name, generic.ObjectHandler[*v1beta1.MachineSet](sync))
}

func (c *MachineSetGenericController) Cache() MachineSetCache {
	return &MachineSetGenericCache{
		c.ControllerInterface.Cache(),
	}
}

type MachineSetGenericCache struct {
	generic.CacheInterface[*v1beta1.MachineSet]
}

func (c MachineSetGenericCache) AddIndexer(indexName string, indexer MachineSetIndexer) {
	c.CacheInterface.AddIndexer(indexName, generic.Indexer[*v1beta1.MachineSet](indexer))
}

type MachineSetStatusHandler func(obj *v1beta1.MachineSet, status v1beta1.MachineSetStatus) (v1beta1.MachineSetStatus, error)

type MachineSetGeneratingHandler func(obj *v1beta1.MachineSet, status v1beta1.MachineSetStatus) ([]runtime.Object, v1beta1.MachineSetStatus, error)

func FromMachineSetHandlerToHandler(sync MachineSetHandler) generic.Handler {
	return generic.FromObjectHandlerToHandler(generic.ObjectHandler[*v1beta1.MachineSet](sync))
}

func RegisterMachineSetStatusHandler(ctx context.Context, controller MachineSetController, condition condition.Cond, name string, handler MachineSetStatusHandler) {
	statusHandler := &machineSetStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromMachineSetHandlerToHandler(statusHandler.sync))
}

func RegisterMachineSetGeneratingHandler(ctx context.Context, controller MachineSetController, apply apply.Apply,
	condition condition.Cond, name string, handler MachineSetGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &machineSetGeneratingHandler{
		MachineSetGeneratingHandler: handler,
		apply:                       apply,
		name:                        name,
		gvk:                         controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterMachineSetStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type machineSetStatusHandler struct {
	client    MachineSetClient
	condition condition.Cond
	handler   MachineSetStatusHandler
}

func (a *machineSetStatusHandler) sync(key string, obj *v1beta1.MachineSet) (*v1beta1.MachineSet, error) {
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

type machineSetGeneratingHandler struct {
	MachineSetGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *machineSetGeneratingHandler) Remove(key string, obj *v1beta1.MachineSet) (*v1beta1.MachineSet, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v1beta1.MachineSet{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *machineSetGeneratingHandler) Handle(obj *v1beta1.MachineSet, status v1beta1.MachineSetStatus) (v1beta1.MachineSetStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.MachineSetGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
