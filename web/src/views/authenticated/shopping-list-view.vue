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
    <div
      v-if="HeaderIsSticky"
      :class="ratherSmallerScreen ? 'ListBar ListBarTop' : 'ListBar'"
    >
      <div class="level">
        <p class="subtitle is-5 mb-0 is-flex" style="float: left">
          <b-icon
            v-if="!completed"
            icon="checkbox-blank-outline"
            type="is-warning"
          />
          <b-icon
            v-else-if="completed"
            icon="checkbox-outline"
            type="is-success"
          />
          <b v-if="hasInitialLoaded">{{ name }}</b>
          <b-skeleton v-else size="is-small" width="35%" :animated="true" />
          <span class="ml-1 mr-1">
            ${{ currentPrice }}/${{ totalPrice }} ({{ totalPercentage }}%)
          </span>
          <span
            class="display-is-editable pointer-cursor-on-hover"
            @click="PatchShoppingListCompleted(id, !completed)"
          />
          <b-button icon-right="magnify" size="is-small" @click="FocusSearch" />
        </p>
      </div>
    </div>
    <div class="container">
      <section class="section">
        <breadcrumb
          back-link-name="Shopping list"
          :current-page-name="$route.name"
          :current-page-name-override="name"
        />
        <h1
          v-if=" hasInitialLoaded "
          id="ListName"
          class="title is-1 is-marginless display-is-editable pointer-cursor-on-hover"
          @click="ActivateEditListModal"
        >
          {{ name }}
        </h1>
        <b-skeleton v-else size="is-medium" width="35%" :animated="true" />
        <div v-if="notes !== '' || notesFromEmpty || editingMeta">
          <div class="content">
            <label class="label">Notes</label>
            <p
              class="display-is-editable subtitle is-4 pointer-cursor-on-hover notes-highlight"
              @click="ActivateEditListModal"
            >
              <i v-if="hasInitialLoaded"> {{ notes }} </i>
              <b-skeleton
                v-else
                size="is-medium"
                width="35%"
                :animated="true"
              />
            </p>
          </div>
        </div>
        <b-button
          v-else
          type="is-text"
          icon-left="pen"
          class="remove-shadow"
          @click="ActivateEditListModal"
        >
          Add notes
        </b-button>
        <div class="mt-3 mb-3">
          <b-tag v-if="completed" type="is-success"> Completed </b-tag>
          <b-tag v-else-if="!completed" type="is-warning"> Uncompleted </b-tag>
        </div>
        <b-tabs
          v-model="itemDisplayState"
          :position="deviceIsMobile ? 'is-centered' : ''"
          class="m-0"
          :expanded="deviceIsMobile"
        >
          <b-tab-item icon="format-list-checks" label="All" />
          <b-tab-item icon="playlist-remove" label="Unobtained" />
          <b-tab-item icon="playlist-check" label="Obtained" />
        </b-tabs>
        <label class="label">Search for items</label>
        <b-field>
          <b-input
            id="search"
            ref="search"
            v-model="itemSearch"
            icon="magnify"
            size="is-medium"
            placeholder="Item name"
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
              <option value="tags">Tags</option>
              <option value="highestPrice">Highest Price</option>
              <option value="lowestPrice">Lowest Price</option>
              <option value="highestQuantity">Highest Quantity</option>
              <option value="lowestQuantity">Lowest Quantity</option>
              <option value="recentlyAdded">Recently Added</option>
              <option value="lastAdded">Last Added</option>
              <option value="recentlyUpdated">Recently Updated</option>
              <option value="lastUpdated">Last Updated</option>
              <option value="alphabeticalDescending">A-z</option>
              <option value="alphabeticalAscending">z-A</option>
            </b-select>
          </p>
        </b-field>
        <div>
          <section>
            <b-modal
              v-model="isEditListModalActive"
              scroll="keep"
              :fullscreen="deviceIsMobile"
              has-modal-card
              :can-cancel="false"
            >
              <shoppinglistEdit
                v-bind="editListProps"
                @close="isEditListModalActive = false"
              />
            </b-modal>
            <b-modal
              v-model="isNewItemModalActive"
              scroll="keep"
              :fullscreen="deviceIsMobile"
              has-modal-card
              :can-cancel="false"
            >
              <shoppinglistItemNew
                v-bind="newItemProps"
                @close="isNewItemModalActive = false"
              />
            </b-modal>
            <b-modal
              v-model="isEditItemModalActive"
              scroll="keep"
              :fullscreen="deviceIsMobile"
              has-modal-card
              :can-cancel="false"
            >
              <shoppinglistItemEdit
                v-bind="editItemProps"
                @close="isEditItemModalActive = false"
              />
            </b-modal>
            <b-modal
              v-model="isEditTagNameModalActive"
              scroll="keep"
              :fullscreen="deviceIsMobile"
              has-modal-card
              :can-cancel="false"
            >
              <shoppingListTagNameEdit
                v-bind="editTagNameProps"
                @close="CloseEditTagNameModal"
              />
            </b-modal>
          </section>
          <section>
            <div
              class="card pointer-cursor-on-hover"
              @click="ActivateNewItemModal()"
            >
              <div class="card-content">
                <div class="media">
                  <div class="media-left">
                    <b-icon icon="plus-box" size="is-medium" />
                  </div>
                  <div class="media-content">
                    <p class="subtitle is-4">Add a new item</p>
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
        <br />
        <div v-if="listItemsFromTags.length > 0">
          <b-loading
            v-model:active="listIsLoading"
            :is-full-page="false"
            :can-cancel="false"
          />
          <div v-if="sortBy === 'tags'">
            <section
              v-for="itemTag in listItemsFromTags"
              :key="itemTag"
              class="mb-5"
            >
              <div class="field mb-0">
                <p class="control">
                  <b-button
                    type="is-text"
                    size="medium"
                    class="title is-5 is-marginless display-is-editable pointer-cursor-on-hover is-paddingless remove-shadow max-width-80"
                    @click="ActivateEditTagNameModal(itemTag.tag)
                          "
                  >
                    {{ itemTag.tag }}
                  </b-button>
                  <b-button
                    style="float: right; display: block"
                    icon-left="plus-box"
                    size="medium"
                    @click="ActivateNewItemModal(itemTag.tag)"
                  />
                </p>
              </div>
              <transition-group name="staggered-fade" tag="div" :css="false">
                <div v-for="(item, index) in itemTag.items" :key="item">
                  <a :id="item.id" />
                  <itemCard
                    :id="item.id"
                    :list="list"
                    :item="item"
                    :index="index"
                    :list-id="id"
                    :device-is-mobile="deviceIsMobile"
                    :item-display-state="itemDisplayState"
                    @view-item="
                      () => {
                        ActivateEditItemModal(item.id);
                      }
                    "
                    @list="
                      (l) => {
                        itemTag.items = l;
                      }
                    "
                    @obtained="
                      (o) => {
                        item.obtained = o;
                      }
                    "
                  />
                </div>
                <br />
              </transition-group>
              <section>
                <div>
                  <p>
                    {{ itemTag.items.length || 0 }} item(s)
                    <span
                      v-if="
                        itemTag.price !== 0 &&
                          typeof itemTag.price !== 'undefined'
                      "
                    >
                      - ${{ itemTag.price.toFixed(2) }}
                      <span v-if="TagIsExcluded(itemTag.tag)">
                        <b-tag type="is-primary">price excluded</b-tag>
                        <infotooltip
                          v-if="participatingFlatmates.length > 1 || manualSplit > 1"
                          :message="
                                  'Split price plus tag price is $' +
                                  (equalPricePerPerson + itemTag.price).toFixed(2)
                                  "
                          @open="LoopStop"
                          @close="LoopStart"
                        />
                      </span>
                    </span>
                  </p>
                </div>
              </section>
            </section>
          </div>
          <div v-else-if="sortBy !== 'tag'">
            <div v-for="(item, index) in listItemsFromPlainList" :key="item">
              <a :id="item.id" />
              <itemCard
                :id="item.id"
                :list="list"
                :item="item"
                :index="index"
                :list-id="id"
                :display-tag="true"
                :device-is-mobile="deviceIsMobile"
                :item-display-state="itemDisplayState"
                @view-item="
                  () => {
                    ActivateEditItemModal(item.id);
                  }
                "
                @list="
                  (l) => {
                    listItemsFromPlainList = l;
                  }
                "
                @obtained="
                  (o) => {
                    item.obtained = o;
                  }
                "
              />
            </div>
            <section>
              <br />
              <p>{{ listItemsFromPlainList.length || 0 }} item(s)</p>
            </section>
            <br />
          </div>
        </div>
        <div v-else>
          <div class="card m-0">
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
                      itemSearch === '' && list.length === 0 && hasInitialLoaded
                    "
                    class="subtitle is-4"
                  >
                    No items added yet.
                  </p>
                  <p
                    v-else-if="
                      itemSearch === '' &&
                        itemDisplayState === 1 &&
                        list.length > 0 &&
                        hasInitialLoaded
                    "
                    class="subtitle is-4"
                  >
                    All items have been obtained.
                  </p>
                  <p
                    v-else-if="
                      itemSearch === '' &&
                        itemDisplayState === 2 &&
                        list.length > 0 &&
                        hasInitialLoaded
                    "
                    class="subtitle is-4"
                  >
                    No items have been obtained yet.
                  </p>
                  <p
                    v-else-if="itemSearch !== '' && hasInitialLoaded"
                    class="subtitle is-4"
                  >
                    No items found.
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
          <br />
        </div>
        <floatingAddButton
          v-if="!(isNewItemModalActive || isEditItemModalActive || isEditTagNameModalActive || isEditListModalActive)"
          :func="ActivateNewItemModal"
        />
        <p class="subtitle is-4">
          <b>Total items</b>: {{ obtainedCount }}/{{ totalItems }}
          <br />
          <b>Total price</b>: ${{ currentPrice }}/${{ totalPrice }} ({{
          totalPercentage }}%)
          <br />
          <span v-if="totalTagExcludeList.length > 0">
            <b>All inclusive price</b>: ${{ totalAllInclusivePrice }}
            <br />
          </span>
          <span v-if="participatingFlatmates.length > 1 || manualSplit > 0">
            <b>Split price</b>: ${{ equalPricePerPerson.toFixed(2) }}
            <infotooltip
              :message="
                manualSplit === 0
                  ? 'Split is the divided price between ' +
                    participatingFlatmates.length +
                    ' flatmates'
                  : 'Split manually by ' + manualSplit
              "
              @open="LoopStop"
              @close="LoopStart"
            />
          </span>
        </p>
        <b-field>
          <b-button
            :icon-left="
              completed === false
                ? 'checkbox-blank-outline'
                : 'checkbox-outline'
            "
            :type="completed === true ? 'is-success' : 'is-warning'"
            size="is-medium"
            expanded
            @click="PatchShoppingListCompleted(id, !completed)"
          >
            {{ completed === false ? "Uncompleted" : "Completed" }}
          </b-button>
          <p class="control">
            <b-button
              icon-left="delete"
              type="is-danger"
              size="is-medium"
              :loading="deleteLoading"
              @click="DeleteShoppingList(id)"
            />
          </p>
        </b-field>
        <p class="subtitle is-6">
          Created
          <span v-if="hasInitialLoaded || typeof authorName !== 'undefined'"
            >{{ TimestampToCalendar(creationTimestamp) }}, by
          </span>
          <b-skeleton v-else size="is-small" width="35%" :animated="true" />
          <router-link
            v-if="hasInitialLoaded || typeof authorName !== 'undefined'"
            tag="a"
            :to="{ name: 'My Flatmates', query: { id: author } }"
          >
            {{ authorNames }}
          </router-link>
          <b-skeleton v-else size="is-small" width="35%" :animated="true" />
          <span v-if="templateId">
            (templated from
            <router-link
              tag="a"
              :to="{ name: 'View shopping list', params: { id: templateId } }"
              >{{ templateListName }}</router-link
            >)
          </span>
          <span v-if="creationTimestamp !== modificationTimestamp">
            <br />
            Last updated
            <span v-if="hasInitialLoaded || typeof authorName !== 'undefined'"
              >{{ TimestampToCalendar(modificationTimestamp) }}, by
            </span>
            <b-skeleton v-else size="is-small" width="35%" :animated="true" />
            <router-link
              v-if="hasInitialLoaded || typeof authorName !== 'undefined'"
              tag="a"
              :to="{ name: 'My Flatmates', query: { id: authorLast } }"
            >
              {{ authorLastNames }}
            </router-link>
            <b-skeleton v-else size="is-small" width="35%" :animated="true" />
          </span>
        </p>
        <b-collapse
          v-model="shoppingListSettingsOpen"
          class="card"
          animation="slide"
          aria-id="contentIdForA11y3"
        >
          <template #trigger="props">
            <div
              slot="trigger"
              slot-scope="props"
              class="card-header"
              role="button"
              aria-controls="contentIdForA11y3"
            >
              <p class="card-header-title">List settings</p>
              <a class="card-header-icon">
                <b-icon :icon="props.open ? 'menu-down' : 'menu-up'" />
              </a>
            </div>
          </template>
          <div class="card-content">
            <div class="content">
              <h2>Exclude tags in price total</h2>
              <h3>Tags in all lists</h3>
              <div>
                <b-checkbox
                  v-for="existingTag in tags"
                  v-model="totalTagExcludeList"
                  :native-value="existingTag"
                  size="is-medium"
                  @change="UpdateShoppingList"
                >
                  {{ existingTag }}
                </b-checkbox>
                <p v-if="tags.length === 0">No tags found</p>
              </div>

              <h3>Tags in this list</h3>
              <div>
                <b-checkbox
                  v-for="existingListTag in tagsList"
                  v-model="totalTagExcludeList"
                  :native-value="existingListTag"
                  size="is-medium"
                  @change="UpdateShoppingList"
                >
                  {{ existingListTag }}
                </b-checkbox>
                <p v-if="tagsList.length === 0">No tags found</p>
              </div>

              <h3>Manual split</h3>
              <b-field
                label="Temporarily use a custom split instead of active flatmates for this list. Set to zero to disable"
              >
                <b-numberinput
                  v-model="manualSplit"
                  size="is-medium"
                  placeholder="Enter how many of this item should be obtained"
                  min="0"
                  expanded
                  required
                  controls-position="compact"
                  icon="numeric"
                />
              </b-field>
            </div>
          </div>
        </b-collapse>
        <br />
      </section>
    </div>
  </div>
