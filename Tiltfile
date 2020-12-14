#!/usr/bin/python

# Tiltfile
#   a fast development flow for FlatTrack

containerRepo = 'registry.gitlab.com/flattrack/flattrack'
containerRepoDocs = 'registry.gitlab.com/flattrack/flattrack/docs'

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
        'ingress.annotations.cert-manager\\.io/cluster-issuer=letsencrypt-prod',
        'ingress.tls[0].secretName=letsencrypt-prod',
        'ingress.tls[0].hosts[0]=' + devHostname,
        'ingress.hosts[0].paths[0]=/',
        'siteURL=http://' + devHostname,
        'realIPHeader=X-Real-Ip'
    ]

namespace='flattrack-dev'
if os.getenv('FLATTRACK_NAMESPACE'):
    namespace = os.get('FLATTRACK_NAMESPACE')

# prepare and use Helm chart
yaml = helm(
  'deployments/flattrack',
  name='flattrack-dev',
  namespace=namespace,
  set=helmSet
  )
k8s_yaml(yaml)

# port-forward FlatTrack and Postgres to ports on the localhost
k8s_resource(workload='flattrack-dev', port_forwards=8080)
k8s_resource(workload='postgres', port_forwards=5432)

# use misc development manifests (postgres, etc...)
k8s_yaml(kustomize('deployments/k8s-manifests/development'))
k8s_yaml(kustomize('deployments/k8s-manifests/docs'))
# if using Kind with Podman
if os.getenv('KIND_EXPERIMENTAL_PROVIDER') == 'podman' and k8s_context() == 'kind-kind':
    custom_build(containerRepo, 'podman build -f build/dev.Dockerfile -t $EXPECTED_REF . && podman save $EXPECTED_REF > /tmp/tilt-containerbuild.tar.gz && kind load image-archive /tmp/tilt-containerbuild.tar.gz', ['.'], disable_push=True, skips_local_docker=True)
    custom_build(containerRepoDocs, 'podman build -f build/docs.Dockerfile -t $EXPECTED_REF . && podman save $EXPECTED_REF > /tmp/tilt-containerbuild.tar.gz && kind load image-archive /tmp/tilt-containerbuild.tar.gz', ['.'], disable_push=True, skips_local_docker=True)
# if using a pair instance
elif os.getenv('SHARINGIO_PAIR_NAME'):
    custom_build(containerRepo, 'docker build -f build/dev.Dockerfile -t $EXPECTED_REF .', ['.'], disable_push=True)
    custom_build(containerRepoDocs, 'docker build -f build/docs.Dockerfile -t $EXPECTED_REF .', ['.'], disable_push=True)
# standard
else:
    docker_build(containerRepo, '.', dockerfile="build/dev.Dockerfile")
    docker_build(containerRepoDocs, '.', dockerfile="build/docs.Dockerfile")

# disallow production clusters
allow_k8s_contexts('in-cluster')
