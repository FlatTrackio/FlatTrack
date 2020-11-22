#!/usr/bin/python

# Tiltfile
#   a fast development flow for FlatTrack

# default values
helmSet=[
      'service.type=NodePort',
      'prometheus.enabled=false',
      'postgres.username=flattrack',
      'postgres.password=flattrack',
      'postgres.host=postgres',
      'postgres.database=flattrack'
]

# when using a pair instance
if os.getenv('SHARINGIO_PAIR_NAME'):
    devHostname=os.getenv("FLATTRACK_DEV_HOSTNAME")
    helmSet = helmSet + [
        'ingress.enabled=true',
        'ingress.hosts[0].host=' + devHostname,
        'ingress.hosts[0].paths[0]=/',
        'siteURL=http://' + devHostname,
        'realIPHeader=X-Real-Ip'
    ]

# prepare and use Helm chart
yaml = helm(
  'deployments/flattrack',
  name='flattrack-dev',
  namespace='flattrack-dev',
  set=helmSet
  )
k8s_yaml(yaml)

# use misc development manifests (postgres, etc...)
k8s_yaml(kustomize('deployments/k8s-manifests/development'))
# if using Kind with Podman
if os.getenv('KIND_EXPERIMENTAL_PROVIDER') == 'podman' and k8s_context() == 'kind-kind':
    custom_build('registry.gitlab.com/flattrack/flattrack', 'podman build -f build/dev.Dockerfile -t $EXPECTED_REF . && podman save $EXPECTED_REF > /tmp/tilt-containerbuild.tar.gz && kind load image-archive /tmp/tilt-containerbuild.tar.gz', ['.'], disable_push=True, skips_local_docker=True)
# if using a pair instance
elif os.getenv('SHARINGIO_PAIR_NAME'):
    custom_build('registry.gitlab.com/flattrack/flattrack', 'docker build -f build/dev.Dockerfile -t $EXPECTED_REF .', ['.'], disable_push=True)
# standard
else:
    docker_build('registry.gitlab.com/flattrack/flattrack', '.', dockerfile="build/dev.Dockerfile")

# disallow production clusters
allow_k8s_contexts('in-cluster')

# port-forward FlatTrack and Postgres to ports on the localhost
local_resource(name="flattrack-port-forward", serve_cmd='kubectl -n flattrack-dev port-forward svc/flattrack-dev 8080:8080')
local_resource(name="postgres-port-forward", serve_cmd='kubectl -n flattrack-dev port-forward svc/postgres 5432:5432')
