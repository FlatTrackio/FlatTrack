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
          back-link-name="Admin home"
          :current-page-name="$route.name"
        />
        <h1 class="title is-1">Settings</h1>
        <p class="subtitle is-4">General FlatTrack settings</p>
        <b-loading
          v-model:active="pageLoading"
          :is-full-page="false"
          :can-cancel="false"
        />
        <div>
          <label class="label">Flat name</label>
          <b-field>
            <b-input
              v-model="flatName"
              type="text"
              maxlength="20"
              placeholder="Enter your flat's name"
              icon="form-textbox"
              size="is-medium"
              icon-right="close-circle"
              icon-right-clickable
              expanded
              required
              @icon-right-click="flatName = ''"
              @keyup.enter.native="PostFlatName"
            />
            <p class="control">
              <b-button
                type="is-primary"
                size="is-medium"
                icon-left="check"
                @click="PostFlatName"
              />
            </p>
          </b-field>
        </div>
        <div>
          <label class="label">Flat notes</label>
          <b-field>
            <b-input
              v-model="flatNotes"
              type="textarea"
              minlength="0"
              maxlength="500"
              placeholder="Enter notes about your flat. e.g: living space, rules, obligations, etc..."
              size="is-medium"
              icon-right="close-circle"
              icon-right-clickable
              expanded
              @icon-right-click="flatNotes = ''"
            />
            <p class="control">
              <b-button
                type="is-primary"
                size="is-medium"
                icon-left="check"
                @click="PutFlatNotes"
              />
            </p>
          </b-field>
        </div>
        <b-field grouped>
          <b-field label="Shopping List Cleanup" expanded>
            <b-field label="Keep lists" expanded>
              <b-select
                placeholder="Forever"
                expanded
                v-model="shoppingListKeepPolicy"
                size="is-medium"
              >
                <option value="Always">Always</option>
                <option value="ThreeMonths">For three months</option>
                <option value="SixMonths">For six months</option>
                <option value="OneYear">For one year</option>
                <option value="TwoYears">For two years</option>
                <option value="Last10">The last 10</option>
                <option value="Last50">The last 50</option>
                <option value="Last100">The last 100</option>
              </b-select>
            </b-field>
          </b-field>
        </b-field>
      </section>
    </div>
  </div>
</template>

<script>
  import flatInfo from "@/requests/authenticated/flatInfo";
  import settings from "@/requests/admin/settings";
  import common from "@/common/common";
  import breadcrumb from "@/components/common/breadcrumb.vue";

  export default {
    name: "AdminSettings",
    components: {
      breadcrumb,
    },
    data() {
      return {
        pageLoading: true,
        flatName: "",
        flatNotes: "",
        shoppingListKeepPolicy: "",
      };
    },
    async beforeMount() {
      flatInfo
        .GetFlatName()
        .then((resp) => {
          this.flatName = resp.data.spec;
          return flatInfo.GetFlatNotes();
        })
        .then((resp) => {
          this.flatNotes = resp.data.spec.notes;
          return settings.GetShoppingListKeepPolicy();
        })
        .then((resp) => {
          this.shoppingListKeepPolicy = resp.data.spec;
          this.pageLoading = false;
        });
    },
    methods: {
      CopyHrefToClipboard() {
        common.CopyHrefToClipboard();
      },
      PostFlatName() {
        if (this.flatName === "") {
          common.DisplayFailureToast(
            this.$buefy,
            "Error: Flat name must not be empty"
          );
          return;
        }
        settings
          .PostFlatName(this.flatName)
          .then((resp) => {
            common.DisplaySuccessToast(this.$buefy, "Set flat name");
          })
          .catch((err) => {
            common.DisplayFailureToast(
              this.$buefy,
              "Failed set the flat name" +
                "<br/>" +
                (err.response.data.metadata.response || err)
            );
          });
      },
      PutFlatNotes() {
        if (this.flatName === "") {
          common.DisplayFailureToast(
            this.$buefy,
            "Error: Flat notes must not be empty"
          );
          return;
        }
        settings
          .PutFlatNotes(this.flatNotes)
          .then((resp) => {
            common.DisplaySuccessToast(this.$buefy, "Set flat notes");
          })
          .catch((err) => {
            common.DisplayFailureToast(
              this.$buefy,
              "Failed set the flat notes" + "<br/>" + err
            );
          });
      },
      TimestampToCalendar(timestamp) {
        return common.TimestampToCalendar(timestamp);
      },
    },
    watch: {
      shoppingListKeepPolicy() {
        if (this.pageLoading === true || this.shoppingListKeepPolicy === "") {
          return;
        }
        settings
          .PutShoppingListKeepPolicy(this.shoppingListKeepPolicy)
          .then((resp) => {
            common.DisplaySuccessToast(
              this.$buefy,
              "Set shopping list keep policy"
            );
          })
          .catch((err) => {
            common.DisplayFailureToast(
              this.$buefy,
              "Failed set the shopping list keep policy" + "<br/>" + err
            );
          });
      },
    },
  };
</script>
