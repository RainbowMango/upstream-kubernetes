/*
Copyright 2020 The Kubernetes Authors.

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

// Package metriclint provides a serials functions to lint a metric.
// The implementation base on [promlint](github.com/prometheus/client_golang/prometheus/testutil/promlint).
//
// metriclint provides a ability to lint a metric at the registry which is different with promlint.
// The lint rules also base on promlint but we may add more rules if necessary.
package metriclint

import (
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// LintResult represents lint result of a specific metric.
type LintResult struct {
	// The FQName of a metric.
	MetricName string

	// one or more lint errors of the metric.
	Issues []string
}

func (lr *LintResult) String() string {
	return lr.MetricName + ":" + strings.Join(lr.Issues, ",")
}

func LintCounter(counterOpts prometheus.CounterOpts) *LintResult {
	result := &LintResult{
		MetricName: prometheus.BuildFQName(counterOpts.Namespace, counterOpts.Subsystem, counterOpts.Name),
	}

	result.Issues = append(result.Issues, commonLint(prometheus.Opts(counterOpts))...)

	// lint names
	result.Issues = append(result.Issues, lintNonHistogramNoBucket(result.MetricName)...)
	result.Issues = append(result.Issues, lintNonHistogramSummaryNoCount(result.MetricName)...)
	result.Issues = append(result.Issues, lintNonHistogramSummaryNoSum(result.MetricName)...)
	result.Issues = append(result.Issues, lintCounterContainsTotal(result.MetricName)...)

	// lint labels
	result.Issues = append(result.Issues, lintNonHistogramNoLabelLe(counterOpts.ConstLabels, nil)...)
	result.Issues = append(result.Issues, lintNonSummaryNoLabelQuantile(counterOpts.ConstLabels, nil)...)
	result.Issues = append(result.Issues, lintLabelNameCamelCase(counterOpts.ConstLabels, nil)...)

	return result
}

func LintCounterVector(counterOpts prometheus.CounterOpts, labelNames []string) *LintResult {
	result := LintCounter(counterOpts)
	result.Issues = append(result.Issues, lintNonHistogramNoLabelLe(nil, labelNames)...)
	result.Issues = append(result.Issues, lintNonSummaryNoLabelQuantile(nil, labelNames)...)
	result.Issues = append(result.Issues, lintLabelNameCamelCase(nil, labelNames)...)

	return result
}

func LintGauge(gaugeOpts prometheus.GaugeOpts) *LintResult {
	result := &LintResult{
		MetricName: prometheus.BuildFQName(gaugeOpts.Namespace, gaugeOpts.Subsystem, gaugeOpts.Name),
	}

	result.Issues = append(result.Issues, commonLint(prometheus.Opts(gaugeOpts))...)
	result.Issues = append(result.Issues, lintNonCounterNoTotal(result.MetricName)...)
	result.Issues = append(result.Issues, lintNonHistogramNoBucket(result.MetricName)...)
	result.Issues = append(result.Issues, lintNonHistogramSummaryNoCount(result.MetricName)...)
	result.Issues = append(result.Issues, lintNonHistogramSummaryNoSum(result.MetricName)...)

	// lint labels
	result.Issues = append(result.Issues, lintNonHistogramNoLabelLe(gaugeOpts.ConstLabels, nil)...)
	result.Issues = append(result.Issues, lintNonSummaryNoLabelQuantile(gaugeOpts.ConstLabels, nil)...)
	result.Issues = append(result.Issues, lintLabelNameCamelCase(gaugeOpts.ConstLabels, nil)...)

	return result
}

func LintGaugeVector(gaugeOpts prometheus.GaugeOpts, labelNames []string) *LintResult {
	result := LintGauge(gaugeOpts)

	result.Issues = append(result.Issues, lintNonHistogramNoLabelLe(nil, labelNames)...)
	result.Issues = append(result.Issues, lintNonSummaryNoLabelQuantile(nil, labelNames)...)
	result.Issues = append(result.Issues, lintLabelNameCamelCase(nil, labelNames)...)

	return result
}

func LintHistogram(histogramOpts prometheus.HistogramOpts) *LintResult {
	result := &LintResult{
		MetricName: prometheus.BuildFQName(histogramOpts.Namespace, histogramOpts.Subsystem, histogramOpts.Name),
	}

	result.Issues = append(result.Issues, commonLint(histogramOpts)...)
	result.Issues = append(result.Issues, lintNonCounterNoTotal(result.MetricName)...)

	// lint labels
	result.Issues = append(result.Issues, lintNonSummaryNoLabelQuantile(histogramOpts.ConstLabels, nil)...)
	result.Issues = append(result.Issues, lintLabelNameCamelCase(histogramOpts.ConstLabels, nil)...)

	return result
}

func LintHistogramVector(histogramOpts prometheus.HistogramOpts, labelNames []string) *LintResult {
	result := LintHistogram(histogramOpts)
	result.Issues = append(result.Issues, lintNonSummaryNoLabelQuantile(nil, labelNames)...)
	result.Issues = append(result.Issues, lintLabelNameCamelCase(nil, labelNames)...)

	return result
}

func LintSummary(summaryOpts prometheus.SummaryOpts) *LintResult {
	result := &LintResult{
		MetricName: prometheus.BuildFQName(summaryOpts.Namespace, summaryOpts.Subsystem, summaryOpts.Name),
	}

	result.Issues = append(result.Issues, commonLint(summaryOpts)...)
	result.Issues = append(result.Issues, lintNonCounterNoTotal(result.MetricName)...)
	result.Issues = append(result.Issues, lintNonHistogramNoBucket(result.MetricName)...)

	// lint labels
	result.Issues = append(result.Issues, lintNonHistogramNoLabelLe(summaryOpts.ConstLabels, nil)...)
	result.Issues = append(result.Issues, lintLabelNameCamelCase(summaryOpts.ConstLabels, nil)...)

	return result
}

func LintSummaryVector(summaryOpts prometheus.SummaryOpts, labelNames []string) *LintResult {
	result := LintSummary(summaryOpts)
	result.Issues = append(result.Issues, lintNonHistogramNoLabelLe(nil, labelNames)...)

	result.Issues = append(result.Issues, lintLabelNameCamelCase(nil, labelNames)...)

	return result
}