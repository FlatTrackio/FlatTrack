<template>
    <div>
        <div class="container">
            <section class="section is-large">
                <h2 class="title is-2">Unknown page!</h2>
                <h4 class="title is-4">Hmmm maybe something went wrong {{ emojiConfused }}</h4>
                <b-button tag="router-link" :to="{ name: 'Home' }" rounded>Go home</b-button>
                <b-button type="is-warning" @click="$router.go(-1)" rounded>Go back</b-button>
            </section>
        </div>
    </div>
</template>

<script>
import emoji from 'node-emoji'

export default {
  name: 'Unknown page',
  data () {
    return {
      emojiConfused: emoji.get('confused')
    }
  }
}
</script>

<style scoped>

</style>
