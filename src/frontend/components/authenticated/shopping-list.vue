<template>
  <div>
    <section v-for="item in items" v-bind:key="item">
      <div class="card">
        <div class="card-content">
          <div class="media">
            <div class="media-left">
              <b-checkbox size="is-medium" v-model="item.obtained" @click="!item.obtained"></b-checkbox>
            </div>
            <div class="media-content pointer-cursor-on-hover" @click="goToRef('/apps/shopping-list/list/' + listId + '/item/' + item.id)">
              <p :class="item.obtained === true ? 'obtained' : ''" class="subtitle is-4">{{ item.name }}</p>
            </div>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script>
export default {
  name: 'shopping list table',
  data () {
    return {
      totalPrice: 0.0
    }
  },
  methods: {
    goToRef (ref) {
      this.$router.push({ path: ref })
    }
  },
  props: {
    items: Array,
    listId: String
  }
}
</script>

<style scoped>
.card-content {
    background-color: transparent;
    padding-left: 1.5em;
    padding-top: 0.6em;
    padding-bottom: 0.6em;
    padding-right: 1.5em;
}

.obtained {
    color: #adadad;
    text-decoration: line-through;
}
</style>
