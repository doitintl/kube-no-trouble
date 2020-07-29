package deprecated120

main[return] {
  resource := input[_]
  old_api := deprecated_resource(resource)
  return := {"Name": resource.metadata.name, 
             # Namespace does not have to be defined in case of local manifests
             "Namespace": get_default(resource.metadata, "namespace", "<undefined>"),
	         "Kind": resource.kind,
	         "ApiVersion": old_api,
	         "RuleSet": "1.20 Deprecated APIs"}
}

deprecated_resource(r) = old_api {
  old_api := deprecated_api(r.kind, r.apiVersion)
}

deprecated_api(kind, api_version) = old_api {
  deprecated_apis = { # -> networking.k8s.io/v1beta1
                      "Ingress":        ["extensions/v1beta1"],
                    }
  deprecated_apis[kind][_] == api_version
  old_api := api_version
}

get_default(val, key, _) = val[key]
get_default(val, key, fallback) = fallback { not val[key] }