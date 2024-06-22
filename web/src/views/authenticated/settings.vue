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
        <nav
          class="breadcrumb is-medium has-arrow-separator"
          aria-label="breadcrumbs"
        >
          <ul>
            <li>
              <router-link :to="{ name: 'Account' }">My account</router-link>
            </li>
            <li class="is-active">
              <router-link :to="{ name: 'Account Settings' }"
                >Settings</router-link
              >
            </li>
          </ul>
          <b-button
            @click="CopyHrefToClipboard()"
            icon-left="content-copy"
            size="is-small"
          ></b-button>
        </nav>
        <h1 class="title is-1">Settings</h1>
        <p class="subtitle is-3">Manage settings for this device</p>
        <b-loading
          :is-full-page="false"
          :active.sync="pageLoading"
          :can-cancel="false"
        ></b-loading>

        <b-field label="Miscellaneous">
          <b-checkbox size="is-medium" v-model="enableAnimations">
            Enable Animations
          </b-checkbox>
        </b-field>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/common/common'

export default {
  name: 'settings-home',
  data () {
    return {
      enableAnimations: common.GetEnableAnimations() !== 'false'
    }
  },
  methods: {
    CopyHrefToClipboard () {
      common.CopyHrefToClipboard()
    }
  },
  watch: {
    enableAnimations () {
      common.WriteEnableAnimations(this.enableAnimations)
      if (this.enableAnimations !== 'false') {
        common.Hooray()
      }
    }
  }
}
</script>
