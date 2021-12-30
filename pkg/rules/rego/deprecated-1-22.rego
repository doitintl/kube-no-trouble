package deprecated122

main[return] {
	resource := input[_]
	api := deprecated_resource(resource)
	return := {
		"Name": get_default(resource.metadata, "name", "<undefined>"),
		# Namespace does not have to be defined in case of local manifests
		"Namespace": get_default(resource.metadata, "namespace", "<undefined>"),
		"Kind": resource.kind,
		"ApiVersion": api.old,
		"ReplaceWith": api.new,
		"RuleSet": "Deprecated APIs removed in 1.22",
		"Since": api.since,
	}
}

deprecated_resource(r) = api {
	api := deprecated_api(r.kind, r.apiVersion)
}

deprecated_api(kind, api_version) = api {
	deprecated_apis = {
		"Ingress": {
			"old": ["extensions/v1beta1", "networking.k8s.io/v1beta1"],
			"new": "networking.k8s.io/v1",
			"since": "1.19",
		},
		"IngressClass": {
			"old": ["networking.k8s.io/v1beta1"],
			"new": "networking.k8s.io/v1",
			"since": "1.19",
		},
		"TokenReview": {
			"old": ["authentication.k8s.io/v1beta1"],
			"new": "authentication.k8s.io/v1",
			"since": "1.19",
		},
		"SubjectAccessReview": {
			"old": ["authorization.k8s.io/v1beta1"],
			"new": "authorization.k8s.io/v1",
			"since": "1.19",
		},
		"SelfSubjectAccessReview": {
			"old": ["authorization.k8s.io/v1beta1"],
			"new": "authorization.k8s.io/v1",
			"since": "1.19",
		},
		"LocalSubjectAccessReview": {
			"old": ["authorization.k8s.io/v1beta1"],
			"new": "authorization.k8s.io/v1",
			"since": "1.19",
		},
		"Lease": {
			"old": ["coordination.k8s.io/v1beta1"],
			"new": "coordination.k8s.io/v1",
			"since": "1.14",
		},
		"PriorityClass": {
			"old": ["scheduling.k8s.io/v1beta1"],
			"new": "scheduling.k8s.io/v1",
			"since": "1.14",
		},
		"CSIDriver": {
			"old": ["storage.k8s.io/v1beta1"],
			"new": "storage.k8s.io/v1",
			"since": "1.19",
		},
		"CSINode": {
			"old": ["storage.k8s.io/v1beta1"],
			"new": "storage.k8s.io/v1",
			"since": "1.17",
		},
		"StorageClass": {
			"old": ["storage.k8s.io/v1beta1"],
			"new": "storage.k8s.io/v1",
			"since": "1.6",
		},
		"VolumeAttachment": {
			"old": ["storage.k8s.io/v1beta1"],
			"new": "storage.k8s.io/v1",
			"since": "1.13",
		},
		"ClusterRole": {
			"old": ["rbac.authorization.k8s.io/v1beta1"],
			"new": "rbac.authorization.k8s.io/v1",
			"since": "1.8",
		},
		"ClusterRoleBinding": {
			"old": ["rbac.authorization.k8s.io/v1beta1"],
			"new": "rbac.authorization.k8s.io/v1",
			"since": "1.8",
		},
		"Role": {
			"old": ["rbac.authorization.k8s.io/v1beta1"],
			"new": "rbac.authorization.k8s.io/v1",
			"since": "1.8",
		},
		"RoleBinding": {
			"old": ["rbac.authorization.k8s.io/v1beta1"],
			"new": "rbac.authorization.k8s.io/v1",
			"since": "1.8",
		},
		"CertificateSigningRequest": {
			"old": ["certificates.k8s.io/v1beta1"],
			"new": "certificates.k8s.io/v1",
			"since": "1.19",
		},
		"APIService": {
			"old": ["apiregistration.k8s.io/v1beta1"],
			"new": "apiregistration.k8s.io/v1",
			"since": "1.10",
		},
		"CustomResourceDefinition": {
			"old": ["apiextensions.k8s.io/v1beta1"],
			"new": "apiextensions.k8s.io/v1",
			"since": "1.16",
		},
		"MutatingWebhookConfiguration": {
			"old": ["admissionregistration.k8s.io/v1beta1"],
			"new": "admissionregistration.k8s.io/v1",
			"since": "1.16",
		},
		"ValidatingWebhookConfiguration": {
			"old": ["admissionregistration.k8s.io/v1beta1"],
			"new": "admissionregistration.k8s.io/v1",
			"since": "1.16",
		},
	}

	deprecated_apis[kind].old[_] == api_version

	api := {
		"old": api_version,
		"new": deprecated_apis[kind].new,
		"since": deprecated_apis[kind].since,
	}
}

get_default(val, key, _) = val[key]

get_default(val, key, fallback) = fallback {
	not val[key]
}
