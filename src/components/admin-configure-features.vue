<template>
    <div>
        <section class="hero is-dark">
            <div class="hero-body">
                <p class="title">
                    FlatTrack (admin)
                </p>
                <p class="subtitle">
                    {{ deploymentName }}
                </p>
            </div>
        </section>
        <div class="container">
            <section class="section">
                <nav class="breadcrumb has-arrow-separator" aria-label="breadcrumbs">
                <ul>
                    <li><a href="/#/admin">Admin Home</a></li>
                    <li class="is-active"><a href="/#/admin/admin-configure-features">Configure Features</a></li>
                </ul>
                </nav>
                <h1 class="title">Configure Features</h1>
                <h2 class="subtitle">Choose the features you'll use</h2>
                <div class="field" v-for="feature in features" v-bind:key="feature">
                    <b-checkbox v-model="featuresEnabled" :native-value="feature.name">{{ feature.name }}</b-checkbox>
                </div>
            </section>
        </div>
    </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'Admin home',
  data () {
    return {
      deploymentName: '',
      features: [
          {
            name: 'Tasks'
          },
          {
            name: 'Shopping List'
          },
          {
            name: 'Noticeboard'
          },
          {
            name: 'Shared Calendar'
          },
          {
            name: 'Recipes'
          },
          {
            name: 'Flatmates'
          },
          {
            name: 'Highfives'
          }
      ],
      featuresEnabled: ['Tasks', 'Shopping List', 'Noticeboard', 'Shared Calendar', 'Recipes', 'Flatmates', 'Highfives']
    }
  },
  created () {
    axios.get(`${location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : '')}/api/settings/deploymentName`)
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