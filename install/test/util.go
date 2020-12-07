package test

import (
	glooTest "github.com/solo-io/gloo/install/test"
	"github.com/solo-io/k8s-utils/manifesttestutils"
)

func GetGlooEServiceAccountPermissions(namespace string) *manifesttestutils.ServiceAccountPermissions {
	// build off of the permissions imported from Gloo
	permissions := glooTest.GetServiceAccountPermissions(namespace)
	ApplyPermissionsForReadOnlyApiserver(namespace, permissions)
	ApplyPermissionsAddedForMutableApiserver(namespace, permissions)
	ApplyPermissionsForGlooEServiceAccounts(namespace, permissions)
	return permissions
}

func GetGlooWithReadOnlyUiServiceAccountPermissions(namespace string) *manifesttestutils.ServiceAccountPermissions {
	// build off of the permissions imported from Gloo
	permissions := glooTest.GetServiceAccountPermissions(namespace)
	ApplyPermissionsForReadOnlyApiserver(namespace, permissions)
	return permissions
}

func ApplyPermissionsForReadOnlyApiserver(namespace string, permissions *manifesttestutils.ServiceAccountPermissions) {
	// Apiserver
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{""},
		[]string{"pods", "services", "configmaps", "namespaces", "secrets"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{"gloo.solo.io"},
		[]string{"artifacts", "upstreams", "upstreamgroups", "proxies", "secrets"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{"gloo.solo.io"},
		[]string{"settings"},
		[]string{"get", "list", "watch", "create"})
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{"gateway.solo.io"},
		[]string{"gateways"},
		[]string{"get", "list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{"gateway.solo.io"},
		[]string{"virtualservices", "routetables"},
		[]string{"get", "list", "watch"})
}

func ApplyPermissionsAddedForMutableApiserver(namespace string, permissions *manifesttestutils.ServiceAccountPermissions) {
	// Apiserver
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{""},
		[]string{"secrets"},
		[]string{"create", "update", "delete"})
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{"gloo.solo.io"},
		[]string{"artifacts", "upstreams", "upstreamgroups", "proxies", "secrets"},
		[]string{"create", "delete", "update"})
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{"gloo.solo.io"},
		[]string{"settings"},
		[]string{"update", "delete"})
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{"gateway.solo.io"},
		[]string{"gateways"},
		[]string{"create", "update", "delete"})
	permissions.AddExpectedPermission(
		"gloo-system.apiserver-ui",
		namespace,
		[]string{"gateway.solo.io"},
		[]string{"virtualservices", "routetables"},
		[]string{"create", "update", "delete"})

}

func ApplyPermissionsForGlooEServiceAccounts(namespace string, permissions *manifesttestutils.ServiceAccountPermissions) {
	// Observability
	permissions.AddExpectedPermission(
		"gloo-system.observability",
		namespace,
		[]string{"gloo.solo.io"},
		[]string{"settings"},
		[]string{"create", "get", "list", "watch"})
	permissions.AddExpectedPermission(
		"gloo-system.observability",
		namespace,
		[]string{"gloo.solo.io"},
		[]string{"upstreams"},
		[]string{"get", "list", "watch"})
}
