package deprecatedFuture

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
		"RuleSet": "Deprecated APIs to be removed in future",
		"Since": api.since,
		"Labels": get_default(resource.metadata, "labels", "<undefined>"),
	}
}

deprecated_resource(r) = api {
	api := deprecated_api(r.kind, r.apiVersion)
}

deprecated_api(kind, api_version) = api {
	deprecated_apis = {
		"VolumeSnapshot": {
			"old": ["snapshot.storage.k8s.io/v1beta1"],
			"new": "snapshot.storage.k8s.io/v1",
			"since": "1.21",
		},
		"VolumeSnapshotClass": {
			"old": ["snapshot.storage.k8s.io/v1beta1"],
			"new": "snapshot.storage.k8s.io/v1",
			"since": "1.21",
		},
		"VolumeSnapshotContent": {
			"old": ["snapshot.storage.k8s.io/v1beta1"],
			"new": "snapshot.storage.k8s.io/v1",
			"since": "1.21",
		},
		"SelfSubjectReview": {
			"old": ["authentication.k8s.io/v1beta1"],
			"new": "authentication.k8s.io/v1",
			"since": "1.28",
		},
		"ValidatingAdmissionPolicy": {
			"old": ["admissionregistration.k8s.io/v1beta1"],
			"new": "admissionregistration.k8s.io/v1",
			"since": "1.30",
		},
		"ValidatingAdmissionPolicyBinding": {
			"old": ["admissionregistration.k8s.io/v1beta1"],
			"new": "admissionregistration.k8s.io/v1",
			"since": "1.30",
		},
		"IPAddress": {
			"old": ["networking.k8s.io/v1beta1"],
			"new": "networking.k8s.io/v1",
			"since": "1.33",
		},
		"ServiceCIDR": {
			"old": ["networking.k8s.io/v1beta1"],
			"new": "networking.k8s.io/v1",
			"since": "1.33",
		},
		"VolumeAttributesClass": {
			"old": ["storage.k8s.io/v1beta1"],
			"new": "storage.k8s.io/v1",
			"since": "1.34",
		},
		"DeviceClass": {
			"old": ["resource.k8s.io/v1beta1", "resource.k8s.io/v1beta2"],
			"new": "resource.k8s.io/v1",
			"since": "1.34",
		},
		"ResourceClaim": {
			"old": ["resource.k8s.io/v1beta1", "resource.k8s.io/v1beta2"],
			"new": "resource.k8s.io/v1",
			"since": "1.34",
		},
		"ResourceClaimTemplate": {
			"old": ["resource.k8s.io/v1beta1", "resource.k8s.io/v1beta2"],
			"new": "resource.k8s.io/v1",
			"since": "1.34",
		},
		"ResourceSlice": {
			"old": ["resource.k8s.io/v1beta1", "resource.k8s.io/v1beta2"],
			"new": "resource.k8s.io/v1",
			"since": "1.34",
		},
		"MutatingAdmissionPolicy": {
			"old": ["admissionregistration.k8s.io/v1beta1"],
			"new": "admissionregistration.k8s.io/v1",
			"since": "1.36",
		},
		"MutatingAdmissionPolicyBinding": {
			"old": ["admissionregistration.k8s.io/v1beta1"],
			"new": "admissionregistration.k8s.io/v1",
			"since": "1.36",
		},
		"ClusterTrustBundle": {
			"old": ["certificates.k8s.io/v1beta1"],
			"new": "certificates.k8s.io/v1",
			"since": "1.37",
		},
		"PodCertificateRequest": {
			"old": ["certificates.k8s.io/v1beta1"],
			"new": "certificates.k8s.io/v1",
			"since": "1.37",
		},
		"StorageVersionMigration": {
			"old": ["storagemigration.k8s.io/v1beta1"],
			"new": "storagemigration.k8s.io/v1",
			"since": "1.37",
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
