package resourcemanagers

import (
	"context"

	cappv1alpha1 "github.com/dana-team/container-app-operator/api/v1alpha1"
	rclient "github.com/dana-team/container-app-operator/internal/kinds/capp/resourceclient"
	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	ProviderConfig = "providerConfig"
)

type ProviderConfigManager struct {
	Ctx           context.Context
	K8sclient     client.Client
	Log           logr.Logger
	EventRecorder interface{}
}

func (r ProviderConfigManager) Manage(capp cappv1alpha1.Capp) error {
	return rclient.EnsureProviderConfig(r.Ctx, r.K8sclient, capp.Namespace)
}

func (r ProviderConfigManager) CleanUp(capp cappv1alpha1.Capp) error {
	return rclient.DeleteProviderConfigIfUnused(r.Ctx, r.K8sclient, capp.Namespace)
}

func (r ProviderConfigManager) IsRequired(capp cappv1alpha1.Capp) bool {
	return true
}

