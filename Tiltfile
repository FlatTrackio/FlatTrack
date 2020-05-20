yaml = helm(
  'k8s-manifests/helm',
  name='flattrack-dev',
  namespace='flattrack-dev',
  set=[
      'service.type=NodePort',
      'prometheus.enabled=false',
      'podDisruptionBudget.enabled=false',
      'horizonalPodAutoscaler.enabled=false'
  ]
  )
k8s_yaml(yaml)
k8s_yaml(kustomize('k8s-manifests/development'))
docker_build('registry.gitlab.com/flattrack/flattrack', '.')
allow_k8s_contexts('in-cluster')
