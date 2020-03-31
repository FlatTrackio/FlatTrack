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
  }
}
</script>

<style src="../../assets/style.css"></style>

<style scoped>

</style>
