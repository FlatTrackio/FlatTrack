<template>
  <div id="topbar">
    <b-navbar fixed-top type="is-info" class="gradient-blue" shadow="true" transparent="false">
      <template slot="brand">
        <b-navbar-item tag="router-link" :to="{ name: 'Home' }">
          <img
            src=""
            alt="FlatTrack logo"
            />
          <h1 class="title is-5" style="color: #fff;">{{ flatName }}</h1>
        </b-navbar-item>
      </template>
      <template slot="start">
        <b-navbar-item tag="router-link" :to="{ name: 'Home' }">
          Home
        </b-navbar-item>
        <b-navbar-item href="https://flattrack.io/help" tag="a" target="_blank" v-if="false">
          FlatTrack Help
        </b-navbar-item>
        <b-navbar-item tag="router-link" :to="{ name: 'My Flatmates', query: { 'group': 'admin' }}">
          Contact admin
        </b-navbar-item>
        <b-navbar-item tag="router-link" :to="{ name: 'My Flat' }">
          My flat: {{ flatName }}
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
  async beforeMount () {
    this.GetFlatInfo()
  }
}
</script>

<style>
#topbar {
    z-index: 100;
}

.gradient-blue {
    background-image: linear-gradient(141deg, #04a6d7, #209cee 71%, #3287f5);
}
</style>
