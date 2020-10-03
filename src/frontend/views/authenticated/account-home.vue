<template>
  <div>
    <div class="container">
      <section class="section">
        <h1 class="title is-1">My Account</h1>
        <p class="subtitle is-4">Manage your account</p>
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
        <div class="card pointer-cursor-on-hover" v-if="typeof CurrentTheme === 'object'" disabled>
          <div class="card-content">
            <div class="media">
              <div class="media-left">
                <b-icon :icon="CurrentTheme.icon" size="is-medium"></b-icon>
              </div>
              <div class="media-content">
                <p class="title is-3">Current theme</p>
                <p class="subtitle is-5">{{ CurrentTheme.name }}</p>
              </div>
            </div>
          </div>
          <div class="content">
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import theme from '@/frontend/common/theme'

export default {
  name: 'Account home',
  data () {
    return {
      apps: [
        {
          name: 'Profile',
          description: 'Manage your general information',
          icon: 'account-circle',
          ref: '/account/profile'
        },
        {
          name: 'Security',
          description: 'Manage your account security',
          icon: 'lock-question',
          ref: '/account/security'
        },
        {
          name: 'Settings',
          description: 'Manage settings for this device',
          icon: 'account-cog',
          ref: '/account/settings'
        }
      ]
    }
  },
  methods: {
    goToApp (ref) {
      this.$router.push({ path: ref })
    }
  },
  computed: {
    CurrentTheme () {
      return theme.GetTheme()
    }
  }
}
</script>

<style scoped>

</style>
