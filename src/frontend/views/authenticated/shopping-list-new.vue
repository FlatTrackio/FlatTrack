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
              required>
            </b-input>
          </b-field>
          <b-field label="Notes">
            <b-input
              type="text"
              size="is-medium"
              v-model="notes"
              icon="text"
              maxlength="100">
            </b-input>
          </b-field>
          <b-field label="Use another list as a template (optional)" v-if="lists.length > 0">
            <b-select
              placeholder="Template a preview list"
              v-model="listTemplate"
              icon="content-copy"
              expanded
              size="is-medium">
              <option
                value="">
              </option>
              <option
                v-for="list in lists"
                :value="list.id"
                :key="list.id">
                {{ list.name }}
                </option>
            </b-select>
          </b-field>
          <div class="field" v-if="listTemplate !== '' && typeof listTemplate !== 'undefined'">
            <b-checkbox v-model="templateListFromOnlyUnobtained">Create list only from unobtained items in template list</b-checkbox>
          </div>
          <br/>
          <b-button
            icon-left="plus"
            type="is-success"
            size="is-medium"
            native-type="submit"
            @click="PostNewShoppingList(name, notes, listTemplate, templateListFromOnlyUnobtained)">
            Create
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
      templateListOnlyFromUnobtained: false,
      lists: []
    }
  },
  methods: {
    PostNewShoppingList (name, notes, listTemplate, templateListFromOnlyUnobtained) {
      if (notes === '') {
        notes = undefined
      }
      shoppinglist.PostShoppingList(name, notes, listTemplate, templateListFromOnlyUnobtained).then(resp => {
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
