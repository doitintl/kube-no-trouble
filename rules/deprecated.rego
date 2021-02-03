package deprecated

main[return] {
	resource := input[_]
	api := deprecated_resource(resource)
	return := {
		"Name": resource.metadata.name,
		# Namespace does not have to be defined in case of local manifests
		"Namespace": get_default(resource.metadata, "namespace", "<undefined>"),
		"Kind": resource.kind,
		"ApiVersion": api.old,
		"ReplaceWith": api.new,
		"RuleSet": "Deprecated APIs without removal annoucement",
		"Since": api.since,
	}
}

deprecated_resource(r) = api {
	api := deprecated_api(r.kind, r.apiVersion)
}

deprecated_api(kind, api_version) = api {
	deprecated_apis = {
		"BootstrapTokenString": {
			"old": ["kubeadm.k8s.io/v1beta1"],
			"new": "kubeadm.k8s.io/v1beta2",
			"since": "1.17",
		},
		"CustomResourceDefinition": {
			"old": ["apiextensions.k8s.io/v1beta1"],
			"new": "apiextensions.k8s.io/v1",
			"since": "1.19",
		},
		"CustomResourceDefinitionList": {
			"old": ["apiextensions.k8s.io/v1beta1"],
			"new": "apiextensions.k8s.io/v1",
			"since": "1.19",
		},
		"APIService": {
			"old": ["apiregistration.k8s.io/v1beta1"],
			"new": "apiregistration.k8s.io/v1",
			"since": "1.19",
		},
		"APIServiceList": {
			"old": ["apiregistration.k8s.io/v1beta1"],
			"new": "apiregistration.k8s.io/v1",
			"since": "1.19",
		},
		"HorizontalPodAutoscaler": {
			"old": ["autoscaling/v2beta1"],
			"new": "autoscaling/v2beta2",
			"since": "1.19",
		},
		"HorizontalPodAutoscalerList": {
			"old": ["autoscaling/v2beta1"],
			"new": "autoscaling/v2beta2",
			"since": "1.19",
		},
		"StorageClass": {
			"old": ["storage.k8s.io/v1beta1"],
			"new": "storage.k8s.io/v1",
			"since": "1.19",
		},
		"StorageClassList": {
			"old": ["storage.k8s.io/v1beta1"],
			"new": "storage.k8s.io/v1",
			"since": "1.19",
		},
		"VolumeAttachment": {
			"old": ["storage.k8s.io/v1beta1"],
			"new": "storage.k8s.io/v1",
			"since": "1.19",
		},
		"VolumeAttachmentList": {
			"old": ["storage.k8s.io/v1beta1"],
			"new": "storage.k8s.io/v1",
			"since": "1.19",
		},
		"CSIDriver": {
			"old": ["storage.k8s.io/v1beta1"],
			"new": "storage.k8s.io/v1",
			"since": "1.19",
		},
		"CSIDriverList": {
			"old": ["storage.k8s.io/v1beta1"],
			"new": "storage.k8s.io/v1",
			"since": "1.19",
		},
		"CSINode": {
			"old": ["storage.k8s.io/v1beta1"],
			"new": "storage.k8s.io/v1",
			"since": "1.19",
		},
		"CSINodeList": {
			"old": ["storage.k8s.io/v1beta1"],
			"new": "storage.k8s.io/v1",
			"since": "1.19",
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
