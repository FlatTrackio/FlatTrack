<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
            <li><a href="/">Home</a></li>
            <li class="is-active"><a href="/apps">Apps</a></li>
            </ul>
        </nav>
        <h1 class="title">Apps</h1>
        <p></p>
        <br>
        <div class="tile is-ancestor">
          <div class="tile is-parent" @click="goToApp('/apps/shopping-list')">
            <div class="tile is-child box">
              <p class="title">Shopping list</p>
              <p>Manage a flat-wide shopping list for all to contribute to</p>
            </div>
            <!-- TODO add '>' arrow at right side of this box && change mouse-over to click (?) -->
          </div>
        </div>
        <div class="tile is-ancestor">
          <div class="tile is-parent" @click="goToApp('/apps/flatmates')">
            <div class="tile is-child box">
              <p class="title">Flatmates</p>
              <p>View your flatmates and their contact details</p>
            </div>
            <!-- TODO add '>' arrow at right side of this box && change mouse-over to click (?) -->
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import emoji from 'node-emoji'

export default {
  name: 'Apps',
  data () {
    return {
      members: [],
      emojiSmile: emoji.get('smile')
    }
  },
  created () {
  },
  methods: {
    goToApp: (ref) => {
      window.location.href = ref
    }
  }
}
</script>

<style src="../../assets/style.css"></style>

<style scoped>

</style>
