# Deployment with Kubernetes

## Configuration

Chart values are located in [../deployments/flattrack/values.yaml](https://gitlab.com/flattrack/flattrack/-/blob/master/deployments/flattrack/values.yaml).

<table border="2" cellspacing="0" cellpadding="6" rules="groups" frame="hsides">


<colgroup>
<col  class="org-left" />

<col  class="org-left" />

<col  class="org-left" />
</colgroup>
<thead>
<tr>
<th scope="col" class="org-left">Parameter</th>
<th scope="col" class="org-left">Description</th>
<th scope="col" class="org-left">Default</th>
</tr>
</thead>

<tbody>
<tr>
<td class="org-left">replicaCount</td>
<td class="org-left">FlatTrack replica Pods</td>
<td class="org-left">1</td>
</tr>


<tr>
<td class="org-left">image.repository</td>
<td class="org-left">The repo where the image lives</td>
<td class="org-left">registry.gitlab.com/flattrack/flattrack</td>
</tr>


<tr>
<td class="org-left">image.tag</td>
<td class="org-left">Specifies a tag of from the image to use</td>
<td class="org-left"><code>nil</code></td>
</tr>


<tr>
<td class="org-left">image.pullPolicy</td>
<td class="org-left">FlatTrack container pull policy</td>
<td class="org-left">IfNotPresent</td>
</tr>


<tr>
<td class="org-left">imagePullSecrets</td>
<td class="org-left">References for the registry secrets to pull FlatTrack from</td>
<td class="org-left"><code>[]</code></td>
</tr>


<tr>
<td class="org-left">nameOverride</td>
<td class="org-left">Expand the name of the chart</td>
<td class="org-left"><code>""</code></td>
</tr>


<tr>
<td class="org-left">fullNameOverride</td>
<td class="org-left">Create a FQDN for the app name</td>
<td class="org-left"><code>""</code></td>
</tr>


<tr>
<td class="org-left">podSecurityContext.readOnlyRootFilesystem</td>
<td class="org-left">Set the rootfs as read-only</td>
<td class="org-left">true</td>
</tr>


<tr>
<td class="org-left">podSecurityContext.runAsUser</td>
<td class="org-left">The user to run as</td>
<td class="org-left">1000</td>
</tr>


<tr>
<td class="org-left">podSecurityContext.runAsGroup</td>
<td class="org-left">The group to run as</td>
<td class="org-left">1000</td>
</tr>


<tr>
<td class="org-left">podSecurityContext.allowPrivilegeEscalation</td>
<td class="org-left">If the process in the container can become root</td>
<td class="org-left">true</td>
</tr>


<tr>
<td class="org-left">service.type</td>
<td class="org-left">In way which the app is exposed</td>
<td class="org-left">ClusterIP</td>
</tr>


<tr>
<td class="org-left">service.port</td>
<td class="org-left">The port to run the app on</td>
<td class="org-left">8080</td>
</tr>


<tr>
<td class="org-left">horizonalPodAutoscaler.enabled</td>
<td class="org-left">If the Pods should autoscale</td>
<td class="org-left">true</td>
</tr>


<tr>
<td class="org-left">horizonalPodAutoscaler.minReplicas</td>
<td class="org-left">Minimum amount of Pods</td>
<td class="org-left">2</td>
</tr>


<tr>
<td class="org-left">horizonalPodAutoscaler.maxReplicas</td>
<td class="org-left">Maximum amount of Pods</td>
<td class="org-left">5</td>
</tr>


<tr>
<td class="org-left">horizonalPodAutoscaler.targetCPUUtilizationPercentage</td>
<td class="org-left">How much resource should be utilized before scaling</td>
<td class="org-left">5</td>
</tr>


<tr>
<td class="org-left">podDisruptionBudget.enabled</td>
<td class="org-left">define the max number of disruptions allowed</td>
<td class="org-left">enable</td>
</tr>


<tr>
<td class="org-left">podDisruptionBudget.minAvailable</td>
<td class="org-left">minimum number of Pods that should always be available</td>
<td class="org-left">2</td>
</tr>


<tr>
<td class="org-left">podDisruptionBudget.maxUnavailable</td>
<td class="org-left">max amount of Pods that are allowed to be unavailable</td>
<td class="org-left">&#xa0;</td>
</tr>


<tr>
<td class="org-left">prometheus.enabled</td>
<td class="org-left">enable Prometheus annotations</td>
<td class="org-left">true</td>
</tr>


<tr>
<td class="org-left">networkPolicy.enabled</td>
<td class="org-left">enable networkPolicy rules</td>
<td class="org-left">false</td>
</tr>


<tr>
<td class="org-left">networkPolicy.ingress</td>
<td class="org-left">configure Ingress policy rules</td>
<td class="org-left">{}</td>
</tr>


<tr>
<td class="org-left">networkPolicy.egress</td>
<td class="org-left">configure Egress policy rules</td>
<td class="org-left">{}</td>
</tr>


<tr>
<td class="org-left">extraEnvVars</td>
<td class="org-left">declare extra environment variables</td>
<td class="org-left">&#xa0;</td>
</tr>


<tr>
<td class="org-left">postgres.enabled</td>
<td class="org-left">use the helm chart to set the database credentials</td>
<td class="org-left">true</td>
</tr>


<tr>
<td class="org-left">postgres.username</td>
<td class="org-left">the username for an existing Postgres databse</td>
<td class="org-left">flattrack</td>
</tr>


<tr>
<td class="org-left">postgres.password</td>
<td class="org-left">the password for an existing Postgres databse</td>
<td class="org-left">flattrack</td>
</tr>


<tr>
<td class="org-left">postgres.host</td>
<td class="org-left">the host for an existing Postgres databse</td>
<td class="org-left">flattrack</td>
</tr>


<tr>
<td class="org-left">postgres.database</td>
<td class="org-left">the database for an existing Postgres databse</td>
<td class="org-left">flattrack</td>
</tr>


<tr>
<td class="org-left">postgres.sslmode</td>
<td class="org-left">the sslmode for an existing Postgres databse</td>
<td class="org-left">disable</td>
</tr>


<tr>
<td class="org-left">smtp.enabled</td>
<td class="org-left">enable SMTP for messages and notifications</td>
<td class="org-left">false</td>
</tr>


<tr>
<td class="org-left">smtp.username</td>
<td class="org-left">SMTP username</td>
<td class="org-left"><code>""</code></td>
</tr>


<tr>
<td class="org-left">smtp.password</td>
<td class="org-left">SMTP password</td>
<td class="org-left"><code>""</code></td>
</tr>


<tr>
<td class="org-left">smtp.host</td>
<td class="org-left">SMTP host</td>
<td class="org-left"><code>""</code></td>
</tr>


<tr>
<td class="org-left">smtp.port</td>
<td class="org-left">SMTP port</td>
<td class="org-left"><code>""</code></td>
</tr>


<tr>
<td class="org-left">siteURL</td>
<td class="org-left">the URL of the site</td>
<td class="org-left"><code>""</code></td>
</tr>


<tr>
<td class="org-left">realIPHeader</td>
<td class="org-left">a HTTP header to use for forwarding the real IP address</td>
<td class="org-left"><code>""</code></td>
</tr>


<tr>
<td class="org-left">embeddedHTML</td>
<td class="org-left">inject custom HTML into index.html</td>
<td class="org-left"><code>""</code></td>
</tr>


<tr>
<td class="org-left">labels</td>
<td class="org-left">declare labels for all resources</td>
<td class="org-left"><code>{}</code></td>
</tr>


<tr>
<td class="org-left">annotations</td>
<td class="org-left">declare annotations for all resources</td>
<td class="org-left"><code>{}</code></td>
</tr>


<tr>
<td class="org-left">ingress.enabled</td>
<td class="org-left">create an ingress manifests</td>
<td class="org-left">false</td>
</tr>


<tr>
<td class="org-left">ingress.annotations</td>
<td class="org-left">set annotations for the ingress manifest</td>
<td class="org-left"><code>{}</code></td>
</tr>


<tr>
<td class="org-left">ingress.hosts</td>
<td class="org-left">the hosts which the ingress endpoint should be accessed from</td>
<td class="org-left">&#xa0;</td>
</tr>


<tr>
<td class="org-left">ingress.tls</td>
<td class="org-left">references to TLS secrets</td>
<td class="org-left"><code>[]</code></td>
</tr>


<tr>
<td class="org-left">resources.limits.cpu</td>
<td class="org-left">max amount of CPU</td>
<td class="org-left">250m</td>
</tr>


<tr>
<td class="org-left">resources.limits.memory</td>
<td class="org-left">max amount of memory</td>
<td class="org-left">40Mi</td>
</tr>


<tr>
<td class="org-left">resources.resources.cpu</td>
<td class="org-left">requested amount of CPU</td>
<td class="org-left">250m</td>
</tr>


<tr>
<td class="org-left">resources.limits.memory</td>
<td class="org-left">max amount of memory</td>
<td class="org-left">40Mi</td>
</tr>


<tr>
<td class="org-left">nodeSelector</td>
<td class="org-left">declare the node labels for Pod scheduling</td>
<td class="org-left"><code>{}</code></td>
</tr>


<tr>
<td class="org-left">tolerations</td>
<td class="org-left">declare the toleration labels for Pod scheduling</td>
<td class="org-left"><code>[]</code></td>
</tr>


<tr>
<td class="org-left">affinity</td>
<td class="org-left">declare the affinity settings for the Pod scheduling</td>
<td class="org-left"><code>{}</code></td>
</tr>
</tbody>
</table>

## Installation

Create the namespace:

    kubectl create namespace flattrack-myflat

Install the release with Helm:

    helm install flattrack-myflat --namespace flattrack-myflat deployments/flattrack

# Extra notes

To configure FlatTrack, please refer to the [configuration guide](./configuration.md).

