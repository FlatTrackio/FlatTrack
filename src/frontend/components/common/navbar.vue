<template>
  <div>
    <section>
      <div class="block">
        <b-sidebar
          type="is-light"
          :fullheight="fullheight"
          :fullwidth="fullwidth"
          :overlay="overlay"
          :right="right"
          can-cancel="false"
          :open.sync="open">
          <div class="navbar-text">
            <h1 class="title is-2">FlatTrack</h1>
          </div>
          <div class="p-1">
            <b-menu>
              <b-menu-list label="General">
                <b-menu-item icon="home" label="Home" tag="router-link" to="/"></b-menu-item>
                <b-menu-item icon="information-outline" :label="flatName" tag="router-link" to="/flat"></b-menu-item>
              </b-menu-list>
              <b-menu-list label="Apps">
                <b-menu-item icon="format-list-checks" label="Shopping list" tag="router-link" to="/apps/shopping-list"></b-menu-item>
                <b-menu-item icon="account-group" label="Flatmates" tag="router-link" to="/apps/flatmates"></b-menu-item>
                <b-menu-item icon="apps" label="Apps" tag="router-link" to="/apps"></b-menu-item>
              </b-menu-list>
              <b-menu-list label="Admin">
                <b-menu-item icon="account-group" label="Flatmates" tag="router-link" to="/admin/flatmates"></b-menu-item>
                <b-menu-item icon="settings" label="Settings" tag="router-link" to="/admin"></b-menu-item>
              </b-menu-list>
              <b-menu-list label="Help">
                <b-menu-item icon="open-in-new" label="FlatTrack help" tag="a" target="_blank" href="https://flattrack.io/help"></b-menu-item>
                <b-menu-item icon="phone" label="Contact admin" tag="router-link" to="/apps/flatmates?group=admin"></b-menu-item>
              </b-menu-list>
              <b-menu-list label="Account">
                <b-menu-item icon="account" label="Profile" tag="router-link" to="/profile"></b-menu-item>
                <b-menu-item icon="exit-to-app" label="Sign out" @click="signOut"></b-menu-item>
              </b-menu-list>
            </b-menu>
          </div>
        </b-sidebar>
      </div>
    </section>
  </div>
</template>

<script>
import common from '@/frontend/common/common'
import flatInfo from '@/frontend/requests/authenticated/flatInfo'
import { DialogProgrammatic as Dialog, LoadingProgrammatic as Loading } from 'buefy'

export default {
  name: 'navbar',
  data () {
    return {
      isActive: true,
      open: true,
      overlay: false,
      fullheight: true,
      fullwidth: false,
      right: false,
      flatName: 'My flat'
    }
  },
  methods: {
    signOut () {
      common.SignoutDialog()
    },
    GetFlatName () {
      flatInfo.GetFlatName().then(resp => {
        this.flatName = resp.data.spec
      })
    }
  },
  async created () {
    this.GetFlatName()
  }
}
</script>

<style>
.p-1 {
    padding-left: 1em;
    padding-bottom: 1em;
    padding-right: 1em;
}

.p-1 .menu {
    margin-top: 23px;
}

.navbar-text {
    text-align: center;
    padding: 15px;
    background-color: #209cee;
}

.navbar-text .title {
    color: white;
}
</style>
