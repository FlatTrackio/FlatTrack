<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><router-link to="/apps/shopping-list">Shopping list</router-link></li>
              <li class="is-active"><router-link to="/apps/shopping-list/new">New shopping list</router-link></li>
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
          <br/>
          <b-button
            icon-left="plus"
            type="is-success"
            size="is-medium"
            native-type="submit"
            expanded
            @click="PostNewShoppingList(name, notes, listTemplate, templateListItemSelector)">
            Create list
          </b-button>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/frontend/common/common'
import shoppinglist from '@/frontend/requests/authenticated/shoppinglist'

export default {
  name: 'shopping-list-new',
  data () {
    return {
      name: '',
      notes: '',
      listTemplate: '',
      templateListItemSelector: undefined,
      lists: []
    }
  },
  methods: {
    PostNewShoppingList (name, notes, listTemplate, templateListItemSelector) {
      if (notes === '') {
        notes = undefined
      }
      shoppinglist.PostShoppingList(name, notes, listTemplate, templateListItemSelector).then(resp => {
        var list = resp.data.spec
        if (list.id !== '' || typeof list.id === 'undefined') {
          this.$router.push({ path: `/apps/shopping-list/list/${list.id}` })
        } else {
          common.DisplayFailureToast('Unable to find created shopping list')
        }
      }).catch(err => {
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
