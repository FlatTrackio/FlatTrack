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
    <headerDisplay />
    <div class="container">
      <section class="section">
        <h1 class="title is-1">Unavailable</h1>
        <b-message has-icon type="is-warning">
          <p
            v-if="typeof message !== 'undefined' && message !== ''"
            class="subtitle is-4"
          >
            {{ message }}
          </p>
          <p v-else class="subtitle is-4">
            This FlatTrack instance is unavailable or not healthy.
          </p>
        </b-message>
        <b-button
          type="is-white"
          size="is-medium"
          rounded
          expanded
          @click="Reload"
        >
          Refresh
        </b-button>
      </section>
    </div>
  </div>
</template>

<script>
  import headerDisplay from "@/components/common/header-display.vue";
  import common from "@/common/common";
  export default {
    name: "UnavailablePage",
    components: {
      headerDisplay,
    },
    data() {
      return {
        message: common.GetMaintenanceModeMessage() || undefined,
      };
    },
    methods: {
      Reload() {
        this.$router.go();
      },
    },
  };
</script>

<style scoped></style>
