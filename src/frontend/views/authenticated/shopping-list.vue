<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><a href="/">Home</a></li>
              <li><a href="/apps">Apps</a></li>
              <li class="is-active"><a href="/apps/shopping-list">Shopping list</a></li>
            </ul>
        </nav>
        <h1 class="title">Shopping list</h1>
        <h2 class="subtitle">Manage your weekly shop</h2>
        <b-tabs position="is-centered" class="block">
          <b-tab-item label="In Progress"></b-tab-item>
          <b-tab-item label="Completed" :disabled="lists.length === 0"></b-tab-item>
        </b-tabs>
        <div id="lists" v-if="!lists.length">
          <section>
            <div class="card">
              <div class="card-content">
                <div class="media">
                  <div class="media-content">
                    <p class="title is-4">A bit empty here</p>
                  </div>
                </div>
                <div class="content">
                  Press the (+) button to add a new shopping list
                  <br />
                </div>
              </div>
            </div>
          </section>
        </div>
        <br/><br/><br/>
        <addButton @click="addList"/>
    </div>
  </div>
</template>

<script>
import common from '@/frontend/common/common'

export default {
  name: 'Shopping List',
  data () {
    return {
      lists: []
    }
  },
  components: {
    addButton: () => import('@/frontend/components/common/floatingAddButton')
  },
  methods: {
    addList: () => {
      console.log('Adding a list')
    }
  }
}
</script>

<style scoped>

</style>
