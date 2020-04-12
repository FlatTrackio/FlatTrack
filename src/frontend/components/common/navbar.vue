<template>
  <div class="navbar-z">
    <section>
      <div class="block">
        <b-sidebar
          type="is-white"
          fullheight="true"
          can-cancel="false"
          open="true">
          <div class="hero is-info is-bold navbar-text">
            <div class="block">
              <img
                src=""
                alt="FlatTrack logo"
                />
            </div>
            <h1 class="title is-2">FlatTrack</h1>
          </div>
          <div class="p-1">
            <b-menu>
              <b-menu-list label="General">
                <b-menu-item icon="home" label="Home" tag="router-link" to="/"></b-menu-item>
                <b-menu-item icon="information-outline" :label="flatName" tag="router-link" to="/flat" v-cloak></b-menu-item>
              </b-menu-list>
              <b-menu-list label="Apps">
                <b-menu-item icon="format-list-checks" label="Shopping list" tag="router-link" to="/apps/shopping-list"></b-menu-item>
                <b-menu-item icon="account-group" label="Flatmates" tag="router-link" to="/apps/flatmates"></b-menu-item>
                <b-menu-item icon="apps" label="Apps" tag="router-link" to="/apps"></b-menu-item>
              </b-menu-list>
              <b-menu-list label="Admin" v-if="canUserAccountAdmin">
                <b-menu-item icon="account-group" label="Flatmates" tag="router-link" to="/admin/accounts"></b-menu-item>
                <b-menu-item icon="settings" label="Admin apps" tag="router-link" to="/admin"></b-menu-item>
              </b-menu-list>
              <b-menu-list label="Help">
                <b-menu-item icon="open-in-new" label="FlatTrack help" tag="a" target="_blank" href="https://flattrack.io/help" disabled></b-menu-item>
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
import cani from '@/frontend/requests/authenticated/can-i'
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
      flatName: 'My flat',
      canUserAccountAdmin: false
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
    },
    CanIadmin () {
      cani.GetCanIgroup('admin').then(resp => {
        this.canUserAccountAdmin = resp.data.spec
      })
    }
  },
  async beforeMount () {
    this.GetFlatName()
    this.CanIadmin()
  }
}
</script>

<style>
.navbar-z {
    z-index: 100;
}

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
}

.navbar-text .title {
    color: white;
}

.align-bottom {
    bottom: 0;
    left: 0;
    width: 100%;
    height: 40px;
    background-color: #496c8a40;
    position: absolute;
}
</style>
