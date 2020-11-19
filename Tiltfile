yaml = helm(
  'deployments/flattrack',
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
k8s_yaml(kustomize('deployments/k8s-manifests/development'))
if os.getenv('KIND_EXPERIMENTAL_PROVIDER') == 'podman' and k8s_context() == 'kind-kind':
    custom_build('registry.gitlab.com/flattrack/flattrack', 'podman build -f build/dev.Dockerfile -t $EXPECTED_REF . && podman save $EXPECTED_REF > /tmp/tilt-containerbuild.tar.gz && kind load image-archive /tmp/tilt-containerbuild.tar.gz', ['.'], disable_push=True, skips_local_docker=True)
elif not os.getenv('SHARINGIO_PAIR_NAME'):
    docker_build('registry.gitlab.com/flattrack/flattrack', '.', dockerfile="build/dev.Dockerfile")
else:
    custom_build('registry.gitlab.com/flattrack/flattrack', 'docker build -f build/dev.Dockerfile -t $EXPECTED_REF .', ['.'], disable_push=True)
allow_k8s_contexts('in-cluster')
local_resource(name="flattrack-port-forward", serve_cmd='kubectl -n flattrack-dev port-forward svc/flattrack-dev 8080:8080')
local_resource(name="postgres-port-forward", serve_cmd='kubectl -n flattrack-dev port-forward svc/postgres 5432:5432')
