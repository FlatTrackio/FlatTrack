<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><router-link to="/apps">Apps</router-link></li>
              <li><router-link to="/apps/shopping-list">Shopping list</router-link></li>
              <li class="is-active"><router-link to="/apps/shopping-list/new">New shopping list</router-link></li>
            </ul>
        </nav>
        <h1 class="title is-1">New shopping list</h1>
        <b-field label="Name">
          <b-input type="text"
                   v-model="name"
                   maxlength="30"
                   required>
          </b-input>
        </b-field>
        <b-field label="Notes">
          <b-input type="text"
                   v-model="notes"
                   maxlength="30"
                   >
          </b-input>
        </b-field>
        <b-button type="is-success" size="is-medium" rounded native-type="submit" @click="PostNewShoppingList(name, notes)">Create</b-button>
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
      flatName: 'New shopping list',
      name: '',
      notes: ''
    }
  },
  methods: {
    PostNewShoppingList (name, notes) {
      if (notes === '') {
        notes = undefined
      }
      shoppinglist.PostShoppingList(name, notes).then(resp => {
        var list = resp.data.spec
        if (list.id !== '' || typeof list.id === 'undefined') {
          this.$router.push({ path: `/apps/shopping-list/lists/${list.id}` })
        } else {
          common.DisplayFailureToast('Unable to find created shopping list')
        }
      }).catch(err => {
        common.DisplayFailureToast(`Failed to create shopping list - ${err.response.data.metadata.response}`)
      })
    }
  }
}
</script>
