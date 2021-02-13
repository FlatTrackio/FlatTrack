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
  <div>
    <div class="container">
      <section class="section">
        <h1 class="title is-1">My Account</h1>
        <p class="subtitle is-4">Manage your account</p>
        <p></p>
        <br>
        <div v-for="app in apps" v-bind:key="app" >
          <div class="card pointer-cursor-on-hover" @click="$router.push({ name: app.routeName })">
            <div class="card-content">
              <div class="media">
                <div class="media-left">
                  <b-icon :icon="app.icon" size="is-medium"></b-icon>
                </div>
                <div class="media-content">
                  <p class="title is-3">{{ app.name }}</p>
                  <p class="subtitle is-5">{{ app.description }}</p>
                </div>
                <div class="media-right">
                  <b-icon icon="chevron-right" size="is-medium" type="is-midgray"></b-icon>
                </div>
              </div>
            </div>
            <div class="content">
            </div>
          </div>
          <br />
        </div>
        <div class="card pointer-cursor-on-hover" v-if="typeof CurrentTheme === 'object'" disabled>
          <div class="card-content">
            <div class="media">
              <div class="media-left">
                <b-icon :icon="CurrentTheme.icon" size="is-medium"></b-icon>
              </div>
              <div class="media-content">
                <p class="title is-3">Current theme</p>
                <p class="subtitle is-5">{{ CurrentTheme.name }}</p>
              </div>
            </div>
          </div>
          <div class="content">
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import theme from '@/common/theme'

export default {
  name: 'Account home',
  data () {
    return {
      apps: [
        {
          name: 'Profile',
          description: 'Manage your general information',
          icon: 'account-circle',
          routeName: 'Account Profile'
        },
        {
          name: 'Security',
          description: 'Manage your account security',
          icon: 'lock-question',
          routeName: 'Account Security'
        },
        {
          name: 'Settings',
          description: 'Manage settings for this device',
          icon: 'account-cog',
          routeName: 'Account Settings'
        }
      ]
    }
  },
  computed: {
    CurrentTheme () {
      return theme.GetTheme()
    }
  }
}
</script>

<style scoped>

</style>
