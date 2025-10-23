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
        <breadcrumb
          back-link-name="Home"
          :current-page-name="$route.name"
        />
        <h1
          v-if="hasInitialLoaded || name !== ''"
          class="title is-1"
        >
          {{ name }}
        </h1>
        <b-skeleton
          v-else
          size="is-medium"
          width="35%"
          :animated="true"
        />
        <p class="subtitle is-3">
          About your flat
        </p>
        <b-message
          v-if="notes !== ''"
          type="is-primary"
        >
          <span
            v-for="line in notesSplit"
            :key="line"
          >
            {{ line }}
            <br>
          </span>
        </b-message>
        <b-message
          v-else
          type="is-warning"
        >
          This section for describing such things as, but not limited to:
          <br>
          <ul style="list-style-type: disc">
            <li>how the flat life is</li>
            <li>rules</li>
            <li>regulations</li>
            <li>culture</li>
          </ul>
        </b-message>
        <b-button
          v-if="canUserAccountAdmin === true"
          icon-left="pencil"
          type="is-warning"
          rounded
          @click="$router.push({ name: 'Admin settings' })"
        >
          Edit message
        </b-button>
      </section>
    </div>
  </div>
</template>

<script>
  import flatInfo from "@/requests/authenticated/flatInfo";
  import cani from "@/requests/authenticated/can-i";
  import common from "@/common/common";
  import breadcrumb from "@/components/common/breadcrumb.vue";

  export default {
    name: "FlatHome",
    components: {
      breadcrumb,
    },
    data() {
      return {
        name: "",
        notes: "",
        notesSplit: "",
        hasInitialLoaded: false,
        canUserAccountAdmin: false,
      };
    },
    async beforeMount() {
      this.name = common.GetFlatnameFromCache() || this.name;
      flatInfo
        .GetFlatName()
        .then((resp) => {
          if (this.name !== resp.data.spec) {
            this.name = resp.data.spec;
            common.WriteFlatnameToCache(resp.data.spec);
          }
          return flatInfo.GetFlatNotes();
        })
        .then((resp) => {
          this.notes = resp.data.spec.notes;
          this.notesSplit = this.notes.split("\n");
          this.hasInitialLoaded = true;
        });
      cani.GetCanIgroup("admin").then((resp) => {
        this.canUserAccountAdmin = resp.data.data;
      });
    },
    methods: {
      CopyHrefToClipboard() {
        common.CopyHrefToClipboard();
      },
    },
  };
</script>
