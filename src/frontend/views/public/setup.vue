<template>
  <div>
    <headerDisplay genericMessage="true"/>
    <div class="section">
      <section class="section">
        <b-steps v-model="step" :has-navigation=false :animated=true>
            <b-step-item label="Set up Database" icon="account-key"></b-step-item>
            <b-step-item label="Add your flatmates" icon="account-plus"></b-step-item>
            <b-step-item label="Add some tasks" icon="account-plus"></b-step-item>
        </b-steps>
        <h1 class="title is-1">Getting started</h1>
        <p class="subtitle is-3">Setup your FlatTrack instance</p>
      </section>
    </div>
  </div>
</template>

<script>
import headerDisplay from '@/frontend/components/header-display'

export default {
  name: 'setup',
  data () {
    return {
      step: 0
    }
  },
  components: {
    headerDisplay
  }
}
</script>

<style scoped>

</style>
