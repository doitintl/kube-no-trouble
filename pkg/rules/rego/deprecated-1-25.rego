package deprecated125

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
		"RuleSet": "Deprecated APIs removed in 1.25",
		"Since": api.since,
	}
}

deprecated_resource(r) = api {
	api := deprecated_api(r.kind, r.apiVersion)
}

deprecated_api(kind, api_version) = api {
	deprecated_apis = {
		"RuntimeClass": {
			"old": ["node.k8s.io/v1beta1"],
			"new": "node.k8s.io/v1",
			"since": "1.20",
		},
		"PodDisruptionBudget": {
			"old": ["policy/v1beta1"],
			"new": "policy/v1",
			"since": "1.21",
		},
		"PodSecurityPolicy": {
			"old": ["policy/v1beta1"],
			"new": "<removed>",
			"since": "1.21",
		},
		"EndpointSlice": {
			"old": ["discovery.k8s.io/v1beta1"],
			"new": "discovery.k8s.io/v1",
			"since": "1.21",
		},
		"CronJob": {
			"old": ["batch/v1beta1"],
			"new": "batch/v1",
			"since": "1.21",
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
