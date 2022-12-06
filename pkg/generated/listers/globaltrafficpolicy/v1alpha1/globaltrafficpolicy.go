/*
 * MIT License
 *
 * Copyright (c) since 2021,  flomesh.io Authors.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */
// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/flomesh-io/ErieCanal/apis/globaltrafficpolicy/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// GlobalTrafficPolicyLister helps list GlobalTrafficPolicies.
// All objects returned here must be treated as read-only.
type GlobalTrafficPolicyLister interface {
	// List lists all GlobalTrafficPolicies in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.GlobalTrafficPolicy, err error)
	// GlobalTrafficPolicies returns an object that can list and get GlobalTrafficPolicies.
	GlobalTrafficPolicies(namespace string) GlobalTrafficPolicyNamespaceLister
	GlobalTrafficPolicyListerExpansion
}

// globalTrafficPolicyLister implements the GlobalTrafficPolicyLister interface.
type globalTrafficPolicyLister struct {
	indexer cache.Indexer
}

// NewGlobalTrafficPolicyLister returns a new GlobalTrafficPolicyLister.
func NewGlobalTrafficPolicyLister(indexer cache.Indexer) GlobalTrafficPolicyLister {
	return &globalTrafficPolicyLister{indexer: indexer}
}

// List lists all GlobalTrafficPolicies in the indexer.
func (s *globalTrafficPolicyLister) List(selector labels.Selector) (ret []*v1alpha1.GlobalTrafficPolicy, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.GlobalTrafficPolicy))
	})
	return ret, err
}

// GlobalTrafficPolicies returns an object that can list and get GlobalTrafficPolicies.
func (s *globalTrafficPolicyLister) GlobalTrafficPolicies(namespace string) GlobalTrafficPolicyNamespaceLister {
	return globalTrafficPolicyNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// GlobalTrafficPolicyNamespaceLister helps list and get GlobalTrafficPolicies.
// All objects returned here must be treated as read-only.
type GlobalTrafficPolicyNamespaceLister interface {
	// List lists all GlobalTrafficPolicies in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.GlobalTrafficPolicy, err error)
	// Get retrieves the GlobalTrafficPolicy from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.GlobalTrafficPolicy, error)
	GlobalTrafficPolicyNamespaceListerExpansion
}

// globalTrafficPolicyNamespaceLister implements the GlobalTrafficPolicyNamespaceLister
// interface.
type globalTrafficPolicyNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all GlobalTrafficPolicies in the indexer for a given namespace.
func (s globalTrafficPolicyNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.GlobalTrafficPolicy, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.GlobalTrafficPolicy))
	})
	return ret, err
}

// Get retrieves the GlobalTrafficPolicy from the indexer for a given namespace and name.
func (s globalTrafficPolicyNamespaceLister) Get(name string) (*v1alpha1.GlobalTrafficPolicy, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("globaltrafficpolicy"), name)
	}
	return obj.(*v1alpha1.GlobalTrafficPolicy), nil
}
