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
        <b-loading
          v-model:active="pageLoading"
          :is-full-page="false"
          :can-cancel="false"
        />
        <h1 class="title is-1">
          Home
        </h1>
        <div
          v-if="notes !== '' || canUserAccountAdmin === true"
          class="my-4"
        >
          <p class="subtitle is-3">
            About {{ name }}
          </p>
          <b-message
            v-if="notes !== ''"
            type="is-primary"
          >
            <span
              v-for="line in notesSplit"
              :key="line"
            >
              {{ line }}
              <br>
            </span>
          </b-message>
          <b-message
            v-else
            type="is-warning"
          >
            This section for describing such things as, but not limited to:
            <br>
            <ul style="list-style-type: disc">
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
            rounded
            @click="$router.push({ name: 'Admin settings' })"
          >
            Edit message
          </b-button>
        </div>
        <p class="subtitle is-3">
          Recent activity
        </p>
        <h1 class="subtitle is-4">
          Shopping lists
        </h1>
        <div v-if="lists.length > 0">
          <shoppingListCardView
            v-for="(list, index) in lists"
            :key="list"
            :list="list"
            :authors="authors"
            :lists="lists"
            :index="index"
            :device-is-mobile="true"
            :mini="true"
          />
          <p class="is-size-5 ml-3 mb-2">
            <b-icon
              icon="party-popper"
              type="is-success"
              size="is-medium"
            />
            You're all caught up! That's all for now
          </p>
          <b-button
            class="has-text-left"
            type="is-info"
            icon-left="close-box-outline"
            expanded
            @click="ClearItems"
          >
            Clear
          </b-button>
        </div>
        <div v-else>
          <b-message type="is-warning">
            You're all caught up. Check back later!
          </b-message>
        </div>
        <div v-if="canUserAccountAdmin === true">
          <br>
          <h1 class="subtitle is-4">
            Admin
          </h1>
          <div
            class="card pointer-cursor-on-hover"
            @click="$router.push({ name: 'Admin accounts' })"
          >
            <div class="card-content card-content-list">
              <div class="media">
                <div class="media-left">
                  <b-icon
                    icon="account-group"
                    size="is-medium"
                  />
                </div>
                <div class="media-content">
                  <div class="block">
                    <p class="subtitle is-4">
                      Review flatmember accounts
                    </p>
                    <p class="subtitle is-6">
                      Add new and remove previous flatmate accounts
                    </p>
                  </div>
                </div>
                <div class="media-right">
                  <b-icon
                    icon="chevron-right"
                    size="is-medium"
                    type="is-midgray"
                  />
                </div>
              </div>
            </div>
            <div class="content" />
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
  name: 'FlattrackHome',
  components: {
    shoppingListCardView: () =>
      import('@/components/authenticated/shopping-list-card-view.vue')
  },
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
      })
    cani.GetCanIgroup('admin').then((resp) => {
      this.canUserAccountAdmin = resp.data.data
    })
    this.GetShoppingLists()
  },
  methods: {
    GetShoppingLists () {
      shoppinglist
        .GetShoppingLists(
          undefined,
          this.sortBy,
          undefined,
          this.modificationTimestampAfter,
          5
        )
        .then((resp) => {
          this.pageLoading = false
          this.lists = resp.data.list || []
        })
        .catch(() => {
          common.DisplayFailureToast(
            'Hmmm seems somethings gone wrong loading the shopping lists'
          )
        })
    },
    ClearItems () {
      common.WriteHomeLastViewedTimestamp(Number(dayjs().unix()))
      this.lists = []
    }
  }
}
</script>

<style scoped></style>
