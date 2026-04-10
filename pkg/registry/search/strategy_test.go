/*
Copyright 2026 The Karmada Authors.

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

package search

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	policyv1alpha1 "github.com/karmada-io/karmada/pkg/apis/policy/v1alpha1"
	searchapis "github.com/karmada-io/karmada/pkg/apis/search"
	searchscheme "github.com/karmada-io/karmada/pkg/apis/search/scheme"
)

func TestStrategyValidate(t *testing.T) {
	strategy := NewStrategy(searchscheme.Scheme)
	tests := []struct {
		name             string
		resourceRegistry *searchapis.ResourceRegistry
		wantErr          bool
	}{
		{
			name:             "valid resource registry",
			resourceRegistry: validResourceRegistry("registry-a"),
			wantErr:          false,
		},
		{
			name: "invalid resource registry selector api version",
			resourceRegistry: &searchapis.ResourceRegistry{
				ObjectMeta: metav1.ObjectMeta{Name: "registry-a"},
				Spec: searchapis.ResourceRegistrySpec{
					TargetCluster: policyv1alpha1.ClusterAffinity{},
					ResourceSelectors: []searchapis.ResourceSelector{
						{APIVersion: "apps/v1/extra", Kind: "Deployment"},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := strategy.Validate(context.Background(), tt.resourceRegistry)
			if tt.wantErr && len(errs) == 0 {
				t.Fatal("expected validation errors, got none")
			}
			if !tt.wantErr && len(errs) != 0 {
				t.Fatalf("expected no validation errors, got: %v", errs)
			}
		})
	}
}

func TestStrategyValidateUpdate(t *testing.T) {
	strategy := NewStrategy(searchscheme.Scheme)
	oldResourceRegistry := validResourceRegistry("registry-a")

	tests := []struct {
		name                string
		newResourceRegistry *searchapis.ResourceRegistry
		wantErr             bool
	}{
		{
			name:                "valid update",
			newResourceRegistry: validResourceRegistry("registry-a"),
			wantErr:             false,
		},
		{
			name:                "invalid update changing metadata name",
			newResourceRegistry: validResourceRegistry("registry-b"),
			wantErr:             true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := strategy.ValidateUpdate(context.Background(), tt.newResourceRegistry, oldResourceRegistry)
			if tt.wantErr && len(errs) == 0 {
				t.Fatal("expected validation errors, got none")
			}
			if !tt.wantErr && len(errs) != 0 {
				t.Fatalf("expected no validation errors, got: %v", errs)
			}
		})
	}
}

func validResourceRegistry(name string) *searchapis.ResourceRegistry {
	return &searchapis.ResourceRegistry{
		ObjectMeta: metav1.ObjectMeta{Name: name, ResourceVersion: "1"},
		Spec: searchapis.ResourceRegistrySpec{
			TargetCluster: policyv1alpha1.ClusterAffinity{},
			ResourceSelectors: []searchapis.ResourceSelector{
				{APIVersion: "apps/v1", Kind: "Deployment", Namespace: "default"},
			},
		},
	}
}
