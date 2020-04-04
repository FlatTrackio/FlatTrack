<template>
  <div id="topbar">
    <b-navbar type="is-info" shadow="true" transparent="true">
      <template slot="brand">
        <b-navbar-item tag="router-link" :to="{ path: '/' }">
          <h1 class="title is-5" style="color: #fff;">FlatTrack ({{ flatName }})</h1>
        </b-navbar-item>
      </template>
      <template slot="start">
        <b-navbar-item tag="router-link" :to="{ path: '/' }">
          Home
        </b-navbar-item>
        <b-navbar-item href="https://flattrack.io/help" tag="a" target="_blank">
          Help
        </b-navbar-item>
        <b-navbar-item @click="signOut">
          Sign out
        </b-navbar-item>
      </template>
    </b-navbar>
  </div>
</template>

<script>
import common from '@/frontend/common/common'
import flatInfo from '@/frontend/requests/authenticated/flatInfo'
import { DialogProgrammatic as Dialog, LoadingProgrammatic as Loading } from 'buefy'

export default {
  name: 'topbar',
  data () {
    return {
      flatName: 'My Flat'
    }
  },
  methods: {
    signOut () {
      common.SignoutDialog()
    },
    GetFlatInfo () {
      flatInfo.GetFlatName().then(resp => {
        this.flatName = resp.data.spec
      })
    }
  },
  async created () {
    this.GetFlatInfo()
  }
}
</script>

<style>
</style>
