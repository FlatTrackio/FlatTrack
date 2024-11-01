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
    <div
      class="modal-card"
      style="width: auto"
    >
      <header class="modal-card-head">
        <p class="modal-card-title">
          New shopping item
        </p>
        <p class="modal-card-subtitle">
          Add an item to the list
        </p>
      </header>
      <section class="modal-card-body">
        <div>
          <b-field
            label="Name"
            class="is-marginless"
          >
            <b-input
              v-model="name"
              type="text"
              size="is-medium"
              maxlength="30"
              icon="text"
              placeholder="Enter a name for this item"
              icon-right="close-circle"
              icon-right-clickable
              autofocus
              required
              @icon-right-click="name = ''"
              @keyup.enter="PostShoppingListItem"
            />
          </b-field>
          <b-field
            label="Notes (optional)"
            class="is-marginless"
          >
            <b-input
              v-model="notes"
              type="text"
              size="is-medium"
              icon="text"
              placeholder="Enter information extra"
              icon-right="close-circle"
              icon-right-clickable
              maxlength="40"
              @keyup.enter="PostShoppingListItem"
              @icon-right-click="notes = ''"
            />
          </b-field>
          <b-field label="Price (optional)">
            <b-input
              v-model="price"
              type="number"
              step="0.01"
              placeholder="0.00"
              icon="currency-usd"
              icon-right="close-circle"
              icon-right-clickable
              size="is-medium"
              @icon-right-click="price = ''"
              @keyup.enter="PostShoppingListItem"
            />
          </b-field>
          <b-field label="Quantity">
            <b-numberinput
              v-model="quantity"
              size="is-medium"
              placeholder="Enter how many of this item should be obtained"
              min="0"
              expanded
              required
              controls-position="compact"
              icon="numeric"
            />
          </b-field>
          <div>
            <div class="field has-addons">
              <label class="label">Tag (optional)</label>
              <p class="control">
                <infotooltip
                  message="To manage tags, navigate to the Apps -> Shopping List -> Manage tags page"
                />
              </p>
            </div>
            <b-field class="is-marginless">
              <b-dropdown>
                <template #trigger>
                  <b-button
                    v-if="tagsList.length > 0 || tags.length > 0"
                    icon-left="menu-down"
                    type="is-primary"
                    size="is-medium"
                  />
                </template>

                <b-dropdown-item
                  v-if="tagsList.length > 0"
                  disabled
                >
                  Tags in this list
                </b-dropdown-item>
                <div
                  v-for="existingListTag in tagsList"
                  :key="existingListTag"
                >
                  <b-dropdown-item
                    v-if="
                      existingListTag !== '' &&
                        existingListTag.length > 0 &&
                        typeof existingListTag !== 'undefined'
                    "
                    :value="existingListTag"
                    @click="tag = existingListTag"
                  >
                    {{ existingListTag }}
                  </b-dropdown-item>
                </div>
                <b-dropdown-item
                  v-if="tags.length > 0"
                  disabled
                >
                  Tags in all lists
                </b-dropdown-item>
                <div
                  v-for="existingTag in tags"
                  :key="existingTag"
                >
                  <b-dropdown-item
                    v-if="
                      existingTag.name !== '' &&
                        existingTag.name.length > 0 &&
                        typeof existingTag.name !== 'undefined'
                    "
                    :value="existingTag.name"
                    @click="tag = existingTag.name"
                  >
                    {{ existingTag.name }}
                  </b-dropdown-item>
                </div>
              </b-dropdown>
              <b-input
                v-model="tag"
                expanded
                type="text"
                icon="tag"
                maxlength="30"
                placeholder="Enter a tag to group the item"
                icon-right="close-circle"
                icon-right-clickable
                size="is-medium"
                @keyup.enter="PostShoppingListItem"
                @icon-right-click="tag = ''"
              />
            </b-field>
          </div>
          <b-field label="Obtained">
            <b-checkbox
              v-model="obtained"
              size="is-medium"
            >
              Obtained
            </b-checkbox>
          </b-field>
          <p
            v-if="typeof price !== 'undefined' && price !== 0 && quantity > 1"
            class="m-1"
          >
            Total price with quantity: ${{ itemCurrentPrice.toFixed(2) }}
          </p>
          <b-field grouped>
            <b-button
              type="is-warning"
              size="is-medium"
              icon-left="arrow-left"
              native-type="submit"
              @click="$parent.close()"
            >
              Back
            </b-button>
            <b-button
              type="is-success"
              size="is-medium"
              icon-left="plus"
              native-type="submit"
              expanded
              :loading="submitLoading"
              :disabled="submitLoading"
              @click="PostShoppingListItem"
            >
              Add item
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

  export default {
    name: "ShoppingItemNew",
    components: {
      infotooltip: () => import("@/components/common/info-tooltip.vue"),
    },
    props: {
      withTag: {
        type: String,
        default: null
      },
      withName: {
        type: String,
        default: null
      }
    },
    data() {
      return {
        shoppingListId: this.$route.params.id,
        shoppingListName: "",
        tags: [],
        tagsList: [],
        submitLoading: false,
        name: "",
        notes: "",
        price: 0,
        quantity: 1,
        tag: "",
        obtained: false,
      };
    },
    computed: {
      itemCurrentPrice() {
        return this.price * this.quantity;
      },
    },
    async beforeMount() {
      if (this.withTag && this.tag === "") {
        this.tag = this.withTag;
      }
      if (this.withName && this.name === "") {
        this.name = this.withName;
      }
      shoppinglist
        .GetShoppingList(this.shoppingListId)
        .then((resp) => {
          var list = resp.data.spec;
          this.shoppingListName = list.name;
          return shoppinglist.GetAllShoppingListItemTags();
        })
        .then((resp) => {
          this.tags = resp.data.list || [];
          return shoppinglist.GetShoppingListItemTags(this.shoppingListId);
        })
        .then((resp) => {
          this.tagsList = resp.data.list || [];
        });
    },
    methods: {
      CopyHrefToClipboard() {
        common.CopyHrefToClipboard();
      },
      PostShoppingListItem() {
        this.submitLoading = true;
        if (this.notes === "") {
          this.notes = undefined;
        }
        if (this.price === 0) {
          this.price = undefined;
        } else {
          this.price = parseFloat(this.price);
        }
        if (this.tag === "") {
          this.tag = "Untagged";
        }

        shoppinglist
          .PostShoppingListItem(
            this.shoppingListId,
            this.name,
            this.notes,
            this.price,
            this.quantity,
            this.tag,
            this.obtained
          )
          .then((resp) => {
            var item = resp.data.spec;
            if (item.id !== "" || typeof item.id === "undefined") {
              this.$parent.close();
            } else {
              this.submitLoading = false;
              common.DisplayFailureToast(
                "Unable to find created shopping item"
              );
            }
          })
          .catch((err) => {
            this.submitLoading = false;
            common.DisplayFailureToast(
              `Failed to add shopping list item - ${err.response.data.metadata.response}`
            );
          });
      },
    },
  };
</script>
