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
      'postgres.database=flattrack',
      'minio.accessKey=flattrack',
      'minio.secretKey=flattrack',
      'minio.host=minio:9000',
      'minio.SSL.enabled=false',
      'minio.bucket=bucket'
  ]
  )
k8s_yaml(yaml)
k8s_yaml(kustomize('k8s-manifests/development'))
load('ext://helm_remote', 'helm_remote')
helm_remote(
  'minio',
  repo_url="https://helm.min.io/",
  namespace='flattrack-dev',
  version="7.2.1",
  set=[
      "accessKey=flattrack",
      "secretKey=flattrack",
      "defaultBucket.enabled=true",
      "resources.requests.memory=1Gi"
  ]
  )
docker_build('registry.gitlab.com/flattrack/flattrack', '.')
allow_k8s_contexts('in-cluster')
