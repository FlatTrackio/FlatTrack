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
          My Account
        </h1>
        <p class="subtitle is-4">
          Manage your account
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
        <div
          v-if="typeof CurrentTheme === 'object'"
          class="card pointer-cursor-on-hover"
        >
          <div
            class="card-content"
            @click="NextTheme"
          >
            <div class="media">
              <div class="media-left">
                <b-icon
                  :icon="CurrentTheme.icon"
                  size="is-medium"
                />
              </div>
              <div class="media-content">
                <p class="title is-3">
                  Current theme
                </p>
                <p class="subtitle is-5">
                  {{ CurrentTheme.name }}
                </p>
              </div>
            </div>
          </div>
          <div class="content" />
        </div>
      </section>
    </div>
  </div>
</template>

<script>
  import theme from "@/common/theme";

  export default {
    name: "AccountHome",
    data() {
      return {
        apps: [
          {
            name: "Profile",
            description: "Manage your general information",
            icon: "account-circle",
            routeName: "Account Profile",
          },
          {
            name: "Security",
            description: "Manage your account security",
            icon: "lock-question",
            routeName: "Account Security",
          },
          {
            name: "Settings",
            description: "Manage settings for this device",
            icon: "account-cog",
            routeName: "Account Settings",
          },
        ],
      CurrentTheme:  theme.GetTheme(),
      };
    },
    methods: {
      NextTheme() {
        const currentTheme = theme.GetTheme();
        const themes = theme.ListThemes();
        const currentThemeIndex = themes.findIndex(i => i.name === currentTheme.name)
        const nextThemeIndex = currentThemeIndex + 1
        if (currentThemeIndex === -1) {
          theme.SetThemeDefault()
          return
        }
        const nextTheme = themes[nextThemeIndex % themes.length]
        theme.SetTheme(nextTheme.name)
        this.CurrentTheme = theme.GetTheme()
      },
    },
  };
</script>

<style scoped></style>
