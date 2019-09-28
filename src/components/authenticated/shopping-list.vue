<template>
  <div>
    <headerDisplay/>
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
import headerDisplay from '../common/header-display'

export default {
  name: 'Shopping List',
  data () {
    return {
    }
  },
  components: {
    headerDisplay
  }
}
</script>

<style scoped>

</style>
