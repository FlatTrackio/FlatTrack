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
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><router-link :to="{ name: 'Shopping list' }">Shopping list</router-link></li>
              <li class="is-active"><router-link :to="{ name: 'New shopping list' }">New shopping list</router-link></li>
            </ul>
        </nav>
        <div>
          <h1 class="title is-1">New shopping list</h1>
          <p class="subtitle is-3">Start a new list for your next shop</p>
          <b-field label="Name">
            <b-input
              type="text"
              v-model="name"
              maxlength="30"
              icon="textbox"
              size="is-medium"
              placeholder="Enter a title for this list"
              autofocus
              icon-right="close-circle"
              icon-right-clickable
              @icon-right-click="name = ''"
              @keyup.enter.native="PostNewShoppingList"
              required>
            </b-input>
          </b-field>
          <b-field label="Notes (optional)">
            <b-input
              type="text"
              size="is-medium"
              v-model="notes"
              icon="text"
              placeholder="Enter extra information"
              @keyup.enter.native="PostNewShoppingList"
              icon-right="close-circle"
              icon-right-clickable
              @icon-right-click="notes = ''"
              maxlength="100">
            </b-input>
          </b-field>
          <b-field label="Template list (optional)" v-if="lists.length > 0">
            <b-select
              placeholder="Optionally select a list to base a new list off"
              v-model="listTemplate"
              icon="content-copy"
              expanded
              size="is-medium">
              <option
                value="">
                No template
              </option>
              <option disabled></option>
              <option
                v-for="list in lists"
                :value="list.id"
                :key="list.id">
                {{ list.name }}
                </option>
            </b-select>
          </b-field>
          <div class="field" v-if="listTemplate !== '' && typeof listTemplate !== 'undefined'">
            <label class="label">
              Select items
            </label>
            <div class="field">
              <b-radio
                v-model="templateListItemSelector"
                size="is-medium"
                name="itemSelector"
                native-value="all">
                All items
              </b-radio>
            </div>
            <div class="field">
              <b-radio
                v-model="templateListItemSelector"
                size="is-medium"
                name="itemSelector"
                native-value="unobtained">
                Only from unobtained
              </b-radio>
            </div>
            <div class="field">
              <b-radio
                v-model="templateListItemSelector"
                size="is-medium"
                name="itemSelector"
                native-value="obtained">
                Only from obtained
              </b-radio>
            </div>
          </div>
          <b-button
            icon-left="plus"
            type="is-success"
            size="is-medium"
            native-type="submit"
            expanded
            :loading="submitLoading"
            :disabled="submitLoading"
            @click="PostNewShoppingList(name, notes, listTemplate, templateListItemSelector)">
            Create list
          </b-button>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/common/common'
import shoppinglist from '@/requests/authenticated/shoppinglist'

export default {
  name: 'shopping-list-new',
  data () {
    return {
      name: '',
      notes: '',
      listTemplate: '',
      templateListItemSelector: 'all',
      lists: []
    }
  },
  methods: {
    PostNewShoppingList () {
      if (this.notes === '') {
        this.notes = undefined
      }
      this.submitLoading = true
      shoppinglist.PostShoppingList(this.name, this.notes, this.listTemplate, this.templateListItemSelector).then(resp => {
        this.submitLoading = false
        var list = resp.data.spec
        if (list.id !== '' || typeof list.id === 'undefined') {
          this.$router.push({ name: 'View shopping list', params: { id: list.id } })
        } else {
          common.DisplayFailureToast('Unable to find created shopping list')
        }
      }).catch(err => {
        this.submitLoading = false
        common.DisplayFailureToast(`Failed to create shopping list - ${err.response.data.metadata.response}`)
      })
    }
  },
  async beforeMount () {
    shoppinglist.GetShoppingLists().then(resp => {
      this.lists = resp.data.list || []
    })
  }
}
</script>
