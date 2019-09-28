<template>
    <div>
      <headerDisplay admin="true"/>
      <div class="container">
          <section class="section">
              <h2 class="title is-2">Hey, {{ login.name }}!</h2>
              <h4 class="title is-4">This is the admin page, you can configure FlatTrack from here</h4>
              <div id="menu-bar-items">
                  <b-menu>
                      <b-menu-list>
                      <b-menu-item v-for="item in pages" v-bind:key="item" :href="item.url" class="menu-bar-item" :label="item.name"></b-menu-item>
                      </b-menu-list>
                  </b-menu>
              </div>
          </section>
        </div>
        <footer class="footer">
        <div class="content has-text-centered">
            <div class="columns is-variable is-1-mobile is-0-tablet is-3-desktop is-8-widescreen is-2-fullhd">
            <div class="column">
                <a href="/#/">Go to general home</a>
            </div>
            </div>
        </div>
      </footer>
    </div>
</template>

<script>
import headerDisplay from '../common/header-display'

export default {
  name: 'Admin home',
  data () {
    return {
      pageLocation: location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : ''),
      pages: [
        {
          name: 'Configure Features',
          url: '#/admin/features'
        },
        {
          name: 'Manage Flatmates',
          url: '#/admin/members'
        },
        {
          name: 'Tasks',
          url: '#/admin/tasks'
        }
      ],
      login: {
        name: 'Person'
      }
    }
  },
  components: {
    headerDisplay
  }
}
</script>

<style scoped>

</style>
