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
        <p class="subtitle is-3">Recent activity of your flat</p>

        <h1 class="subtitle is-4">Shopping lists</h1>
        <div v-if="lists.length > 0">
          <shoppingListCardView :list="list" :authors="authors" :lists="lists" :index="index" v-for="(list, index) in lists" v-bind:key="list" :deviceIsMobile="true" :mini="true" />
          <p class="is-size-5 ml-3">
            <b-icon icon="party-popper" type="is-success" size="is-medium"></b-icon>
            You're all caught up! That's all for now
          </p>
          <!-- TODO add clear button -->
        </div>
        <div v-else>
          <p class="is-size-5">Nothing recent yet. Check back in later!</p>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/common/common'
import shoppinglist from '@/requests/authenticated/shoppinglist'
export default {
  name: 'home',
  data () {
    return {
      authors: {},
      lists: [],
      pageLoading: true,
      deviceIsMobile: false,
      sortBy: 'recentlyUpdated'
    }
  },
  methods: {
    // TODO add modificationTimestampAfter value
    GetShoppingLists () {
      shoppinglist.GetShoppingLists(undefined, this.sortBy, undefined, undefined).then(resp => {
        this.pageLoading = false
        this.lists = resp.data.list || []
        console.log(this.lists)
      }).catch(() => {
        common.DisplayFailureToast('Hmmm seems somethings gone wrong loading the shopping lists')
      })
    }
  },
  components: {
    shoppingListCardView: () => import('@/components/authenticated/shopping-list-card-view.vue')
  },
  async beforeMount () {
    this.GetShoppingLists()
  }
}
</script>

<style scoped>
</style>
