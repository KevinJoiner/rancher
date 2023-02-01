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

type UserHandler func(string, *v3.User) (*v3.User, error)

type UserController interface {
	generic.ControllerMeta
	UserClient

	OnChange(ctx context.Context, name string, sync UserHandler)
	OnRemove(ctx context.Context, name string, sync UserHandler)
	Enqueue(name string)
	EnqueueAfter(name string, duration time.Duration)

	Cache() UserCache
}

type UserClient interface {
	Create(*v3.User) (*v3.User, error)
	Update(*v3.User) (*v3.User, error)
	UpdateStatus(*v3.User) (*v3.User, error)
	Delete(name string, options *metav1.DeleteOptions) error
	Get(name string, options metav1.GetOptions) (*v3.User, error)
	List(opts metav1.ListOptions) (*v3.UserList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v3.User, err error)
}

type UserCache interface {
	Get(name string) (*v3.User, error)
	List(selector labels.Selector) ([]*v3.User, error)

	AddIndexer(indexName string, indexer UserIndexer)
	GetByIndex(indexName, key string) ([]*v3.User, error)
}

type UserIndexer func(obj *v3.User) ([]string, error)

type UserGenericController struct {
	generic.NonNamespacedControllerInterface[*v3.User, *v3.UserList]
}

func (c *UserGenericController) OnChange(ctx context.Context, name string, sync UserHandler) {
	c.NonNamespacedControllerInterface.OnChange(ctx, name, generic.ObjectHandler[*v3.User](sync))
}

func (c *UserGenericController) OnRemove(ctx context.Context, name string, sync UserHandler) {
	c.NonNamespacedControllerInterface.OnRemove(ctx, name, generic.ObjectHandler[*v3.User](sync))
}

func (c *UserGenericController) Cache() UserCache {
	return &UserGenericCache{
		c.NonNamespacedControllerInterface.Cache(),
	}
}

type UserGenericCache struct {
	generic.NonNamespacedCacheInterface[*v3.User]
}

func (c UserGenericCache) AddIndexer(indexName string, indexer UserIndexer) {
	c.NonNamespacedCacheInterface.AddIndexer(indexName, generic.Indexer[*v3.User](indexer))
}

type UserStatusHandler func(obj *v3.User, status v3.UserStatus) (v3.UserStatus, error)

type UserGeneratingHandler func(obj *v3.User, status v3.UserStatus) ([]runtime.Object, v3.UserStatus, error)

func FromUserHandlerToHandler(sync UserHandler) generic.Handler {
	return generic.FromObjectHandlerToHandler(generic.ObjectHandler[*v3.User](sync))
}

func RegisterUserStatusHandler(ctx context.Context, controller UserController, condition condition.Cond, name string, handler UserStatusHandler) {
	statusHandler := &userStatusHandler{
		client:    controller,
		condition: condition,
		handler:   handler,
	}
	controller.AddGenericHandler(ctx, name, FromUserHandlerToHandler(statusHandler.sync))
}

func RegisterUserGeneratingHandler(ctx context.Context, controller UserController, apply apply.Apply,
	condition condition.Cond, name string, handler UserGeneratingHandler, opts *generic.GeneratingHandlerOptions) {
	statusHandler := &userGeneratingHandler{
		UserGeneratingHandler: handler,
		apply:                 apply,
		name:                  name,
		gvk:                   controller.GroupVersionKind(),
	}
	if opts != nil {
		statusHandler.opts = *opts
	}
	controller.OnChange(ctx, name, statusHandler.Remove)
	RegisterUserStatusHandler(ctx, controller, condition, name, statusHandler.Handle)
}

type userStatusHandler struct {
	client    UserClient
	condition condition.Cond
	handler   UserStatusHandler
}

func (a *userStatusHandler) sync(key string, obj *v3.User) (*v3.User, error) {
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

type userGeneratingHandler struct {
	UserGeneratingHandler
	apply apply.Apply
	opts  generic.GeneratingHandlerOptions
	gvk   schema.GroupVersionKind
	name  string
}

func (a *userGeneratingHandler) Remove(key string, obj *v3.User) (*v3.User, error) {
	if obj != nil {
		return obj, nil
	}

	obj = &v3.User{}
	obj.Namespace, obj.Name = kv.RSplit(key, "/")
	obj.SetGroupVersionKind(a.gvk)

	return nil, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects()
}

func (a *userGeneratingHandler) Handle(obj *v3.User, status v3.UserStatus) (v3.UserStatus, error) {
	if !obj.DeletionTimestamp.IsZero() {
		return status, nil
	}

	objs, newStatus, err := a.UserGeneratingHandler(obj, status)
	if err != nil {
		return newStatus, err
	}

	return newStatus, generic.ConfigureApplyForObject(a.apply, obj, &a.opts).
		WithOwner(obj).
		WithSetID(a.name).
		ApplyObjects(objs...)
}
