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
        <b-loading :is-full-page="false" :active.sync="pageLoading" :can-cancel="false"></b-loading>
        <h1 class="title is-1">Home</h1>
        <div v-if="notes !== '' || canUserAccountAdmin === true" class="my-4">
          <p class="subtitle is-3">About {{ name }}</p>
          <b-message type="is-primary" v-if="notes !== ''">
            <span v-for="line in notesSplit" v-bind:key="line">
              {{ line }}
              <br/>
            </span>
          </b-message>
          <b-message type="is-warning" v-else>
            This section for describing such things as, but not limited to:
            <br/>
            <ul style="list-style-type: disc;">
              <li>how the flat life is</li>
              <li>rules</li>
              <li>regulations</li>
              <li>culture</li>
            </ul>
          </b-message>
          <b-button
            v-if="canUserAccountAdmin === true"
            icon-left="pencil"
            type="is-warning"
            @click="$router.push({ name: 'Admin settings' })" rounded>
            Edit message
          </b-button>
        </div>
        <p class="subtitle is-3">Recent activity</p>
        <h1 class="subtitle is-4">Shopping lists</h1>
        <div v-if="lists.length > 0">
          <shoppingListCardView :list="list" :authors="authors" :lists="lists" :index="index" v-for="(list, index) in lists" v-bind:key="list" :deviceIsMobile="true" :mini="true" />
          <p class="is-size-5 ml-3 mb-2">
            <b-icon icon="party-popper" type="is-success" size="is-medium"></b-icon>
            You're all caught up! That's all for now
          </p>
          <b-button
            class="has-text-left"
            @click="ClearItems"
            type="is-info"
            icon-left="close-box-outline"
            expanded>
            Clear
          </b-button>
        </div>
        <div v-else>
          <b-message type="is-warning">
            You're all caught up. Check back later!
          </b-message>
        </div>
        <div v-if="canUserAccountAdmin === true">
          <br/>
          <h1 class="subtitle is-4">Admin</h1>
          <div class="card pointer-cursor-on-hover" @click="$router.push({ name: 'Admin accounts' })">
            <div class="card-content card-content-list">
              <div class="media">
                <div class="media-left">
                  <b-icon icon="account-group" size="is-medium"></b-icon>
                </div>
                <div class="media-content">
                  <div class="block">
                    <p class="subtitle is-4">Review flatmember accounts</p>
                    <p class="subtitle is-6">Add new and remove previous flatmate accounts</p>
                  </div>
                </div>
                <div class="media-right">
                  <b-icon icon="chevron-right" size="is-medium" type="is-midgray"></b-icon>
                </div>
              </div>
            </div>
            <div class="content">
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/common/common'
import shoppinglist from '@/requests/authenticated/shoppinglist'
import flatInfo from '@/requests/authenticated/flatInfo'
import cani from '@/requests/authenticated/can-i'
import dayjs from 'dayjs'

export default {
  name: 'home',
  data () {
    return {
      name: '',
      authors: {},
      pageLoading: true,
      lists: [],
      deviceIsMobile: false,
      modificationTimestampAfter: common.GetHomeLastViewedTimestamp(),
      notes: '',
      notesSplit: '',
      sortBy: 'recentlyUpdated',
      canUserAccountAdmin: false
    }
  },
  methods: {
    GetShoppingLists () {
      shoppinglist.GetShoppingLists(undefined, this.sortBy, undefined, this.modificationTimestampAfter, 5).then(resp => {
        this.pageLoading = false
        this.lists = resp.data.list || []
      }).catch(() => {
        common.DisplayFailureToast('Hmmm seems somethings gone wrong loading the shopping lists')
      })
    },
    ClearItems () {
      common.WriteHomeLastViewedTimestamp(Number(dayjs().unix()))
      this.lists = []
    }
  },
  components: {
    shoppingListCardView: () => import('@/components/authenticated/shopping-list-card-view.vue')
  },
  async beforeMount () {
    this.name = common.GetFlatnameFromCache() || this.name
    flatInfo.GetFlatName().then(resp => {
      if (this.name !== resp.data.spec) {
        this.name = resp.data.spec
        common.WriteFlatnameToCache(resp.data.spec)
      }
      return flatInfo.GetFlatNotes()
    }).then(resp => {
      this.notes = resp.data.spec.notes
      this.notesSplit = this.notes.split('\n')
    })
    cani.GetCanIgroup('admin').then(resp => {
      this.canUserAccountAdmin = resp.data.data
    })
    this.GetShoppingLists()
  }
}
</script>

<style scoped>
</style>
