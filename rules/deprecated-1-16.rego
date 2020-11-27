package deprecated116

main[return] {
	resource := input[_]
	old_api := deprecated_resource(resource)
	return := {
		"Name": resource.metadata.name,
		# Namespace does not have to be defined in case of local manifests
		"Namespace": get_default(resource.metadata, "namespace", "<undefined>"),
		"Kind": resource.kind,
		"ApiVersion": old_api,
		"RuleSet": "Deprecated APIs removed in 1.16",
	}
}

deprecated_resource(r) = old_api {
	old_api := deprecated_api(r.kind, r.apiVersion)
}

deprecated_api(kind, api_version) = old_api {
	deprecated_apis = {
		"Deployment": ["extensions/v1beta1", "apps/v1beta1", "apps/v1beta2"], # -> apps/v1
		# -> networking.k8s.io/v1
		"NetworkPolicy": ["extensions/v1beta1"],
		# -> policy/v1beta1
		"PodSecurityPolicy": ["extensions/v1beta1"],
		# -> apps/v1
		"DaemonSet": ["extensions/v1beta1", "apps/v1beta2"],
		# -> apps/v1
		"StatefulSet": ["apps/v1beta1", "apps/v1beta2"],
		# -> apps/v1
		"ReplicaSet": ["extensions/v1beta1", "apps/v1beta1", "apps/v1beta2"],
	}

	deprecated_apis[kind][_] == api_version
	old_api := api_version
}

get_default(val, key, _) = val[key]

get_default(val, key, fallback) = fallback {
	not val[key]
}
