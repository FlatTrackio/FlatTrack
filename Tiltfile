k8s_yaml(kustomize('k8s-manifests/development'))
docker_build('registry.gitlab.com/flattrack/flattrack', '.')
allow_k8s_contexts('in-cluster')
