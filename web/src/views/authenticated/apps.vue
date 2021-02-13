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
        <h1 class="title is-1">Apps</h1>
        <p class="subtitle is-4">Manage your flat</p>
        <p></p>
        <br>
        <div v-for="app in apps" v-bind:key="app" >
          <div class="card pointer-cursor-on-hover" @click="$router.push({ name: app.routeName })">
            <div class="card-content">
              <div class="media">
                <div class="media-left">
                  <b-icon :icon="app.icon" size="is-medium"></b-icon>
                </div>
                <div class="media-content">
                  <p class="title is-3">{{ app.name }}</p>
                  <p class="subtitle is-5">{{ app.description }}</p>
                </div>
                <div class="media-right">
                  <b-icon icon="chevron-right" size="is-medium" type="is-midgray"></b-icon>
                </div>
              </div>
            </div>
            <div class="content">
            </div>
          </div>
          <br />
        </div>
      </section>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Apps',
  data () {
    return {
      apps: [
        {
          name: 'Shopping list',
          description: 'Create and manage lists to use for your next shop',
          icon: 'format-list-checks',
          routeName: 'Shopping list'
        },
        {
          name: 'Flatmates',
          description: 'View your flatmates and their contact details',
          icon: 'account-group',
          routeName: 'My Flatmates'
        }
      ]
    }
  }
}
</script>

<style scoped>

</style>
