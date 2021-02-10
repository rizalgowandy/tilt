/*
Copyright 2020 The Tilt Dev Authors

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

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	"time"

	v1alpha1 "github.com/tilt-dev/tilt/pkg/apis/core/v1alpha1"
	scheme "github.com/tilt-dev/tilt/pkg/clientset/tiltapi/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// ManifestsGetter has a method to return a ManifestInterface.
// A group's client should implement this interface.
type ManifestsGetter interface {
	Manifests() ManifestInterface
}

// ManifestInterface has methods to work with Manifest resources.
type ManifestInterface interface {
	Create(ctx context.Context, manifest *v1alpha1.Manifest, opts v1.CreateOptions) (*v1alpha1.Manifest, error)
	Update(ctx context.Context, manifest *v1alpha1.Manifest, opts v1.UpdateOptions) (*v1alpha1.Manifest, error)
	UpdateStatus(ctx context.Context, manifest *v1alpha1.Manifest, opts v1.UpdateOptions) (*v1alpha1.Manifest, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.Manifest, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.ManifestList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Manifest, err error)
	ManifestExpansion
}

// manifests implements ManifestInterface
type manifests struct {
	client rest.Interface
}

// newManifests returns a Manifests
func newManifests(c *CoreV1alpha1Client) *manifests {
	return &manifests{
		client: c.RESTClient(),
	}
}

// Get takes name of the manifest, and returns the corresponding manifest object, and an error if there is any.
func (c *manifests) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Manifest, err error) {
	result = &v1alpha1.Manifest{}
	err = c.client.Get().
		Resource("manifests").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Manifests that match those selectors.
func (c *manifests) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ManifestList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.ManifestList{}
	err = c.client.Get().
		Resource("manifests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested manifests.
func (c *manifests) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("manifests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a manifest and creates it.  Returns the server's representation of the manifest, and an error, if there is any.
func (c *manifests) Create(ctx context.Context, manifest *v1alpha1.Manifest, opts v1.CreateOptions) (result *v1alpha1.Manifest, err error) {
	result = &v1alpha1.Manifest{}
	err = c.client.Post().
		Resource("manifests").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(manifest).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a manifest and updates it. Returns the server's representation of the manifest, and an error, if there is any.
func (c *manifests) Update(ctx context.Context, manifest *v1alpha1.Manifest, opts v1.UpdateOptions) (result *v1alpha1.Manifest, err error) {
	result = &v1alpha1.Manifest{}
	err = c.client.Put().
		Resource("manifests").
		Name(manifest.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(manifest).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *manifests) UpdateStatus(ctx context.Context, manifest *v1alpha1.Manifest, opts v1.UpdateOptions) (result *v1alpha1.Manifest, err error) {
	result = &v1alpha1.Manifest{}
	err = c.client.Put().
		Resource("manifests").
		Name(manifest.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(manifest).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the manifest and deletes it. Returns an error if one occurs.
func (c *manifests) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Resource("manifests").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *manifests) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("manifests").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched manifest.
func (c *manifests) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Manifest, err error) {
	result = &v1alpha1.Manifest{}
	err = c.client.Patch(pt).
		Resource("manifests").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
