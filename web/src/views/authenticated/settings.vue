<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><router-link to="/account">My account</router-link></li>
              <li class="is-active"><router-link to="/account/settings">Settings</router-link></li>
            </ul>
        </nav>
        <h1 class="title is-1">Settings</h1>
        <p class="subtitle is-3">Manage settings for this device</p>
        <b-loading :is-full-page="false" :active.sync="pageLoading" :can-cancel="false"></b-loading>

        <b-field
          label="Miscellaneous">
          <b-checkbox
            size="is-medium"
            v-model="enableAnimations">
            Enable Animations
          </b-checkbox>
        </b-field>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/common/common'

export default {
  name: 'profile',
  data () {
    return {
      enableAnimations: common.GetEnableAnimations() === 'true'
    }
  },
  watch: {
    enableAnimations () {
      common.WriteEnableAnimations(this.enableAnimations)
      if (this.enableAnimations === 'true') {
        common.Hooray()
      }
    }
  }
}
</script>
