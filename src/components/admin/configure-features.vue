<template>
    <div>
      <headerDisplay admin="true"/>
      <div class="container">
        <section class="section">
          <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
              <ul>
                  <li><a href="/admin">Admin Home</a></li>
                  <li class="is-active"><a href="/admin/admin-configure-features">Configure Features</a></li>
              </ul>
          </nav>
          <h1 class="title is-2">Configure Features</h1>
          <h2 class="subtitle is-4">Choose the features you'll use</h2>
          <div v-for="feature in features" v-bind:key="feature">
            <div class="field">
                <b-checkbox v-model="featuresEnabled" :native-value="feature.name" size="is-medium">
                  {{ feature.name }}
                </b-checkbox>
                <b-tooltip :label="feature.description" position="is-right">
                [?]
                </b-tooltip>
            </div>
          </div>
        </section>
      </div>
    </div>
</template>

<script>
import headerDisplay from '../common/header-display'

export default {
  name: 'Admin home',
  data () {
    return {
      features: [
        {
          name: 'Tasks',
          description: 'Self management of tasks'
        },
        {
          name: 'Shopping List',
          description: 'A list of items for your flat\'s weekly shop'
        },
        {
          name: 'Noticeboard',
          description: 'Post to your flatmates'
        },
        {
          name: 'Shared Calendar',
          description: 'Keep track of events'
        },
        {
          name: 'Recipes',
          description: 'Create and share recipes'
        },
        {
          name: 'Flatmates',
          description: 'Basic information'
        },
        {
          name: 'Highfives',
          description: 'Celebrate your flatmates'
        }
      ],
      featuresEnabled: ['Tasks', 'Shopping List', 'Noticeboard', 'Shared Calendar', 'Recipes', 'Flatmates', 'Highfives']
    }
  },
  components: {
    headerDisplay
  }
}
</script>

<style scoped>

</style>
