<template>
  <div>
    <div class="container">
      <section class="section">
        <h1 class="title is-1">Apps</h1>
        <p class="subtitle is-4">Manage your flat</p>
        <p></p>
        <br>
        <div v-for="app in apps" v-bind:key="app" >
          <div class="card pointer-cursor-on-hover" @click="goToApp(app.ref)">
            <div class="card-content">
              <div class="media">
                <div class="media-left">
                  <b-icon :icon="app.icon" size="is-medium"></b-icon>
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
          description: 'Create and manage lists to use for your next shop',
          icon: 'format-list-checks',
          ref: 'Shopping list'
        },
        {
          name: 'Flatmates',
          description: 'View your flatmates and their contact details',
          icon: 'account-group',
          ref: 'My Flatmates'
        }
      ]
    }
  },
  methods: {
    goToApp (ref) {
      this.$router.push({ name: ref })
    }
  }
}
</script>

<style scoped>

</style>
