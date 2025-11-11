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
        <breadcrumb back-link-name="Apps" :current-page-name="$route.name" />
        <h1 class="title is-1">Shopping list</h1>
        <p class="subtitle is-3">Manage your weekly shop</p>
        <div
          class="mb-5"
          v-if="(notes !== '' || canUserAccountAdmin) && !pageLoading"
        >
          <div class="content">
            <label class="label">Notes</label>
            <p
              :class="
                canUserAccountAdmin
                  ? 'display-is-editable pointer-cursor-on-hover'
                  : ''
              "
              class="subtitle is-4 notes-highlight"
              @click="EditShoppingListNotes"
            >
              <i> {{ notes || "Add notes" }} </i>
            </p>
          </div>
        </div>
        <b-skeleton
          class="mb-5"
          v-else
          size="is-medium"
          width="35%"
          :animated="true"
        />
        <b-button
          class="has-text-left"
          type="is-info"
          icon-left="tag-multiple"
          expanded
          @click="$router.push({ name: 'Manage shopping tags' })"
        >
          Manage tags
        </b-button>
        <br />
        <div>
          <b-tabs
            v-model="listDisplayState"
            :position="deviceIsMobile ? 'is-centered' : ''"
            class="block is-marginless"
          >
            <b-tab-item icon="format-list-checks" label="All" />
            <b-tab-item icon="playlist-remove" label="Uncompleted" />
            <b-tab-item icon="playlist-check" label="Completed" />
          </b-tabs>
          <label class="label">Search for lists</label>
          <b-field>
            <b-input
              ref="search"
              v-model="listSearch"
              icon="magnify"
              size="is-medium"
              placeholder="Enter a list name"
              type="search"
              expanded
            />
            <p class="control">
              <b-select
                v-model="sortBy"
                placeholder="Sort by"
                icon="sort"
                size="is-medium"
                expanded
              >
                <option value="recentlyAdded">Recently Added</option>
                <option value="lastAdded">Last Added</option>
                <option value="recentlyUpdated">Recently Updated</option>
                <option value="lastUpdated">Last Updated</option>
                <option value="alphabeticalDescending">A-z</option>
                <option value="alphabeticalAscending">z-A</option>
              </b-select>
            </p>
          </b-field>
          <b-loading
            v-model:active="pageLoading"
            :is-full-page="false"
            :can-cancel="false"
          />
          <section>
            <div
              class="card pointer-cursor-on-hover mb-5"
              @click="
                $router.push({
                  name: 'New shopping list',
                  query: { name: listSearch || undefined },
                })
              "
            >
              <div class="card-content">
                <div class="media">
                  <div class="media-left">
                    <b-icon icon="cart-plus" size="is-medium" />
                  </div>
                  <div class="media-content">
                    <p class="title is-4">Add a new list</p>
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
            </div>
          </section>
        </div>
        <floatingAddButton
          v-if="displayFloatingAddButton"
          :router-link="{
            name: 'New shopping list',
            query: { name: listSearch || undefined },
          }"
        />
        <div v-if="listsFiltered.length > 0">
          <shoppingListCardView
            v-for="(list, index) in listsFiltered"
            :key="list"
            :list="list"
            :authors="authors"
            :lists="lists"
            :index="index"
            :device-is-mobile="deviceIsMobile"
            @lists="
              (l) => {
                lists = l;
              }
            "
          />
          <div class="m-3">
            <p>{{ listsFilteredTotal.length }} shopping list(s)</p>
          </div>
          <b-pagination
            :total="listsFilteredTotal.length"
            v-model="page"
            order="is-centered"
            :per-page="itemsPerPage"
          />
        </div>
        <div v-else>
          <div class="card">
            <div class="card-content card-content-list">
              <div class="media">
                <div class="media-left">
                  <b-icon
                    icon="cart-remove"
                    size="is-medium"
                    type="is-midgray"
                  />
                </div>
                <div class="media-content">
                  <p
                    v-if="
                      listSearch === '' && lists.length === 0 && !pageLoading
                    "
                    class="subtitle is-4"
                  >
                    No lists added yet.
                  </p>
                  <p
                    v-else-if="
                      listSearch === '' &&
                        listDisplayState === 1 &&
                        lists.length > 0 &&
                        !pageLoading
                    "
                    class="subtitle is-4"
                  >
                    All lists have been completed.
                  </p>
                  <p
                    v-else-if="
                      listSearch === '' &&
                        listDisplayState === 2 &&
                        lists.length > 0 &&
                        !pageLoading
                    "
                    class="subtitle is-4"
                  >
                    No lists have been completed yet.
                  </p>
                  <p
                    v-else-if="listSearch !== '' && !pageLoading"
                    class="subtitle is-4"
                  >
                    No lists found.
                  </p>
                  <b-skeleton
                    class="mb-5"
                    v-else
                    size="is-medium"
                    width="35%"
                    :animated="true"
                  />
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
  import common from "@/common/common";
  import shoppinglist from "@/requests/authenticated/shoppinglist";
  import flatmates from "@/requests/authenticated/flatmates";
  import shoppingListCardView from "@/components/authenticated/shopping-list-card-view.vue";
  import floatingAddButton from "@/components/common/floating-add-button.vue";
  import cani from "@/requests/authenticated/can-i";
  import breadcrumb from "@/components/common/breadcrumb.vue";

  export default {
    name: "ShoppingList",
    components: {
      shoppingListCardView,
      floatingAddButton,
      breadcrumb,
    },
    data() {
      return {
        displayFloatingAddButton: true,
        canUserAccountAdmin: false,
        notes: "",
        lists: [],
        authors: {},
        listDisplayState: 0,
        deviceIsMobile: false,
        listSearch: "",
        pageLoading: true,
        sortBy: "recentlyUpdated",
        totalAmount: 0,
        page: 1,
        itemsPerPage: 5,
      };
    },
    computed: {
      listsFilteredTotal() {
        return this.lists.filter((item) => {
          return this.ListDisplayState(item);
        });
      },
      listsFiltered() {
        return this.lists
          .filter((item) => {
            return this.ListDisplayState(item);
          })
          .slice(
            (this.page - 1) * this.itemsPerPage,
            this.page * this.itemsPerPage
          );
      },
    },
    watch: {
      sortBy() {
        this.listIsLoading = true;
        this.GetShoppingLists();
      },
    },
    async beforeMount() {
      cani.GetCanIgroup("admin").then((resp) => {
        this.canUserAccountAdmin = resp.data.data;
      });
      this.GetShoppingLists();
      this.GetShoppingListNotes();
    },
    beforeUnmount() {
      window.removeEventListener("resize", this.CheckDeviceIsMobile, true);
    },
    async created() {
      this.CheckDeviceIsMobile();
      window.addEventListener("resize", this.CheckDeviceIsMobile, true);
    },
    methods: {
      CopyHrefToClipboard() {
        common.CopyHrefToClipboard();
      },
      GetShoppingLists() {
        shoppinglist
          .GetShoppingLists(
            undefined,
            this.sortBy,
            undefined,
            undefined,
            undefined,
            undefined
          )
          .then((resp) => {
            this.lists = resp.data.list || [];
            this.pageLoading = false;
            this.totalAmount = resp.data.data || 0;
          })
          .catch(() => {
            common.DisplayFailureToast(
              "Hmmm seems somethings gone wrong loading the shopping lists"
            );
          });
      },
      GetShoppingListNotes() {
        shoppinglist
          .GetShoppingListNotes()
          .then((resp) => {
            this.notes = resp.data.spec || "";
          })
          .catch(() => {
            common.DisplayFailureToast(
              "Hmmm seems somethings gone wrong loading the notes for shopping lists"
            );
          });
      },
      EditShoppingListNotes() {
        if (this.canUserAccountAdmin !== true) {
          return;
        }
        this.displayFloatingAddButton = false;
        this.$buefy.dialog.prompt({
          title: "Shopping list notes",
          message: `Enter notes that are useful for shopping in your flat.`,
          container: null,
          icon: "text",
          hasIcon: true,
          inputAttrs: {
            placeholder:
              "e.g. Our budget is $200/w. Please make sure to bring the supermarket card.",
            maxlength: 250,
            required: false,
            value: this.notes || undefined,
          },
          trapFocus: true,
          onConfirm: (value) => {
            this.pageLoading = true;
            shoppinglist
              .PutShoppingListNotes(value)
              .then(() => {
                this.pageLoading = false;
                this.displayFloatingAddButton = true;
                this.GetShoppingListNotes();
              })
              .catch((err) => {
                common.DisplayFailureToast(
                  "Failed to update notes" +
                    `<br/>${err.response.data.metadata.response}`
                );
                this.displayFloatingAddButton = true;
                this.pageLoading = false;
              });
          },
          onCancel: () => {
            this.displayFloatingAddButton = true;
          },
        });
      },
      GetFlatmateName(id) {
        flatmates
          .GetFlatmate(id)
          .then((resp) => {
            return resp.data.spec.names;
          })
          .catch((err) => {
            common.DisplayFailureToast(
              "Failed to fetch user account" +
                `<br/>${err.response.data.metadata.response}`
            );
            return id;
          });
      },
      ListDisplayState(list) {
        if (this.listDisplayState === 1 && list.completed === false) {
          return this.ItemByNameInList(list);
        } else if (this.listDisplayState === 2 && list.completed === true) {
          return this.ItemByNameInList(list);
        } else if (this.listDisplayState === 0) {
          return this.ItemByNameInList(list);
        }
      },
      ItemByNameInList(item) {
        var vm = this;
        return (
          item.name.toLowerCase().indexOf(vm.listSearch.toLowerCase()) !== -1
        );
      },
      CheckDeviceIsMobile() {
        this.deviceIsMobile = common.DeviceIsMobile();
      },
    },
  };
</script>

<style scoped></style>
