<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><router-link to="/admin">Admin</router-link></li>
              <li class="is-active"><router-link to="/admin/settings">Settings</router-link></li>
            </ul>
        </nav>
        <h1 class="title is-1">Settings</h1>
        <p class="subtitle is-4">General FlatTrack settings</p>
        <b-loading :is-full-page="false" :active.sync="pageLoading" :can-cancel="false"></b-loading>
        <div>
          <label class="label">Flat name</label>
          <b-field>
            <b-input
              type="text"
              v-model="flatName"
              maxlength="20"
              placeholder="Enter your flat's name"
              icon="textbox"
              size="is-medium"
              expanded
              required>
            </b-input>
            <p class="control">
              <b-button
                type="is-primary"
                size="is-medium"
                icon-left="check"
                @keyup.enter.native="PostFlatName"
                @click="PostFlatName">
              </b-button>
            </p>
        </b-field>
        </div>
        <br/>
      </section>
    </div>
  </div>
</template>

<script>
import flatInfo from '@/frontend/requests/authenticated/flatInfo'
import settings from '@/frontend/requests/admin/settings'
import common from '@/frontend/common/common'

export default {
  name: 'Flatmates',
  data () {
    return {
      pageLoading: true,
      flatName: ''
    }
  },
  async beforeMount () {
    flatInfo.GetFlatName().then(resp => {
      this.flatName = resp.data.spec
      this.pageLoading = false
    })
  },
  methods: {
    goToRef (ref) {
      this.$router.push({ path: ref })
    },
    PostFlatName () {
      if (this.flatName === '') {
        common.DisplayFailureToast('Error: Flat name must not be empty')
        return
      }
      settings.PostFlatName(this.flatName).then(resp => {
        common.DisplaySuccessToast('Set flat name')
      }).catch(err => {
        common.DisplayFailureToast('Failed set the flat name' + '<br/>' + (err.response.data.metadata.response || err))
      })
    },
    TimestampToCalendar (timestamp) {
      return common.TimestampToCalendar(timestamp)
    }
  }
}
</script>
