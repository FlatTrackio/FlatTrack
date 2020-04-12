<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
            <li><router-link to="/">Home</router-link></li>
            <li class="is-active"><router-link to="/admin">Admin</router-link></li>
            </ul>
        </nav>
        <h1 class="title is-1">Admin apps</h1>
        <p class="subtitle is-3">Manage FlatTrack</p>
        <p></p>
        <br>
        <!-- TODO add icons to the left sides of cards -->
        <div class="tile is-ancestor pointer-cursor-on-hover" v-for="app in apps" v-bind:key="app">
          <div class="tile is-parent" @click="goToApp(app.ref)">
            <div class="tile is-child box">
              <p class="title is-2">{{ app.name }}</p>
              <p class="subtitle is-4">{{ app.description }}</p>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Admin home',
  data () {
    return {
      apps: [
        {
          name: 'Flatmates',
          description: 'Manage your flatmates',
          ref: '/admin/accounts'
        }
      ]
    }
  },
  methods: {
    goToApp (ref) {
      this.$router.push({ path: ref })
    }
  }
}
</script>

<style scoped>

</style>
