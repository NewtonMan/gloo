package kube

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

type NamespaceClient interface {
	ListNamespaces(ctx context.Context) ([]string, error)
}

type namespaceClient struct {
	namespacesGetter v1.NamespacesGetter
}

func (n *namespaceClient) ListNamespaces(ctx context.Context) ([]string, error) {
	var namespaces []string
	nsList, err := n.namespacesGetter.Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, ns := range nsList.Items {
		namespaces = append(namespaces, ns.Name)
	}
	return namespaces, nil
}

func NewNamespaceClient(namespacesGetter v1.NamespacesGetter) NamespaceClient {
	return &namespaceClient{
		namespacesGetter: namespacesGetter,
	}
}