</template>

<script>
  import common from "@/common/common";
  import shoppinglistCommon from "@/common/shoppinglist";
  import shoppinglist from "@/requests/authenticated/shoppinglist";
  import flatmates from "@/requests/authenticated/flatmates";
  import itemCard from "@/components/authenticated/shopping-list-item-card-view.vue";
  import floatingAddButton from "@/components/common/floating-add-button.vue";
  import shoppinglistEdit from "@/components/authenticated/shopping-list-edit.vue";
  import shoppinglistItemNew from "@/components/authenticated/shopping-list-item-new.vue";
  import shoppinglistItemEdit from "@/components/authenticated/shopping-list-item-edit.vue";
  import shoppingListTagNameEdit from "@/components/authenticated/shopping-list-tag-name-edit.vue";
  import infotooltip from "@/components/common/info-tooltip.vue";
  import breadcrumb from "@/components/common/breadcrumb.vue";
  import websockets from "@/requests/authenticated/websockets";

  export default {
    name: "ShoppingList",
    components: {
      itemCard,
      floatingAddButton,
      shoppinglistEdit,
      shoppinglistItemNew,
      shoppinglistItemEdit,
      shoppingListTagNameEdit,
      infotooltip,
      breadcrumb,
    },
    data() {
      return {
        intervalLoop: null,
        editing: false,
        editingMeta: false,
        notesFromEmpty: false,
        itemSearch: shoppinglistCommon.GetShoppingListSearch(this.id) || "",
        authorNames: "",
        authorLastNames: "",
        totalItems: 0,
        loopCreated: new Date(),
        sortBy: shoppinglistCommon.GetShoppingListSortBy() || "tags",
        itemDisplayState:
          shoppinglistCommon.GetShoppingListObtainedFilter(this.id) || 0,
        deviceIsMobile: false,
        HeaderIsSticky: false,
        TagTmp: "",
        editingTag: "",
        listIsLoading:
          shoppinglistCommon.GetShoppingListFromCache(this.id).length > 0,
        hasInitialLoaded: false,
        deleteLoading: false,
        ratherSmallerScreen: false,
        templateListName: "",
        canAnimate: false,
        id: this.$route.params.id,
        name: "Unnamed list",
        notes: "",
        author: "",
        authorLast: "",
        completed: false,
        creationTimestamp: 0,
        modificationTimestamp: 0,
        templateId: undefined,
        listFull: [],
        shoppingListSettingsOpen: false,
        totalTagExcludeList: [],
        tags: [],
        tagsList: [],
        flatmates: [],
        manualSplit: 0,
        websocket: null,
        websocketReady: false,
        isEditListModalActive: false,
        isNewItemModalActive: false,
        isEditItemModalActive: false,
        isEditTagNameModalActive: false,
        editListProps: {
          shoppingListId: "",
          existingName: "",
          existingNotes: "",
          completed: false,
          totalTagExcludeList: [],
        },
        newItemProps: {
          withName: "",
          withTag: "",
        },
        editItemProps: {
          listId: "",
          id: "",
        },
        editTagNameProps: {
          shoppingListId: "",
          shoppingListName: "",
          existingName: "",
        },
      };
    },
    computed: {
      ItemId() {
        return this.$route.query.itemId;
      },
    list() {
      var obtained;
      switch (this.itemDisplayState) {
        case 1:
          obtained = false;
          break;
        case 2:
          obtained = true;
          break;
      }
      return this.listFull.filter(
        (item) => item.obtained === obtained || typeof obtained === "undefined"
      );
      },
      listItemsFromTags() {
        return this.RestructureShoppingListToTags(
          this.list.filter((item) => {
            return this.ItemByNameInList(item);
          })
        );
      },
      listItemsFromPlainList() {
        return this.list.filter((item) => {
          return this.ItemByNameInList(item);
        });
      },
      obtainedCount() {
        if (this.listFull.length === 0) {
          return 0;
        }
        var obtained = 0;
        this.listFull.forEach((item) => {
          obtained += item.obtained === true ? 1 : 0;
        });
        return obtained;
      },
      currentPrice() {
        if (this.listFull.length === 0) {
          return 0;
        }
        var currentPrice = 0;
        this.listFull.forEach((item) => {
          if (
            item.obtained !== true ||
            this.totalTagExcludeList.includes(item.tag)
          ) {
            return;
          }
          if (typeof item.price !== "number") {
            item.price = 0;
          }
          currentPrice += (item.price || 0) * item.quantity;
        });
        currentPrice = currentPrice.toFixed(2);
        return currentPrice;
      },
      totalPrice() {
        if (this.listFull.length === 0) {
          return 0;
        }
        var totalPrice = 0;
        this.listFull.forEach((item) => {
          if (this.totalTagExcludeList.includes(item.tag)) {
            return;
          }
          if (typeof item.price !== "number") {
            item.price = 0;
          }
          totalPrice += (item.price || 0) * item.quantity;
        });
        totalPrice = totalPrice.toFixed(2);
        return totalPrice;
      },
      totalAllInclusivePrice() {
        if (this.listFull.length === 0) {
          return 0;
        }
        var totalPrice = 0;
        this.listFull.forEach((item) => {
          if (typeof item.price !== "number") {
            item.price = 0;
          }
          totalPrice += (item.price || 0) * item.quantity;
        });
        totalPrice = totalPrice.toFixed(2);
        return totalPrice;
      },
      participatingFlatmates() {
        return this.flatmates.filter(
          (u) => u.disabled !== true && u.registered === true
        );
      },
      equalPricePerPerson() {
        if (this.manualSplit !== 0) {
          return this.totalPrice / this.manualSplit;
        }
        return this.totalPrice / this.participatingFlatmates.length;
      },
      totalPercentage() {
        return Math.round((100 * this.currentPrice) / this.totalPrice) || 0;
      },
    },
    watch: {
      id() {
        console.log("ID changed");
        window.location.reload(false);
      },
      sortBy() {
        shoppinglistCommon.WriteShoppingListSortBy(this.sortBy);
        this.listIsLoading = true;
        this.ResetLoopTime();
        this.LoopStop();
        this.LoopStart();
      },
      itemDisplayState() {
        this.listIsLoading = true;
        this.GetShoppingListItems();
        shoppinglistCommon.WriteShoppingListObtainedFilter(
          this.id,
          this.itemDisplayState
        );
      },
      itemSearch() {
        shoppinglistCommon.WriteShoppingListSearch(this.id, this.itemSearch);
      },
      hasInitialLoaded() {
        this.canAnimate = true;
      },
      completed() {
        var enableAnimations = common.GetEnableAnimations();
        if (
          this.completed === true &&
          enableAnimations !== "false" &&
          this.canAnimate === true
        ) {
          common.Hooray();
        }
      },
      author() {
        if (this.authorNames !== "") {
          return;
        }
        flatmates.GetFlatmate(this.author).then((resp) => {
          this.authorNames = resp.data.spec.names;
        });
      },
      authorLast() {
        if (this.authorLastNames !== "") {
          return;
        }
        flatmates.GetFlatmate(this.authorLast).then((resp) => {
          this.authorLastNames = resp.data.spec.names;
        });
      },
      templateId() {
        if (typeof this.templateId === "undefined" || this.templateId === "") {
          return;
        }
        shoppinglist.GetShoppingList(this.templateId).then((resp) => {
          this.templateListName = resp.data.spec.name;
        });
      },
      isNewItemModalActive() {
        if (this.isNewItemModalActive !== false) {
          return;
        }
        this.GetShoppingListItems();
      },
      isEditItemModalActive() {
        if (this.isEditItemModalActive !== false) {
          return;
        }
        this.GetShoppingListItems();
      },
    websocket() {
      if (this.websocket === null) {
        this.OpenWebSocket();
      }
    },
    },
    async beforeMount() {
      this.GetShoppingList();
      this.GetShoppingListItems();
      this.GetFlatmates();
      if (window.innerWidth <= 330) {
        this.ratherSmallerScreen = true;
      }
    },
    async created() {
      this.CheckDeviceIsMobile();
      window.addEventListener("resize", this.CheckDeviceIsMobile, true);
      window.addEventListener("scroll", this.ManageStickyHeader, true);
      this.LoopStart();
      window.addEventListener("focus", this.ResetLoopTime, true);
      // TODO better way to do this? why does this not pull in through the data state function?
      this.itemDisplayState = shoppinglistCommon.GetShoppingListObtainedFilter(
        this.id
      );
    this.OpenWebSocket();
    this.websocket.onopen = (evt) => {
      this.websocketReady = true;
      this.websocket.send("session opened.");
      console.log("WS ready!");
    };
    this.websocket.onmessage = (evt) => {
      console.log("msg", { evt });
      this.GetShoppingListItemsFromWS(evt);
    };
    this.websocket.onclose = (evt) => {
      console.log("WebSocket closed.");
      this.websocket.send("session closed.");
      setTimeout(this.OpenWebSocket, 100);
    };
    this.websocket.onerror = (evt) => {
      console.log("err", { evt });
    };
    },
    mounted() {
      if (typeof this.ItemId !== "undefined") {
        var el = this.$refs[this.ItemId][0].$el;
        window.scrollTo(0, el.offsetTop);
      }
    },
    beforeUnmount() {
      this.LoopStop();
      window.removeEventListener("resize", this.CheckDeviceIsMobile, true);
      window.removeEventListener("scroll", this.ManageStickyHeader, true);
      window.removeEventListener("focus", this.ResetLoopTime, true);
    },
    methods: {
      CopyHrefToClipboard() {
        common.CopyHrefToClipboard();
      },
      ActivateEditListModal() {
        this.editListProps = {
          shoppingListId: this.id,
          existingName: this.name,
          existingNotes: this.notes,
          completed: this.completed,
          totalTagExcludeList: this.totalTagExcludeList,
        };
        this.isEditListModalActive = true;
      },
      ActivateNewItemModal(tag) {
        this.newItemProps = {
          withName: this.itemSearch,
          withTag: tag || "",
        };
        this.isNewItemModalActive = true;
        this.itemSearch = "";
      },
      ActivateEditItemModal(itemId) {
        this.editItemProps = {
          shoppingListId: this.id,
          id: itemId,
        };
        this.isEditItemModalActive = true;
      },
      ActivateEditTagNameModal(tag) {
        this.editTagNameProps = {
          shoppingListName: this.name,
          shoppingListId: this.id,
          existingName: tag,
        };
        this.isEditTagNameModalActive = true;
      },
      CloseEditTagNameModal() {
        this.isEditTagNameModalActive = false;
        this.GetShoppingListItems();
      },
      ItemByNameInList(item) {
        var vm = this;
        return (
          item.name.toLowerCase().indexOf(vm.itemSearch.toLowerCase()) !== -1
        );
      },
      FocusSearchBox() {
        this.$refs.search.$el.focus();
      },
      RestructureShoppingListToTags(list) {
        return shoppinglistCommon.RestructureShoppingListToTags(list);
      },
      GetShoppingList() {
        if (this.editing === true) {
          return;
        }
        var id = this.id;
        shoppinglist
          .GetShoppingList(id)
          .then((resp) => {
            this.name = resp.data.spec.name;
            this.notes = resp.data.spec.notes || "";
            this.author = resp.data.spec.author;
            this.authorLast = resp.data.spec.authorLast;
            this.completed = resp.data.spec.completed;
            this.creationTimestamp = resp.data.spec.creationTimestamp;
            this.modificationTimestamp = resp.data.spec.modificationTimestamp;
            this.templateId = resp.data.spec.templateId;
            this.totalTagExcludeList = resp.data.spec.totalTagExclude || [];
          })
          .catch((err) => {
            if (err.response.status === 404) {
              common.DisplayFailureToast(
                this.$buefy,
                "Error list not found" +
                  "<br/>" +
                  err.response.data.metadata.response
              );
              this.$router.push({ name: "Shopping list" });
              return;
            }
            common.DisplayFailureToast(
              this.$buefy,
              "Error loading the shopping list" +
                "<br/>" +
                err.response.data.metadata.response
            );
          });

        shoppinglist.GetAllShoppingListItemTags().then((resp) => {
          this.itemIsLoading = false;
          if (resp.data.list === null) {
            return;
          }
          this.tags = resp.data.list.map((i) => i.name) || [];
        });
        shoppinglist.GetShoppingListItemTags(this.id).then((resp) => {
          this.tagsList = resp.data.list || [];
        });
      },
      UpdateShoppingList() {
        this.notesFromEmpty = false;
        this.editingMeta = false;
        this.editing = false;

        var id = this.id;
        shoppinglist
          .UpdateShoppingList(
            id,
            this.name,
            this.notes,
            this.completed,
            this.totalTagExcludeList
          )
          .catch((err) => {
            common.DisplayFailureToast(
              this.$buefy,
              "Failed to update shopping list" +
                "<br/>" +
                err.response.data.metadata.response
            );
          });
      },
      PatchShoppingListCompleted(id, completed) {
        shoppinglist
          .PatchShoppingListCompleted(id, completed)
          .then((resp) => {
            this.completed = resp.data.spec.completed;
          })
          .catch((err) => {
            common.DisplayFailureToast(
              this.$buefy,
              "Failed to set list as completed" +
                "<br/>" +
                err.response.data.metadata.response
            );
          });
      },
      DeleteShoppingList(id) {
        this.$buefy.dialog.confirm({
          title: "Delete shopping list",
          message:
            "Are you sure that you wish to delete this shopping list?" +
            "<br/>" +
            "This action cannot be undone.",
          confirmText: "Delete shopping list",
          type: "is-danger",
          hasIcon: true,
          onConfirm: () => {
            this.deleteLoading = true;
            window.clearInterval(this.intervalLoop);
            shoppinglist
              .DeleteShoppingList(id)
              .then((resp) => {
                common.DisplaySuccessToast(
                  this.$buefy,
                  "Deleted the shopping list"
                );
                shoppinglistCommon.DeleteShoppingListFromCache(id);
                setTimeout(() => {
                  this.$router.push({ name: "Shopping list" });
                }, 1 * 1000);
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
      GetShoppingListItems() {
        var obtained;
        switch (this.itemDisplayState) {
          case 1:
            obtained = false;
            break;
          case 2:
            obtained = true;
            break;
        }

        shoppinglist
          .GetShoppingListItems(this.id, this.sortBy, undefined)
          .then((resp) => {
            var responseList = resp.data.list || [];
            this.totalItems = responseList === null ? 0 : responseList.length;
            if (this.list === null) {
              this.list = [];
            }

            if (responseList !== this.list) {
              this.listFull = responseList;
              this.list = responseList.filter(
                (item) =>
                  item.obtained === obtained || typeof obtained === "undefined"
              );
              shoppinglistCommon.WriteShoppingListToCache(this.id, this.list);
              this.listIsLoading = false;
              this.hasInitialLoaded = true;
            }
          });
      },
      TimestampToCalendar(timestamp) {
        return common.TimestampToCalendar(timestamp);
      },
      LoopStart() {
        if (shoppinglistCommon.GetShoppingListAutoRefresh() === "false") {
          return;
        }
        this.intervalLoop = window.setInterval(() => {
          if (this.editing === true) {
            return;
          }
          this.GetShoppingList();
          this.GetShoppingListItems();

          var now = new Date();
          var timePassed =
            now.getTime() / 1000 - this.loopCreated.getTime() / 1000;
          if (timePassed >= 3600 / 4) {
            window.clearInterval(this.intervalLoop);
          }
        }, 3 * 1000);
      },
      LoopStop() {
        window.clearInterval(this.intervalLoop);
      },
      CheckDeviceIsMobile() {
        this.deviceIsMobile = common.DeviceIsMobile();
      },
      ManageStickyHeader() {
        this.HeaderIsSticky =
          window.pageYOffset >
          document.getElementById("ListName").offsetTop + 30;
      },
      ResetLoopTime() {
        this.loopCreated = new Date();
      },
      FocusEl(name) {
        this.$nextTick(() => {
          this.$refs[name].focus();
        });
      },
      FocusName() {
        this.FocusEl("name");
      },
      FocusNotes() {
        this.FocusEl("notes");
      },
      FocusSearch() {
        this.FocusEl("search");
      },
      TagIsExcluded(tag) {
        return this.totalTagExcludeList.includes(tag);
      },
      GetFlatmates() {
        flatmates.GetAllFlatmates().then((resp) => {
          if (resp.data.list === null) {
            this.flatmates = [];
            return;
          }
          this.flatmates = resp.data.list;
        });
      },
    OpenWebSocket() {
      console.log("Opening WS conn...");
      this.websocket = websockets.GetWebSocket(
        "shoppingitem",
        this.id,
        { name: "sortBy", value: this.sortBy },
        { name: "obtained", value: this.obtained }
      );
    },
    SendWebSocketListUpdated() {
      this.websocket.send("updated");
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

  .ListBar {
    position: fixed;
    height: auto;
    width: 100%;
    z-index: 20;
    padding: 10px;
    box-shadow: black 0px -45px 71px;
    display: block;
    background-color: hsla(0, 0%, 100%, 0.73);
    backdrop-filter: blur(5px);
  }

  .ListBarTop {
    top: 0;
  }
  @media (prefers-color-scheme: dark) {
    .ListBar {
      background-color: #333333ba;
    }
  }
  [data-theme="dark"] {
    .ListBar {
      background-color: #333333ba;
    }
  }
  @media (prefers-color-scheme: light) {
    .ListBar {
      background-color: hsla(0, 0%, 100%, 0.73);
    }
  }
  [data-theme="light"] {
    .ListBar {
      background-color: hsla(0, 0%, 100%, 0.73);
    }
  }
  .max-width-80 {
    max-width: 80%;
  }
</style>
