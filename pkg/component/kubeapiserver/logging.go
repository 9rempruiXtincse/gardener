// Copyright 2021 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubeapiserver

import (
	"fmt"

	fluentbitv1alpha2 "github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2"
	fluentbitv1alpha2filter "github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/filter"
	fluentbitv1alpha2parser "github.com/fluent/fluent-operator/v2/apis/fluentbit/v1alpha2/plugins/parser"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/pointer"

	v1beta1constants "github.com/gardener/gardener/pkg/apis/core/v1beta1/constants"
	"github.com/gardener/gardener/pkg/component"
)

// CentralLoggingConfiguration returns a fluent-bit parser and filter for the kube-apiserver logs.
func CentralLoggingConfiguration() (component.CentralLoggingConfig, error) {
	return component.CentralLoggingConfig{Filters: generateClusterFilters(), Parsers: generateClusterParsers()}, nil
}

func generateClusterFilters() []*fluentbitv1alpha2.ClusterFilter {
	return []*fluentbitv1alpha2.ClusterFilter{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:   fmt.Sprintf("%s--%s", v1beta1constants.DeploymentNameKubeAPIServer, ContainerNameKubeAPIServer),
				Labels: map[string]string{v1beta1constants.LabelKeyCustomLoggingResource: v1beta1constants.LabelValueCustomLoggingResource},
			},
			Spec: fluentbitv1alpha2.FilterSpec{
				Match: fmt.Sprintf("kubernetes.*%s*%s*", v1beta1constants.DeploymentNameKubeAPIServer, ContainerNameKubeAPIServer),
				FilterItems: []fluentbitv1alpha2.FilterItem{
					{
						Parser: &fluentbitv1alpha2filter.Parser{
							KeyName:     "log",
							Parser:      ContainerNameKubeAPIServer + "-parser",
							ReserveData: pointer.Bool(true),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:   fmt.Sprintf("%s--%s", v1beta1constants.DeploymentNameKubeAPIServer, containerNameAPIServerProxyPodMutator),
				Labels: map[string]string{v1beta1constants.LabelKeyCustomLoggingResource: v1beta1constants.LabelValueCustomLoggingResource},
			},
			Spec: fluentbitv1alpha2.FilterSpec{
				Match: fmt.Sprintf("kubernetes.*%s*%s*", v1beta1constants.DeploymentNameKubeAPIServer, containerNameAPIServerProxyPodMutator),
				FilterItems: []fluentbitv1alpha2.FilterItem{
					{
						Parser: &fluentbitv1alpha2filter.Parser{
							KeyName:     "log",
							Parser:      containerNameAPIServerProxyPodMutator + "-parser",
							ReserveData: pointer.Bool(true),
						},
					},
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:   containerNameAPIServerProxyPodMutator + "-modify-severity",
				Labels: map[string]string{v1beta1constants.LabelKeyCustomLoggingResource: v1beta1constants.LabelValueCustomLoggingResource},
			},
			Spec: fluentbitv1alpha2.FilterSpec{
				Match: fmt.Sprintf("kubernetes.*%s*%s*", v1beta1constants.DeploymentNameKubeAPIServer, containerNameAPIServerProxyPodMutator),
				FilterItems: []fluentbitv1alpha2.FilterItem{
					{
						Modify: &fluentbitv1alpha2filter.Modify{
							Rules: []fluentbitv1alpha2filter.Rule{
								{
									Copy: map[string]string{"level": "severity"},
								},
							},
						},
					},
				},
			},
		},
	}
}

func generateClusterParsers() []*fluentbitv1alpha2.ClusterParser {
	return []*fluentbitv1alpha2.ClusterParser{
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:   ContainerNameKubeAPIServer + "-parser",
				Labels: map[string]string{v1beta1constants.LabelKeyCustomLoggingResource: v1beta1constants.LabelValueCustomLoggingResource},
			},
			Spec: fluentbitv1alpha2.ParserSpec{
				Regex: &fluentbitv1alpha2parser.Regex{
					Regex:      "^(?<severity>\\w)(?<time>\\d{4} [^\\s]*)\\s+(?<pid>\\d+)\\s+(?<source>[^ \\]]+)\\] (?<log>.*)$",
					TimeKey:    "time",
					TimeFormat: "%m%d %H:%M:%S.%L",
				},
			},
		},
		{
			ObjectMeta: metav1.ObjectMeta{
				Name:   containerNameAPIServerProxyPodMutator + "-parser",
				Labels: map[string]string{v1beta1constants.LabelKeyCustomLoggingResource: v1beta1constants.LabelValueCustomLoggingResource},
			},
			Spec: fluentbitv1alpha2.ParserSpec{
				JSON: &fluentbitv1alpha2parser.JSON{
					TimeKey: "ts",
				},
			},
		},
	}
}
