/*
Copyright 2022.

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
/*
 * SPDX-FileCopyrightText: 2024 Samir Zeort <samir.zeort@sap.com>
 *
 * SPDX-License-Identifier: Apache-2.0
 */

// Package v1alpha1 contains API Schema definitions for the component v1alpha1 API group
// +kubebuilder:object:generate=true
// +groupName=operator.kyma-project.io
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +kubebuilder:object:root=true

type Managed struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
}

// +kubebuilder:object:root=true

// ManagedList contains a list of Managed.
type ManagedList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Managed `json:"items"`
}
