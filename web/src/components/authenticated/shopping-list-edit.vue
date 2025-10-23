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
  <div class="item-page">
    <div class="modal-card" style="width: auto">
      <header class="modal-card-head is-flex-wrap-wrap">
        <p v-if="!pageLoading" class="modal-card-title">
          {{ name || existingName }}
        </p>
        <b-skeleton v-else size="is-small" width="35%" :animated="true" />
        <p class="modal-card-subtitle">Rename shopping list and add notes</p>
      </header>
      <section class="modal-card-body">
        <div>
          <h3 class="title is-3">{{ shoppingListName }}</h3>
          <b-loading
            v-model:active="pageLoading"
            :is-full-page="false"
            :can-cancel="false"
          />
          <b-field label="Name">
            <b-input
              v-model="name"
              type="text"
              size="is-medium"
              maxlength="30"
              icon="text"
              placeholder="Enter a new name for this tag"
              icon-right="close-circle"
              icon-right-clickable
              required
              @icon-right-click="name = ''"
              @keyup.enter.native="UpdateShoppingList"
            />
          </b-field>
          <b-field label="Notes">
            <b-input
              ref="notes"
              v-model="notes"
              icon="text"
              size="is-medium"
              maxlength="100"
              type="text"
              placeholder="Enter extra information"
              icon-right="close-circle"
              icon-right-clickable
              @icon-right-click="notes = ''"
              @keyup.enter.native="UpdateShoppingList"
            />
          </b-field>
          <b-field addons>
            <b-button
              type="is-warning"
              size="is-medium"
              icon-left="arrow-left"
              native-type="submit"
              @click="$emit('close')"
            >
              Back
            </b-button>
            <b-button
              type="is-success"
              size="is-medium"
              icon-left="delta"
              native-type="submit"
              expanded
              :loading="submitLoading"
              :disabled="submitLoading"
              @click="UpdateShoppingList"
            >
              Update
            </b-button>
          </b-field>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
  import common from "@/common/common";
  import shoppinglist from "@/requests/authenticated/shoppinglist";
  import infotooltip from "@/components/common/info-tooltip.vue";

  export default {
    name: "ShoppingListEdit",
    components: {
      infotooltip,
    },
    props: {
      shoppingListId: String,
      existingName: String,
      existingNotes: String,
      completed: Boolean,
      totalTagExcludeList: [String],
    },
    data() {
      return {
        pageLoading: false,
        submitLoading: false,
        name: "",
        notes: "",
      };
    },
    methods: {
      UpdateShoppingList() {
        this.notesFromEmpty = false;
        this.editingMeta = false;
        this.editing = false;

        shoppinglist
          .UpdateShoppingList(this.shoppingListId, this.name, this.notes, this.completed, this.totalTagExcludeList)
          .then((resp) => {
            this.$emit("close");
          })
          .catch((err) => {
            common.DisplayFailureToast(
              this.$buefy,
              "Failed to update shopping list" +
                "<br/>" +
                err.response.data.metadata.response
            );
          });
      },
    },
    async mounted() {
      this.name = this.existingName;
      this.notes = this.existingNotes;
    },
  };
</script>
