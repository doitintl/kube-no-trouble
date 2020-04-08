package deprecated116

main[return] {
  resource := input[_]
  old_api := deprecated_resource(resource)
  return := {"Name": resource.metadata.name, 
             "Namespace": resource.metadata.namespace,
	     "Kind": resource.kind,
	     "ApiVersion": old_api,
	     "RuleSet": "1.16 Deprecated APIs"}
}

deprecated_resource(r) = old_api {
  last_applied := json.unmarshal(r.metadata.annotations["kubectl.kubernetes.io/last-applied-configuration"])
  old_api := deprecated_api(r.kind, last_applied.apiVersion)
}

deprecated_resource(r) = old_api {
  old_api := deprecated_api(r.kind, r.apiVersion)
}

deprecated_api(kind, api_version) = old_api {
  deprecated_apis = { # -> apps/v1
                      "Deployment":        ["extensions/v1beta1", "apps/v1beta1", "apps/v1beta2"],
                      # -> networking.k8s.io/v1
                      "NetworkPolicy":     ["extensions/v1beta1"],
                      # -> policy/v1beta1
                      "PodSecurityPolicy": ["extensions/v1beta1"],
                      # -> apps/v1
                      "DaemonSet":         ["extensions/v1beta1", "apps/v1beta2"],
                      # -> apps/v1
                      "StatefulSet":       ["apps/v1beta1", "apps/v1beta2"],
                      # -> apps/v1
                      "ReplicaSet":        ["extensions/v1beta1", "apps/v1beta1", "apps/v1beta2"]
                    }
  deprecated_apis[kind][_] == api_version
  old_api := api_version
}
