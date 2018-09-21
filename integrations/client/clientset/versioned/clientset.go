/*
Copyright 2018 Weaveworks Ltd.

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
package versioned

import (
	fluxv1beta1 "github.com/weaveworks/flux/integrations/client/clientset/versioned/typed/flux.weave.works/v1beta1"
	helmv1alpha2 "github.com/weaveworks/flux/integrations/client/clientset/versioned/typed/helm.integrations.flux.weave.works/v1alpha2"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	FluxV1beta1() fluxv1beta1.FluxV1beta1Interface
	// Deprecated: please explicitly pick a version if possible.
	Flux() fluxv1beta1.FluxV1beta1Interface
	HelmV1alpha2() helmv1alpha2.HelmV1alpha2Interface
	// Deprecated: please explicitly pick a version if possible.
	Helm() helmv1alpha2.HelmV1alpha2Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	fluxV1beta1  *fluxv1beta1.FluxV1beta1Client
	helmV1alpha2 *helmv1alpha2.HelmV1alpha2Client
}

// FluxV1beta1 retrieves the FluxV1beta1Client
func (c *Clientset) FluxV1beta1() fluxv1beta1.FluxV1beta1Interface {
	return c.fluxV1beta1
}

// Deprecated: Flux retrieves the default version of FluxClient.
// Please explicitly pick a version.
func (c *Clientset) Flux() fluxv1beta1.FluxV1beta1Interface {
	return c.fluxV1beta1
}

// HelmV1alpha2 retrieves the HelmV1alpha2Client
func (c *Clientset) HelmV1alpha2() helmv1alpha2.HelmV1alpha2Interface {
	return c.helmV1alpha2
}

// Deprecated: Helm retrieves the default version of HelmClient.
// Please explicitly pick a version.
func (c *Clientset) Helm() helmv1alpha2.HelmV1alpha2Interface {
	return c.helmV1alpha2
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.fluxV1beta1, err = fluxv1beta1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.helmV1alpha2, err = helmv1alpha2.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.fluxV1beta1 = fluxv1beta1.NewForConfigOrDie(c)
	cs.helmV1alpha2 = helmv1alpha2.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.fluxV1beta1 = fluxv1beta1.New(c)
	cs.helmV1alpha2 = helmv1alpha2.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
