<template>
    <div>
        <section class="hero is-info">
            <div class="hero-body">
                <p class="title">
                    FlatTrack
                </p>
                <p class="subtitle">
                    {{ deploymentName }}
                </p>
            </div>
        </section>
        <div class="container">
            <section class="section is-large">
                <h2 class="title is-2">Unknown page!</h2>
                <h4 class="title is-4">Hmmm maybe something went wrong :/</h4>
            </section>
        </div>
        <footer class="footer">
            <div class="content has-text-centered">
                <div class="columns is-variable is-1-mobile is-0-tablet is-3-desktop is-8-widescreen is-2-fullhd">
                    <div class="column">
                        <a href="/#/">Go home</a>
                    </div>
                </div>
            </div>
        </footer>
    </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'Unknown page',
  deploymentName: 'Keep track of your flat',
  pageLocation: location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : ''),
  data () {
    return {
    }
  },
  created () {
    axios.get(`/api/settings/deploymentName`)
      .then(response => {
        this.deploymentName = response.data.value
      })
      .catch(err => {
        this.pageErrors.push(err)
      })
  }
}
</script>

<style scoped>

</style>
