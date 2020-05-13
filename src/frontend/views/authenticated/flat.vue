<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><router-link to="/">Home</router-link></li>
              <li class="is-active"><router-link to="/flat">My flat</router-link></li>
            </ul>
        </nav>
        <h1 v-if="hasInitialLoaded || flatName !== ''" class="title is-1">{{ flatName }}</h1>
        <b-skeleton v-else size="is-medium" width="35%" :animated="true"></b-skeleton>
        <p class="subtitle is-3">About your flat</p>
        <b-message type="is-warning">
          This section for describing such things as, but not limited to:
          <br/>
          <ul style="list-style-type: disc;">
            <li>how the flat life is</li>
            <li>rules</li>
            <li>regulations</li>
            <li>culture</li>
          </ul>
        </b-message>
      </section>
    </div>
  </div>
</template>

<script>
import flatInfo from '@/frontend/requests/authenticated/flatInfo'
import common from '@/frontend/common/common'

export default {
  name: 'flat',
  data () {
    return {
      flatName: '',
      hasInitialLoaded: false
    }
  },
  methods: {
    GetFlatName () {
      flatInfo.GetFlatName().then(resp => {
        if (this.flatName !== resp.data.spec) {
          this.hasInitialLoaded = true
          this.flatName = resp.data.spec
          common.WriteFlatnameToCache(resp.data.spec)
        }
      })
    }
  },
  async beforeMount () {
    this.flatName = common.GetFlatnameFromCache() || this.flatName
    this.GetFlatName()
  }
}
</script>
