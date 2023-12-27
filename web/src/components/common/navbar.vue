<!--
     This program is free software: you can redistribute it and/or modify
     it under the terms of the Affero GNU General Public License as published by
     the Free Software Foundation, either version 3 of the License, or
     (at your option) any later version.

     This program is distributed in the hope that it will be useful,
     but WITHOUT ANY WARRANTY; without even the implied warranty of
     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
     GNU General Public License for more details.

     You should have received a copy of the Affero GNU General Public License
     along with this program.  If not, see <https://www.gnu.org/licenses/>.
-->

<template>
  <div class="navbar-z">
    <b-loading
      :is-full-page="false"
      :active.sync="pageLoading"
      :can-cancel="false"
    ></b-loading>
    <section>
      <div class="block">
        <b-sidebar
          type="is-white"
          fullheight="true"
          can-cancel="false"
          open="true"
        >
          <div class="hero is-info is-bold navbar-text navbar-shadow">
            <!-- <div class="block"> -->
            <!--   <img -->
            <!--     src="" -->
            <!--     alt="FlatTrack logo" -->
            <!--     /> -->
            <!-- </div> -->
            <h1 class="title is-2">FlatTrack</h1>
          </div>
          <div class="p-1">
            <b-menu>
              <b-menu-list label="General">
                <b-menu-item
                  icon="home"
                  label="Home"
                  tag="router-link"
                  to="/"
                ></b-menu-item>
                <b-menu-item
                  icon="information-outline"
                  label="My flat"
                  tag="router-link"
                  :to="{ name: 'My Flat' }"
                ></b-menu-item>
                <b-menu-item
                  icon="information-outline"
                  label="About FlatTrack"
                  tag="router-link"
                  :to="{ name: 'About FlatTrack' }"
                ></b-menu-item>
              </b-menu-list>
              <b-menu-list label="Apps">
                <b-menu-item
                  icon="format-list-checks"
                  label="Shopping list"
                  tag="router-link"
                  :to="{ name: 'Shopping list' }"
                ></b-menu-item>
                <b-menu-item
                  icon="account-group"
                  label="Flatmates"
                  tag="router-link"
                  :to="{ name: 'My Flatmates' }"
                ></b-menu-item>
                <b-menu-item
                  icon="apps"
                  label="Apps"
                  tag="router-link"
                  :to="{ name: 'Apps' }"
                ></b-menu-item>
              </b-menu-list>
              <b-menu-list label="Admin" v-if="canUserAccountAdmin">
                <b-menu-item
                  icon="account-group"
                  label="Flatmates"
                  tag="router-link"
                  :to="{ name: 'Admin accounts' }"
                ></b-menu-item>
                <b-menu-item
                  icon="settings"
                  label="Admin apps"
                  tag="router-link"
                  :to="{ name: 'Admin home' }"
                ></b-menu-item>
              </b-menu-list>
              <b-menu-list label="Help">
                <b-menu-item
                  icon="open-in-new"
                  label="FlatTrack help"
                  tag="a"
                  target="_blank"
                  href="https://flattrack.io/help"
                  disabled
                ></b-menu-item>
                <b-menu-item
                  icon="phone"
                  label="Contact admin"
                  tag="router-link"
                  :to="{ name: 'My Flatmates', query: { group: 'admin' } }"
                ></b-menu-item>
              </b-menu-list>
              <b-menu-list label="Account">
                <b-menu-item
                  icon="account-circle"
                  label="Profile"
                  tag="router-link"
                  :to="{ name: 'Account Profile' }"
                ></b-menu-item>
                <b-menu-item
                  icon="account"
                  label="My Account"
                  tag="router-link"
                  :to="{ name: 'Account' }"
                ></b-menu-item>
                <b-menu-item
                  icon="exit-to-app"
                  label="Sign out"
                  @click="signOut"
                ></b-menu-item>
              </b-menu-list>
            </b-menu>
          </div>
        </b-sidebar>
      </div>
    </section>
  </div>
</template>

<script>
import common from '@/common/common'
import cani from '@/requests/authenticated/can-i'

export default {
  name: 'nav-bar',
  data () {
    return {
      isActive: true,
      open: true,
      overlay: false,
      fullheight: true,
      fullwidth: false,
      pageLoading: true,
      flatName: 'My flat',
      canUserAccountAdmin: false
    }
  },
  methods: {
    signOut () {
      common.SignoutDialog()
    }
  },
  async beforeMount () {
    cani.GetCanIgroup('admin').then((resp) => {
      this.canUserAccountAdmin = resp.data.data
      this.pageLoading = false
    })
  }
}
</script>

<style>
.navbar-z {
  z-index: 100;
}

.p-1 {
  padding-left: 1em !important;
  padding-bottom: 1em !important;
  padding-right: 1em !important;
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
