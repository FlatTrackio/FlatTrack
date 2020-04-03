<template>
  <div id="topbar">
    <b-navbar type="is-info" shadow="true" transparent="true">
      <template slot="brand">
        <b-navbar-item tag="router-link" :to="{ path: '/' }">
          <h1 class="title is-5" style="color: #fff;">FlatTrack (My Flat)</h1>
        </b-navbar-item>
      </template>
      <template slot="start">
        <b-navbar-item tag="router-link" :to="{ path: '/' }">
          Home
        </b-navbar-item>
        <b-navbar-item href="https://flattrack.io/help">
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
import { DialogProgrammatic as Dialog, LoadingProgrammatic as Loading } from 'buefy'

export default {
  name: 'topbar',
  methods: {
    signOut () {
      Dialog.confirm({
        message: 'Are you sure you want to sign out?',
        onConfirm: () => {
          const loadingComponent = Loading.open({
            container: null
          })
          setTimeout(() => {
            common.DeleteAuthToken()
            window.location.href = '/login'
          }, 1 * 1000)
        }
      })
    }
  }
}
</script>

<style>
</style>
