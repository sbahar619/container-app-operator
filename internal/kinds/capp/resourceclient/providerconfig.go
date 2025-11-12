package resourceclient

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	providerConfigGroup   = "dns-v2.m.crossplane.io"
	providerConfigVersion = "v1beta1"
	ProviderConfigKind    = "ProviderConfig"
	ClusterProviderConfigKind    = "ClusterProviderConfig"
	providerConfigUsageListKind = "ProviderConfigUsageList"
	providerConfigName    = "default"
	secretNamespace       = "crossplane-system"
	secretName            = "dns-creds"
	secretKey             = "credentials"
)

// newProviderConfig returns an Unstructured ProviderConfig object for the given namespace.
func newProviderConfig(namespace string) *unstructured.Unstructured {
	pc := &unstructured.Unstructured{}
	pc.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   providerConfigGroup,
		Version: providerConfigVersion,
		Kind:    ProviderConfigKind,
	})
	pc.SetNamespace(namespace)
	pc.SetName(providerConfigName)
	return pc
}

// desiredSecretRef returns the desired secretRef subtree for the ProviderConfig spec.
func desiredSecretRef() map[string]any {
	return map[string]any{
		"namespace": secretNamespace,
		"name":      secretName,
		"key":       secretKey,
	}
}

// ensureDesiredSecretRef sets the desired secretRef in the given ProviderConfig if it differs.
func ensureDesiredSecretRef(pc *unstructured.Unstructured) bool {
	current, _, _ := unstructured.NestedMap(pc.Object, "spec", "credentials", "secretRef")
	desired := desiredSecretRef()
	if !equality.Semantic.DeepEqual(current, desired) {
		_ = unstructured.SetNestedMap(pc.Object, desired, "spec", "credentials", "secretRef")
		return true
	}
	return false
}

// EnsureProviderConfig ensures a shared, namespaced ProviderConfig named "default" exists with the desired secretRef.
func EnsureProviderConfig(ctx context.Context, c client.Client, namespace string) error {
	pc := newProviderConfig(namespace)

	err := c.Get(ctx, client.ObjectKey{Namespace: namespace, Name: providerConfigName}, pc)
	if apierrors.IsNotFound(err) {
		desired := map[string]any{
			"credentials": map[string]any{
				"source": "Secret",
				"secretRef": desiredSecretRef(),
			},
		}
		_ = unstructured.SetNestedMap(pc.Object, desired, "spec")
		return c.Create(ctx, pc)
	}
	if err != nil {
		return err
	}

	if ensureDesiredSecretRef(pc) {
		return c.Update(ctx, pc)
	}
	return nil
}

// DeleteProviderConfigIfUnused deletes the namespaced shared ProviderConfig only if it is unused in the namespace.
func DeleteProviderConfigIfUnused(ctx context.Context, c client.Client, namespace string) error {
	usages := &unstructured.UnstructuredList{}
	usages.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   providerConfigGroup,
		Version: providerConfigVersion,
		Kind:    providerConfigUsageListKind,
	})
	if err := c.List(ctx, usages, client.InNamespace(namespace)); err != nil {
		return err
	}
	for _, u := range usages.Items {
		name, _, _ := unstructured.NestedString(u.Object, "spec", "providerConfigRef", "name")
		if name == providerConfigName {
			return nil
		}
	}

	if err := c.Delete(ctx, newProviderConfig(namespace)); err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
		return err
	}
	return nil
}


