<template>
    <div>
      <headerDisplay admin="true"/>
      <div class="container">
        <section class="section">
          <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><a href="/">Home</a></li>
              <li class="is-active"><a href="/admin/tasks">Admin</a></li>
            </ul>
          </nav>
          <h2 class="title is-2">Admin</h2>
          <p class="subtitle is-4">This is the admin page, you can configure FlatTrack from here</p>
          <div>
            <div v-for="item in pages" v-bind:key="item">
              <b-button type="is-white" tag="a" :href="item.url" size="is-large" expanded>{{ item.name }}</b-button>
            </div>
          </div>
        </section>
      </div>
    </div>
</template>

<script>
import headerDisplay from '@/frontend/components/header-display'

export default {
  name: 'Admin home',
  data () {
    return {
      pages: [
        {
          name: 'Configure Features',
          url: 'admin/features'
        },
        {
          name: 'Manage Flatmates',
          url: 'admin/members'
        },
        {
          name: 'Manage Tasks',
          url: 'admin/tasks'
        },
        {
          name: 'Assigned tasks and entries',
          url: 'admin/entries'
        }
      ],
      login: {
        names: 'Person'
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
