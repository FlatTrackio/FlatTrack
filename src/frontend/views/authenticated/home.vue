<template>
  <div>
    <div class="container">
      <section class="section">
        <h1 class="title is-1">Home</h1>
        <b-message type="is-warning">
          This page will be a dashboard of recently added or updated things.
          <br/>
          As it's not available yet, please navigate to other pages which do things other than a message like this.
        </b-message>
      </section>
    </div>
  </div>
</template>

<script>
export default {
  name: 'home',
  data () {
    return {
    }
  },
  methods: {
  },
  async beforeMount () {
  }
}
</script>

<style scoped>
</style>
