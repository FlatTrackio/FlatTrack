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
      <nav class="breadcrumb has-arrow-separator" aria-label="breadcrumbs">
          <ul>
          <li><a href="/#/">Home</a></li>
          <li class="is-active"><a href="/#/shoppinglist">Shopping List</a></li>
          </ul>
      </nav>
      <h1 class="title">Shopping List</h1>
      <p>Page under construction, coming soon</p>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Shopping List',
  deploymentName: 'Keep track of your flat',
  data () {
    return {
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
