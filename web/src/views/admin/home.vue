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
        <h1 class="title is-1">
          Admin apps
        </h1>
        <p class="subtitle is-4">
          Manage FlatTrack
        </p>
        <p />
        <br>
        <div
          v-for="app in apps"
          :key="app"
          class="mb-5"
        >
          <div
            class="card pointer-cursor-on-hover"
            @click="$router.push({ name: app.routeName })"
          >
            <div class="card-content">
              <div class="media">
                <div class="media-left">
                  <b-icon
                    :icon="app.icon"
                    size="is-medium"
                  />
                </div>
                <div class="media-content">
                  <p class="title is-3">
                    {{ app.name }}
                  </p>
                  <p class="subtitle is-5">
                    {{ app.description }}
                  </p>
                </div>
                <div class="media-right">
                  <b-icon
                    icon="chevron-right"
                    size="is-medium"
                    type="is-midgray"
                  />
                </div>
              </div>
            </div>
            <div class="content" />
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
  export default {
    name: "AdminHome",
    data() {
      return {
        apps: [
          {
            name: "Accounts",
            description: "Manage your flatmates accounts",
            icon: "account-group",
            routeName: "Admin accounts",
          },
          {
            name: "Settings",
            description: "General FlatTrack settings",
            icon: "cogs",
            routeName: "Admin settings",
          },
        ],
      };
    },
  };
</script>

<style scoped></style>
