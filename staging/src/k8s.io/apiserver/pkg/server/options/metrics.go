/*
Copyright 2019 The Kubernetes Authors.

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

package options

import (
	"github.com/spf13/pflag"
	"k8s.io/component-base/metrics"
)

// MetricsOptions holds the api server metrics options.
type MetricsOptions struct {
	// ShowHiddenMetrics is the version(x.y) for which you want to show hidden metrics.
	ShowHiddenMetrics string
}

// NewMetricsOptions creates a new instance of MetricsOptions
func NewMetricsOptions() *MetricsOptions {
	return &MetricsOptions{}
}

// AddFlags adds flags related to metrics for a specific APIServer to the specified FlagSet
func (o *MetricsOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}

	fs.StringVar(&o.ShowHiddenMetrics, "show-hidden-metrics", o.ShowHiddenMetrics,
		"The kube-apiserver version(x.y) for which you want to show hidden metrics. Only previous minor version is allowed.")
}

// Validate verifies flags passed to MetricsOptions.
func (o *MetricsOptions) Validate() []error {
	if o == nil {
		return nil
	}

	var errs []error

	if o.ShowHiddenMetrics != "" {
		errs = append(errs, metrics.ValidateShowHiddenMetricsVersion(o.ShowHiddenMetrics)...)
	}

	return errs
}
