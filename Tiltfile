yaml = helm(
  'k8s-manifests/helm',
  name='flattrack-dev',
  namespace='flattrack-dev',
  set=[
      'service.type=NodePort',
      'prometheus.enabled=false',
      'horizonalPodAutoscaler.enabled=false',
      'postgres.username=flattrack',
      'postgres.password=flattrack',
      'postgres.host=postgres',
      'postgres.database=flattrack'
  ]
  )
k8s_yaml(yaml)
k8s_yaml(kustomize('k8s-manifests/development'))
docker_build('registry.gitlab.com/flattrack/flattrack', '.')
allow_k8s_contexts('in-cluster')
