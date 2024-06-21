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
            <li><router-link :to="{ name: 'Home' }">Home</router-link></li>
            <li class="is-active">
              <router-link :to="{ name: 'My Flat' }">My flat</router-link>
            </li>
          </ul>
          <b-button
            @click="CopyHrefToClipboard()"
            icon-left="content-copy"
            size="is-small"
          ></b-button>
        </nav>
        <h1 v-if="hasInitialLoaded || name !== ''" class="title is-1">
          {{ name }}
        </h1>
        <b-skeleton
          v-else
          size="is-medium"
          width="35%"
          :animated="true"
        ></b-skeleton>
        <p class="subtitle is-3">About your flat</p>
        <b-message type="is-primary" v-if="notes !== ''">
          <span v-for="line in notesSplit" v-bind:key="line">
            {{ line }}
            <br />
          </span>
        </b-message>
        <b-message type="is-warning" v-else>
          This section for describing such things as, but not limited to:
          <br />
          <ul style="list-style-type: disc">
            <li>how the flat life is</li>
            <li>rules</li>
            <li>regulations</li>
            <li>culture</li>
          </ul>
        </b-message>
        <b-button
          icon-left="pencil"
          v-if="canUserAccountAdmin === true"
          type="is-warning"
          @click="$router.push({ name: 'Admin settings' })"
          rounded
          >Edit message</b-button
        >
      </section>
    </div>
  </div>
</template>

<script>
import flatInfo from '@/requests/authenticated/flatInfo'
import cani from '@/requests/authenticated/can-i'
import common from '@/common/common'

export default {
  name: 'flat-home',
  data () {
    return {
      name: '',
      notes: '',
      notesSplit: '',
      hasInitialLoaded: false,
      canUserAccountAdmin: false
    }
  },
  methods: {
    ClipboardCopy () {
      common.ClipboardCopy()
    }
  },
  async beforeMount () {
    this.name = common.GetFlatnameFromCache() || this.name
    flatInfo
      .GetFlatName()
      .then((resp) => {
        if (this.name !== resp.data.spec) {
          this.name = resp.data.spec
          common.WriteFlatnameToCache(resp.data.spec)
        }
        return flatInfo.GetFlatNotes()
      })
      .then((resp) => {
        this.notes = resp.data.spec.notes
        this.notesSplit = this.notes.split('\n')
        this.hasInitialLoaded = true
      })
    cani.GetCanIgroup('admin').then((resp) => {
      this.canUserAccountAdmin = resp.data.data
    })
  }
}
</script>

<style src="../../assets/style.css"></style>
