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
      <section class="section is-large">
        <h2 class="title is-2">
          Unknown page!
        </h2>
        <h4 class="title is-4">
          Hmmm maybe something went wrong {{ emojiConfused }}
        </h4>
        <b-button
          tag="router-link"
          :to="{ name: 'Home' }"
          rounded
        >
          Go home
        </b-button>
        <b-button
          type="is-warning"
          rounded
          @click="$router.go(-1)"
        >
          Go back
        </b-button>
      </section>
    </div>
  </div>
</template>

<script>
import emoji from 'node-emoji'

export default {
  name: 'UnknownPage',
  data () {
    return {
      emojiConfused: emoji.get('confused')
    }
  }
}
</script>

<style scoped></style>
