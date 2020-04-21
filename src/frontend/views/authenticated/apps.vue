<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
            <li><router-link to="/">Home</router-link></li>
            <li class="is-active"><router-link to="/apps">Apps</router-link></li>
            </ul>
        </nav>
        <h1 class="title is-1">Apps</h1>
        <p class="subtitle is-3">Manage your flat</p>
        <p></p>
        <br>
        <!-- TODO add icons to the left sides of cards -->
        <div v-for="app in apps" v-bind:key="app" >
          <div class="card pointer-cursor-on-hover" @click="goToApp(app.ref)">
            <div class="card-content">
              <div class="media">
                <div class="media-left">
                </div>
                <div class="media-content">
                  <p class="title is-3">{{ app.name }}</p>
                  <p class="subtitle is-5">{{ app.description }}</p>
                </div>
                <div class="media-right">
                  <b-icon icon="chevron-right" size="is-medium" type="is-midgray"></b-icon>
                </div>
              </div>
            </div>
            <div class="content">
            </div>
          </div>
          <br />
        </div>
      </section>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Apps',
  data () {
    return {
      apps: [
        {
          name: 'Shopping list',
          description: 'Manage a flat-wide shopping list for all to contribute to',
          ref: '/apps/shopping-list'
        },
        {
          name: 'Flatmates',
          description: 'View your flatmates and their contact details',
          ref: '/apps/flatmates'
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
