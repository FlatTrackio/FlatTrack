<template>
  <div>
    <div class="container">
      <section class="section">
        <h1 class="title is-1">Admin apps</h1>
        <p class="subtitle is-4">Manage FlatTrack</p>
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
  name: 'Admin home',
  data () {
    return {
      apps: [
        {
          name: 'Accounts',
          description: 'Manage your flatmates accounts',
          icon: 'account-group',
          ref: 'Admin accounts'
        },
        {
          name: 'Settings',
          description: 'General FlatTrack settings',
          icon: 'settings',
          ref: 'Admin settings'
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
