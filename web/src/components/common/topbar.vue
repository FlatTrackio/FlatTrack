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
  <div id="topbar">
    <b-navbar class="navbar-shadow gradient-blue" :fixed-top="ratherSmallerScreen !== true" type="is-info" transparent="false">
      <template slot="brand">
        <b-navbar-item tag="router-link" :to="{ name: 'Home' }">
          <!-- <img -->
          <!--   src="" -->
          <!--   alt="FlatTrack logo" -->
          <!--   /> -->
          <h1 class="title is-5" style="color: #fff;">
            FlatTrack
            <span style="font-weight: normal;">
              {{ flatName }}
            </span>
          </h1>
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
          My flat
        </b-navbar-item>
        <b-navbar-item tag="router-link" :to="{ name: 'About FlatTrack' }">
          About FlatTrack
        </b-navbar-item>
        <b-navbar-item @click="signOut">
          Sign out
        </b-navbar-item>
      </template>
    </b-navbar>
  </div>
</template>

<script>
import common from '@/common/common'
import flatInfo from '@/requests/authenticated/flatInfo'
import { DialogProgrammatic as Dialog, LoadingProgrammatic as Loading } from 'buefy'

export default {
  name: 'topbar',
  data () {
    return {
      flatName: 'My Flat',
      ratherSmallerScreen: false
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
    if (window.innerWidth <= 330) {
      this.ratherSmallerScreen = true
    }
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
