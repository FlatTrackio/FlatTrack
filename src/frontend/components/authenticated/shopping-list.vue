<template>
  <div>
    <headerDisplay />
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
          <ul>
            <li>
              <a href="/">Home</a>
            </li>
            <li class="is-active">
              <a href="/shopping-list">Shopping list</a>
            </li>
          </ul>
        </nav>
        <h1 class="title">Shopping list</h1>
        <h2 class="subtitle">Manage your weekly shop</h2>
        <div id="tasks" v-if="items && items.length">
          <div class="card-margin" v-for="item of items" v-bind:key="item">
            <div class="card">
              <div class="card-content">
                <div class="media">
                  <div class="media-content">
                    <div class="field">
                      <b-checkbox size="is-large" class="title is-4">
                        {{ item.name }}
                      </b-checkbox>
                      <div class="content">
                        ${{ item.price }}
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div id="tasks" v-if="!items.length">
          <section class="section">
            <div class="card">
              <div class="card-content">
                <div class="media">
                  <div class="media-content">
                    <p class="title is-4">Nothing's on the list</p>
                  </div>
                </div>
                <div class="content">
                  Feel free to get started on your shopping list
                  <br />
                </div>
              </div>
            </div>
          </section>
        </div>
    </div>
  </div>
</template>

<script>
import headerDisplay from '@/components/header-display'

export default {
  name: 'Shopping List',
  data () {
    return {
      items: [
        {
          id: 0,
          name: 'Cheese',
          amount: 1,
          price: '10',
          week: 42,
          timestamp: 1571818843,
          addedBy: 'a798f50a-dfef-475c-958f-86b87c915a96',
          standard: true,
          obtained: false
        },
        {
          id: 1,
          name: 'Bananas',
          amount: 2,
          price: '3',
          week: 42,
          timestamp: 1571818843,
          addedBy: 'a798f50a-dfef-475c-958f-86b87c915a96',
          standard: true,
          obtained: false
        }
      ]
    }
  },
  components: {
    headerDisplay
  }
}
</script>

<style scoped>

</style>
