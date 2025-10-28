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
        <p class="modal-card-subtitle">
          Rename list tag or delete items by tag
        </p>
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
              @keyup.enter.native="UpdateShoppingListItemTag"
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
              @click="UpdateShoppingListItemTag"
            >
              Update tag
            </b-button>
            <b-button
              type="is-danger"
              size="is-medium"
              icon-left="delete"
              native-type="submit"
              :loading="deleteLoading"
              @click="DeleteShoppingListTagItems"
            />
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
    name: "ShoppingListTagNameEdit",
    components: {
      infotooltip,
    },
    props: {
      shoppingListId: String,
      shoppingListName: String,
      existingName: String,
    },
    data() {
      return {
        pageLoading: false,
        submitLoading: false,
        deleteLoading: false,
        name: "",
      };
    },
    methods: {
      UpdateShoppingListItemTag() {
        this.submitLoading = true;
        shoppinglist
          .UpdateShoppingListItemTag(
            this.shoppingListId,
            this.existingName,
            this.name
          )
          .then((resp) => {
            this.$emit("close");
          })
          .catch((err) => {
            common.DisplayFailureToast(
              this.$buefy,
              "Failed to update the shopping list tag" +
                "<br/>" +
                err.response.data.metadata.response
            );
          });
      },
      DeleteShoppingListTagItems() {
        this.$buefy.dialog.confirm({
          title: "Delete shopping list tag items",
          message:
            `Are you sure that you wish to delete '${this.name}'?` +
            "<br/>" +
            "This action cannot be undone.",
          confirmText: "Delete tag items",
          type: "is-danger",
          hasIcon: true,
          onConfirm: () => {
            this.deleteLoading = true;
            shoppinglist
              .DeleteShoppingListTagItems(this.shoppingListId, this.name)
              .then((resp) => {
                this.deleteLoading = false;
                common.DisplaySuccessToast(
                  this.$buefy,
                  "Deleted the shopping list tag"
                );
                this.$emit("close");
              })
              .catch((err) => {
                this.deleteLoading = false;
                common.DisplayFailureToast(
                  this.$buefy,
                  "Failed to delete the shopping list" +
                    "<br/>" +
                    err.response.data.metadata.response
                );
              });
          },
        });
      },
    },
    async mounted() {
      this.name = this.existingName;
    },
  };
</script>
