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
          back-link-name="Shopping list"
          :current-page-name="$route.name"
        />
        <div>
          <h1 class="title is-1">
            New shopping list
          </h1>
          <p class="subtitle is-3">
            Start a new list for your next shop
          </p>
          <b-field label="Name">
            <b-input
              v-model="name"
              type="text"
              maxlength="30"
              icon="form-textbox"
              size="is-medium"
              placeholder="Enter a title for this list"
              autofocus
              icon-right="close-circle"
              icon-right-clickable
              required
              @icon-right-click="name = ''"
              @keyup.enter.native="PostNewShoppingList"
            />
          </b-field>
          <b-field label="Notes (optional)">
            <b-input
              v-model="notes"
              type="text"
              size="is-medium"
              icon="text"
              placeholder="Enter extra information"
              icon-right="close-circle"
              icon-right-clickable
              maxlength="100"
              @keyup.enter.native="PostNewShoppingList"
              @icon-right-click="notes = ''"
            />
          </b-field>
          <b-field
            v-if="lists.length > 0"
            label="Template list (optional)"
          >
            <b-select
              v-model="listTemplate"
              placeholder="Optionally select a list to base a new list off"
              icon="content-copy"
              expanded
              size="is-medium"
            >
              <option value="">
                No template
              </option>
              <option disabled />
              <option
                v-for="list in lists"
                :key="list.id"
                :value="list.id"
              >
                {{ list.name }}
              </option>
            </b-select>
          </b-field>
          <div
            v-if="listTemplate !== '' && typeof listTemplate !== 'undefined'"
            class="field"
          >
            <label class="label"> Select items </label>
            <div class="field">
              <b-radio
                v-model="templateListItemSelector"
                size="is-medium"
                name="itemSelector"
                native-value="all"
              >
                All items
              </b-radio>
            </div>
            <div class="field">
              <b-radio
                v-model="templateListItemSelector"
                size="is-medium"
                name="itemSelector"
                native-value="unobtained"
              >
                Only from unobtained
              </b-radio>
            </div>
            <div class="field">
              <b-radio
                v-model="templateListItemSelector"
                size="is-medium"
                name="itemSelector"
                native-value="obtained"
              >
                Only from obtained
              </b-radio>
            </div>
          </div>
          <b-button
            icon-left="plus"
            type="is-success"
            size="is-medium"
            native-type="submit"
            expanded
            :loading="submitLoading"
            :disabled="submitLoading"
            @click="
              PostNewShoppingList(
                name,
                notes,
                listTemplate,
                templateListItemSelector
              )
            "
          >
            Create list
          </b-button>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
  import common from "@/common/common";
  import shoppinglist from "@/requests/authenticated/shoppinglist";
  import breadcrumb from "@/components/common/breadcrumb.vue";

  export default {
    name: "ShoppingListNew",
    components: {
      breadcrumb,
    },
    data() {
      return {
        name: "",
        notes: "",
        listTemplate: "",
        templateListItemSelector: "all",
        lists: [],
      };
    },
    async beforeMount() {
      shoppinglist.GetShoppingLists(undefined, "templated").then((resp) => {
        this.lists = resp.data.list || [];
      });
      this.name = this.$route.query.name;
    },
    methods: {
      CopyHrefToClipboard() {
        common.CopyHrefToClipboard();
      },
      PostNewShoppingList() {
        if (this.notes === "") {
          this.notes = undefined;
        }
        this.submitLoading = true;
        shoppinglist
          .PostShoppingList(
            this.name,
            this.notes,
            this.listTemplate,
            this.templateListItemSelector
          )
          .then((resp) => {
            this.submitLoading = false;
            var list = resp.data.spec;
            if (list.id !== "" || typeof list.id === "undefined") {
              this.$router.push({
                name: "View shopping list",
                params: { id: list.id },
              });
            } else {
              common.DisplayFailureToast(
                "Unable to find created shopping list"
              );
            }
          })
          .catch((err) => {
            this.submitLoading = false;
            common.DisplayFailureToast(
              `Failed to create shopping list - ${err.response.data.metadata.response}`
            );
          });
      },
    },
  };
</script>
