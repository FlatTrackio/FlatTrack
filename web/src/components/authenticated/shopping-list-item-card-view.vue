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
    <div class="card">
      <div class="card-content card-content-list">
        <div class="media">
          <div
            class="media-left"
            @click="PatchItemObtained(item.id, !obtained)"
          >
            <b-checkbox v-model="obtained" size="is-medium" />
          </div>
          <div
            class="media-content pointer-cursor-on-hover"
            @click="$emit('viewItem')"
          >
            <div class="block">
              <p
                :class="obtained === true ? 'obtained' : ''"
                class="subtitle is-4 m-0"
              >
                {{ item.name }}
                <span
                  v-if="typeof item.price !== 'undefined' && item.price !== 0"
                >
                  (${{ item.price.toFixed(2) }})
                </span>
                <b v-if="item.quantity > 1">x{{ item.quantity }} </b>
                <b-icon
                  v-if="typeof item.price === 'undefined' || item.price === 0"
                  icon="currency-usd-off"
                  type="is-lightred"
                  size="is-small"
                />
              </p>
              <span>
                <p class="subtitle is-6">
                  <b-icon
                    v-if="displayTag === true"
                    icon="tag-multiple"
                    type="is-info"
                    size="is-small"
                  />
                  <span v-if="displayTag === true"> {{ item.tag }} </span>
                  <span
                    v-if="
                      displayTag === true &&
                        typeof item.tag !== 'undefined' &&
                        typeof item.notes !== 'undefined' &&
                        item.notes !== ''
                    "
                  >
                    -
                  </span>
                  <i> {{ item.notes }} </i>
                </p>
              </span>
            </div>
          </div>
          <div class="media-right is-flex">
            <b-field>
              <b-tooltip label="Delete" class="is-paddingless" :delay="200">
                <b-button
                  size="is-small"
                  type="is-danger"
                  icon-right="delete"
                  :loading="itemDeleting"
                  @click="DeleteShoppingListItem(item.id, index)"
                />
              </b-tooltip>
              <span class="pointer-cursor-on-hover" @click="$emit('viewItem')">
                <b-icon
                  icon="chevron-right"
                  size="is-medium"
                  type="is-midgray"
                />
              </span>
            </b-field>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import common from "@/common/common";
  import shoppinglist from "@/requests/authenticated/shoppinglist";

  export default {
    name: "ShoppingListItemCardView",
    props: {
      item: Object,
      listId: String,
      list: Object,
      index: Number,
      displayTag: Boolean,
      deviceIsMobile: Boolean,
      itemDisplayState: Number,
    },
    data() {
      return {
        itemDeleting: false,
        obtained: false,
      };
    },
    created() {
      this.obtained = this.item.obtained;
    },
    methods: {
      PatchItemObtained(itemId, obtained) {
        this.$emit("obtained", obtained);
        shoppinglist
          .PatchShoppingListItemObtained(this.listId, itemId, obtained)
          .then(() => {
            var displayAll =
              typeof this.itemDisplayState === "number" &&
              this.itemDisplayState === 0;
            if (displayAll === true) {
              return;
            }
            let removedFromList = this.list;
            removedFromList.splice(this.index, 1);
            this.$emit("list", removedFromList);
          })
          .catch((err) => {
            common.DisplayFailureToast(
              this.$buefy,
              "Failed to patch the obtained field of this item" +
                "<br/>" +
                err.response.data.metadata.response
            );
          });
      },
      DeleteShoppingListItem(itemId, index) {
        this.$buefy.dialog.confirm({
          title: "Delete item",
          message:
            "Are you sure that you wish to delete this shopping list item?" +
            "<br/>" +
            "This action cannot be undone.",
          confirmText: "Delete item",
          type: "is-danger",
          hasIcon: true,
          onConfirm: () => {
            this.itemDeleting = true;
            shoppinglist
              .DeleteShoppingListItem(this.listId, itemId)
              .then((resp) => {
                common.DisplaySuccessToast(
                  this.$buefy,
                  resp.data.metadata.response
                );
                let removedFromList = this.list;
                removedFromList.splice(this.index, 1);
                this.$emit("list", removedFromList);
              })
              .catch((err) => {
                common.DisplayFailureToast(
                  this.$buefy,
                  "Failed to delete shopping list item" +
                    " - " +
                    err.response.data.metadata.response
                );
                this.itemDeleting = false;
              });
          },
        });
      },
    },
  };
</script>

<style scoped>
  .display-is-editable:hover {
    text-decoration: underline dotted;
    -webkit-transition: width 0.5s ease-in;
  }
  .card-content-list {
    background-color: transparent;
    padding-left: 1.5em;
    padding-top: 0.6em;
    padding-bottom: 0.6em;
    padding-right: 1.5em;
  }

  .obtained {
    color: #adadad;
    text-decoration: line-through;
  }
</style>
