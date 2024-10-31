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
      v-if="HeaderIsSticky && !editing"
      :class="ratherSmallerScreen ? 'ListBar ListBarTop' : 'ListBar'"
    >
      <div class="level">
        <p class="subtitle is-5 mb-0" style="float: left">
          <b>{{ name }}</b>
          ${{ currentPrice }}/${{ totalPrice }} ({{ totalPercentage }}%)
          <span
            @click="PatchShoppingListCompleted(id, !completed)"
            class="display-is-editable pointer-cursor-on-hover"
          >
            <b-tag type="is-warning" v-if="!completed">Uncompleted</b-tag>
            <b-tag type="is-success" v-else-if="completed">Completed</b-tag>
          </span>
        </p>
        <div style="float: right">
          <b-button
            icon-right="magnify"
            size="is-small"
            @click="FocusSearch"
          ></b-button>
        </div>
      </div>
    </div>
    <div class="container">
      <section class="section">
        <nav
          class="breadcrumb is-medium has-arrow-separator"
          aria-label="breadcrumbs"
        >
          <ul>
            <li>
              <router-link :to="{ name: 'Shopping list' }"
                >Shopping list</router-link
              >
            </li>
            <li
              v-if="
                hasInitialLoaded || name !== '' || typeof name !== 'undefined'
              "
              class="is-active"
            >
              <router-link
                :to="{ name: 'View shopping list', params: { id: id } }"
                >{{ name || "Unnamed list" }}</router-link
              >
            </li>
            <b-skeleton
              v-else
              size="is-small"
              width="35%"
              :animated="true"
            ></b-skeleton>
          </ul>
          <b-button
            @click="CopyHrefToClipboard()"
            icon-left="content-copy"
            size="is-small"
          ></b-button>
        </nav>
        <div v-if="editingMeta">
          <b-field label="Name">
            <b-input
              type="text"
              icon="format-title"
              size="is-medium"
              ref="name"
              placeholder="Enter a title for this list"
              @keyup.enter.native="UpdateShoppingList"
              @keyup.esc.native="
                editing = false;
                editingMeta = false;
              "
              icon-right="close-circle"
              icon-right-clickable
              @icon-right-click="name = ''"
              v-model="name"
              required
            >
            </b-input>
          </b-field>
          <br />
        </div>
        <div v-else>
          <h1
            v-if="
              hasInitialLoaded || name !== '' || typeof name !== 'undefined'
            "
            id="ListName"
            class="title is-1 is-marginless display-is-editable pointer-cursor-on-hover"
            @click="
              editing = true;
              editingMeta = true;
              FocusName();
            "
          >
            {{ name }}
          </h1>
          <b-skeleton
            v-else
            size="is-medium"
            width="35%"
            :animated="true"
          ></b-skeleton>
        </div>
        <div v-if="notes !== '' || notesFromEmpty || editingMeta">
          <div v-if="editingMeta">
            <b-field label="Notes">
              <b-input
                icon="text"
                size="is-medium"
                maxlength="100"
                type="text"
                ref="notes"
                placeholder="Enter extra information"
                @keyup.enter.native="UpdateShoppingList"
                @keyup.esc.native="
                  editing = false;
                  editingMeta = false;
                "
                icon-right="close-circle"
                icon-right-clickable
                @icon-right-click="notes = ''"
                v-model="notes"
              >
              </b-input>
            </b-field>
          </div>
          <div v-else>
            <br />
            <div>
              <div class="content">
                <label class="label">Notes</label>
                <p
                  class="display-is-editable subtitle is-4 pointer-cursor-on-hover notes-highlight"
                  @click="
                    editingMeta = true;
                    editing = true;
                    FocusNotes();
                  "
                >
                  <i>
                    {{ notes }}
                  </i>
                </p>
              </div>
            </div>
          </div>
        </div>
        <b-button
          type="is-text"
          @click="
            () => {
              notesFromEmpty = true;
              editingMeta = true;
              editing = true;
              FocusNotes();
            }
          "
          v-if="!editingMeta && notes.length == 0"
          class="remove-shadow"
          >Add notes</b-button
        >
        <div v-if="editingMeta">
          <b-button type="is-info" @click="UpdateShoppingList()">Done</b-button>
          <br />
        </div>
        <br />
        <b-tag type="is-success" v-if="completed">Completed</b-tag>
        <b-tag type="is-warning" v-else-if="!completed">Uncompleted</b-tag>
        <br />
        <b-tabs
          :position="deviceIsMobile ? 'is-centered' : ''"
          class="block is-marginless"
          v-model="itemDisplayState"
          :expanded="deviceIsMobile"
        >
          <b-tab-item icon="format-list-checks" label="All"></b-tab-item>
          <b-tab-item icon="playlist-remove" label="Unobtained"></b-tab-item>
          <b-tab-item icon="playlist-check" label="Obtained"></b-tab-item>
        </b-tabs>
        <label class="label">Search for items</label>
        <b-field>
          <b-input
            id="search"
            icon="magnify"
            size="is-medium"
            placeholder="Item name"
            type="search"
            v-model="itemSearch"
            ref="search"
            expanded
          >
          </b-input>
          <p class="control">
            <b-select
              placeholder="Sort by"
              icon="sort"
              v-model="sortBy"
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
              scroll="keep"
              v-model="isNewItemModalActive"
              :fullscreen="deviceIsMobile"
              has-modal-card
              :can-cancel="false"
            >
              <shoppinglistItemNew v-bind="newItemProps"></shoppinglistItemNew>
            </b-modal>
            <b-modal
              scroll="keep"
              v-model="isEditItemModalActive"
              :fullscreen="deviceIsMobile"
              has-modal-card
              :can-cancel="false"
            >
              <shoppinglistItemEdit
                v-bind="editItemProps"
              ></shoppinglistItemEdit>
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
                    <b-icon icon="plus-box" size="is-medium"> </b-icon>
                  </div>
                  <div class="media-content">
                    <p class="subtitle is-4">Add a new item</p>
                  </div>
                  <div class="media-right">
                    <b-icon
                      icon="chevron-right"
                      size="is-medium"
                      type="is-midgray"
                    ></b-icon>
                  </div>
                </div>
              </div>
            </div>
          </section>
        </div>
        <br />
        <div v-if="listItemsFromTags.length > 0">
          <b-loading
            :is-full-page="false"
            :active.sync="listIsLoading"
            :can-cancel="false"
          ></b-loading>
          <div v-if="sortBy === 'tags'">
            <section v-for="itemTag in listItemsFromTags" v-bind:key="itemTag">
              <div v-if="editingTag === itemTag.tag">
                <label class="label">Tag name</label>
                <b-field>
                  <b-button
                    type="is-danger"
                    icon-right="delete"
                    size="is-medium"
                    @click="
                      DeleteShoppingListTagItems(itemTag.tag);
                      editingTag = '';
                      editing = false;
                    "
                  ></b-button>
                  <b-button
                    type="is-danger"
                    size="is-medium"
                    @click="
                      editingTag = '';
                      editing = false;
                    "
                  >
                    X
                  </b-button>
                  <b-input
                    type="text"
                    icon="format-title"
                    size="is-medium"
                    maxlength="30"
                    placeholder="Enter a tag name"
                    expanded
                    @keyup.enter.native="
                      editingTag = '';
                      UpdateShoppingListItemTag(itemTag.tag, TagTmp);
                      itemTag.tag = TagTmp;
                      TagTmp = '';
                      editing = false;
                    "
                    @keyup.esc.native="
                      editingTag = '';
                      editing = false;
                    "
                    icon-right="close-circle"
                    icon-right-clickable
                    @icon-right-click="TagTmp = ''"
                    v-model="TagTmp"
                    required
                  >
                  </b-input>
                  <p class="control">
                    <b-button
                      type="is-primary"
                      size="is-medium"
                      icon-left="check"
                      @click="
                        editingTag = '';
                        UpdateShoppingListItemTag(itemTag.tag, TagTmp);
                        itemTag.tag = TagTmp;
                        TagTmp = '';
                        editing = false;
                      "
                    >
                    </b-button>
                  </p>
                </b-field>
                <br />
              </div>
              <div v-else>
                <div class="field">
                  <p class="control">
                    <b-button
                      type="is-text"
                      size="medium"
                      class="title is-5 is-marginless display-is-editable pointer-cursor-on-hover is-paddingless remove-shadow"
                      @click="
                        TagTmp = itemTag.tag;
                        editingTag = itemTag.tag;
                        editing = true;
                      "
                    >
                      {{ itemTag.tag }}
                    </b-button>
                    <b-button
                      style="float: right; display: block"
                      icon-left="plus-box"
                      size="medium"
                      @click="ActivateNewItemModal(itemTag.tag)"
                    >
                    </b-button>
                  </p>
                </div>
              </div>
              <transition-group
                name="staggered-fade"
                tag="div"
                v-bind:css="false"
              >
                <div v-for="(item, index) in itemTag.items" v-bind:key="item">
                  <a :id="item.id"></a>
                  <itemCard
                    :list="list"
                    :item="item"
                    :index="index"
                    :listId="id"
                    :deviceIsMobile="deviceIsMobile"
                    :id="item.id"
                    :itemDisplayState="itemDisplayState"
                    @viewItem="
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
                <br />
                <div class="level">
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
                        />
                      </span>
                    </span>
                  </p>
                </div>
              </section>
              <br />
            </section>
          </div>
          <div v-else-if="sortBy !== 'tag'">
            <div
              v-for="(item, index) in listItemsFromPlainList"
              v-bind:key="item"
            >
              <a :id="item.id"></a>
              <itemCard
                :list="list"
                :item="item"
                :index="index"
                :listId="id"
                :displayTag="true"
                :deviceIsMobile="deviceIsMobile"
                :id="item.id"
                :itemDisplayState="itemDisplayState"
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
          <div class="card">
            <div class="card-content card-content-list">
              <div class="media">
                <div class="media-left">
                  <b-icon
                    icon="cart-remove"
                    size="is-medium"
                    type="is-midgray"
                  ></b-icon>
                </div>
                <div class="media-content">
                  <p
                    class="subtitle is-4"
                    v-if="
                      itemSearch === '' && list.length === 0 && !pageLoading
                    "
                  >
                    No items added yet.
                  </p>
                  <p
                    class="subtitle is-4"
                    v-else-if="
                      itemSearch === '' &&
                      itemDisplayState === 1 &&
                      list.length > 0 &&
                      !pageLoading
                    "
                  >
                    All items have been obtained.
                  </p>
                  <p
                    class="subtitle is-4"
                    v-else-if="
                      itemSearch === '' &&
                      itemDisplayState === 2 &&
                      list.length > 0 &&
                      !pageLoading
                    "
                  >
                    No items have been obtained yet.
                  </p>
                  <p
                    class="subtitle is-4"
                    v-else-if="itemSearch !== '' && !pageLoading"
                  >
                    No items found.
                  </p>
                  <p class="subtitle is-4" v-else-if="pageLoading">
                    Loading items...
                  </p>
                </div>
              </div>
            </div>
          </div>
          <br />
        </div>
        <floatingAddButton
          v-if="!(isNewItemModalActive || isEditItemModalActive)"
          :func="ActivateNewItemModal"
        />
        <p class="subtitle is-4">
          <b>Total items</b>: {{ obtainedCount }}/{{ totalItems }}
          <br />
          <b>Total price</b>: ${{ currentPrice }}/${{ totalPrice }} ({{
            totalPercentage
          }}%)
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
            />
          </span>
        </p>
        <b-field>
          <b-button
            :icon-left="
              completed === false
                ? 'checkbox-blank-outline'
                : 'check-box-outline'
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
            >
            </b-button>
          </p>
        </b-field>
        <p class="subtitle is-6">
          Created
          <span v-if="hasInitialLoaded || typeof authorName !== 'undefined'"
            >{{ TimestampToCalendar(creationTimestamp) }}, by</span
          >
          <b-skeleton
            v-else
            size="is-small"
            width="35%"
            :animated="true"
          ></b-skeleton>
          <router-link
            v-if="hasInitialLoaded || typeof authorName !== 'undefined'"
            tag="a"
            :to="{ name: 'My Flatmates', query: { id: author } }"
          >
            {{ authorNames }}
          </router-link>
          <b-skeleton
            v-else
            size="is-small"
            width="35%"
            :animated="true"
          ></b-skeleton>
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
              >{{ TimestampToCalendar(modificationTimestamp) }}, by</span
            >
            <b-skeleton
              v-else
              size="is-small"
              width="35%"
              :animated="true"
            ></b-skeleton>
            <router-link
              v-if="hasInitialLoaded || typeof authorName !== 'undefined'"
              tag="a"
              :to="{ name: 'My Flatmates', query: { id: authorLast } }"
            >
              {{ authorLastNames }}
            </router-link>
            <b-skeleton
              v-else
              size="is-small"
              width="35%"
              :animated="true"
            ></b-skeleton>
          </span>
        </p>
        <b-collapse
          class="card"
          animation="slide"
          aria-id="contentIdForA11y3"
          :open="shoppingListSettingsOpen"
        >
          <div
            slot="trigger"
            slot-scope="props"
            class="card-header"
            role="button"
            aria-controls="contentIdForA11y3"
          >
            <p class="card-header-title">List settings</p>
            <a class="card-header-icon">
              <b-icon :icon="props.open ? 'menu-down' : 'menu-up'"> </b-icon>
            </a>
          </div>
          <div class="card-content">
            <div class="content">
              <h2>Exclude tags in price total</h2>
              <h3>Tags in all lists</h3>
              <div>
                <b-checkbox
                  v-for="existingTag in tags"
                  :key="existingTag"
                  v-model="totalTagExcludeList"
                  @input="UpdateShoppingList"
                  :native-value="existingTag"
                  size="is-medium"
                >
                  {{ existingTag }}
                </b-checkbox>
                <p v-if="tags.length === 0">No tags found</p>
              </div>

              <h3>Tags in this list</h3>
              <div>
                <b-checkbox
                  v-for="existingListTag in tagsList"
                  :key="existingListTag"
                  v-model="totalTagExcludeList"
                  @input="UpdateShoppingList"
                  :native-value="existingListTag"
                  size="is-medium"
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
                >
                </b-numberinput>
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
import common from '@/common/common'
import shoppinglistCommon from '@/common/shoppinglist'
import shoppinglist from '@/requests/authenticated/shoppinglist'
import flatmates from '@/requests/authenticated/flatmates'
import { DialogProgrammatic as Dialog } from 'buefy'

export default {
  name: 'shopping-list',
  data () {
    return {
      intervalLoop: null,
      editing: false,
      editingMeta: false,
      notesFromEmpty: false,
      itemSearch: shoppinglistCommon.GetShoppingListSearch(this.id) || '',
      authorNames: '',
      authorLastNames: '',
      totalItems: 0,
      loopCreated: new Date(),
      sortBy: shoppinglistCommon.GetShoppingListSortBy() || 'tags',
      itemDisplayState:
        shoppinglistCommon.GetShoppingListObtainedFilter(this.id) || 0,
      deviceIsMobile: false,
      HeaderIsSticky: false,
      TagTmp: '',
      editingTag: '',
      listIsLoading:
        shoppinglistCommon.GetShoppingListFromCache(this.id).length > 0,
      hasInitialLoaded: false,
      deleteLoading: false,
      ratherSmallerScreen: false,
      templateListName: '',
      canAnimate: false,
      id: this.$route.params.id,
      name: 'Unnamed list',
      notes: '',
      author: '',
      authorLast: '',
      completed: false,
      creationTimestamp: 0,
      modificationTimestamp: 0,
      templateId: undefined,
      list: shoppinglistCommon.GetShoppingListFromCache(this.id) || [],
      listFull: [],
      shoppingListSettingsOpen: false,
      totalTagExcludeList: [],
      tags: [],
      tagsList: [],
      flatmates: [],
      manualSplit: 0,
      isNewItemModalActive: false,
      isEditItemModalActive: false,
      newItemProps: {
        withName: '',
        withTag: ''
      },
      editItemProps: {
        listId: '',
        id: ''
      }
    }
  },
  components: {
    itemCard: () =>
      import('@/components/authenticated/shopping-list-item-card-view.vue'),
    floatingAddButton: () =>
      import('@/components/common/floating-add-button.vue'),
    shoppinglistItemNew: () =>
      import('@/components/authenticated/shopping-list-item-new.vue'),
    shoppinglistItemEdit: () =>
      import('@/components/authenticated/shopping-list-item-edit.vue'),
    infotooltip: () => import('@/components/common/info-tooltip.vue')
  },
  computed: {
    ItemId () {
      return this.$route.query.itemId
    },
    listItemsFromTags () {
      return this.RestructureShoppingListToTags(
        this.list.filter((item) => {
          return this.ItemByNameInList(item)
        })
      )
    },
    listItemsFromPlainList () {
      return this.list.filter((item) => {
        return this.ItemByNameInList(item)
      })
    },
    obtainedCount () {
      if (this.listFull.length === 0) {
        return 0
      }
      var obtained = 0
      this.listFull.forEach((item) => {
        obtained += item.obtained === true ? 1 : 0
      })
      return obtained
    },
    currentPrice () {
      if (this.listFull.length === 0) {
        return 0
      }
      var currentPrice = 0
      this.listFull.forEach((item) => {
        if (
          item.obtained !== true ||
          this.totalTagExcludeList.includes(item.tag)
        ) {
          return
        }
        if (typeof item.price !== 'number') {
          item.price = 0
        }
        currentPrice += (item.price || 0) * item.quantity
      })
      currentPrice = currentPrice.toFixed(2)
      return currentPrice
    },
    totalPrice () {
      if (this.listFull.length === 0) {
        return 0
      }
      var totalPrice = 0
      this.listFull.forEach((item) => {
        if (this.totalTagExcludeList.includes(item.tag)) {
          return
        }
        if (typeof item.price !== 'number') {
          item.price = 0
        }
        totalPrice += (item.price || 0) * item.quantity
      })
      totalPrice = totalPrice.toFixed(2)
      return totalPrice
    },
    totalAllInclusivePrice () {
      if (this.listFull.length === 0) {
        return 0
      }
      var totalPrice = 0
      this.listFull.forEach((item) => {
        if (typeof item.price !== 'number') {
          item.price = 0
        }
        totalPrice += (item.price || 0) * item.quantity
      })
      totalPrice = totalPrice.toFixed(2)
      return totalPrice
    },
    participatingFlatmates () {
      return this.flatmates.filter(
        (u) => u.disabled !== true && u.registered === true
      )
    },
    equalPricePerPerson () {
      if (this.manualSplit !== 0) {
        return this.totalPrice / this.manualSplit
      }
      return this.totalPrice / this.participatingFlatmates.length
    },
    totalPercentage () {
      return Math.round((100 * this.currentPrice) / this.totalPrice) || 0
    }
  },
  methods: {
    CopyHrefToClipboard () {
      common.CopyHrefToClipboard()
    },
    ActivateNewItemModal (tag) {
      this.newItemProps = {
        withName: this.itemSearch,
        withTag: tag || ''
      }
      this.isNewItemModalActive = true
      this.itemSearch = ''
    },
    ActivateEditItemModal (itemId) {
      this.editItemProps = {
        shoppingListId: this.id,
        id: itemId
      }
      this.isEditItemModalActive = true
    },
    ItemByNameInList (item) {
      var vm = this
      return (
        item.name.toLowerCase().indexOf(vm.itemSearch.toLowerCase()) !== -1
      )
    },
    FocusSearchBox () {
      this.$refs.search.$el.focus()
    },
    RestructureShoppingListToTags (list) {
      return shoppinglistCommon.RestructureShoppingListToTags(list)
    },
    GetShoppingList () {
      if (this.editing === true) {
        return
      }
      var id = this.id
      shoppinglist
        .GetShoppingList(id)
        .then((resp) => {
          this.name = resp.data.spec.name
          this.notes = resp.data.spec.notes || ''
          this.author = resp.data.spec.author
          this.authorLast = resp.data.spec.authorLast
          this.completed = resp.data.spec.completed
          this.creationTimestamp = resp.data.spec.creationTimestamp
          this.modificationTimestamp = resp.data.spec.modificationTimestamp
          this.templateId = resp.data.spec.templateId
          this.totalTagExcludeList = resp.data.spec.totalTagExclude || []
        })
        .catch((err) => {
          if (err.response.status === 404) {
            common.DisplayFailureToast(
              'Error list not found' +
                '<br/>' +
                err.response.data.metadata.response
            )
            this.$router.push({ name: 'Shopping list' })
            return
          }
          common.DisplayFailureToast(
            'Error loading the shopping list' +
              '<br/>' +
              err.response.data.metadata.response
          )
        })

      shoppinglist.GetAllShoppingListItemTags().then((resp) => {
        this.itemIsLoading = false
        if (resp.data.list === null) {
          return
        }
        this.tags = resp.data.list.map((i) => i.name) || []
      })
      shoppinglist.GetShoppingListItemTags(this.id).then((resp) => {
        this.tagsList = resp.data.list || []
      })
    },
    UpdateShoppingList () {
      this.notesFromEmpty = false
      this.editingMeta = false
      this.editing = false

      var id = this.id
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
            'Failed to update shopping list' +
              '<br/>' +
              err.response.data.metadata.response
          )
        })
    },
    PatchShoppingListCompleted (id, completed) {
      shoppinglist
        .PatchShoppingListCompleted(id, completed)
        .then((resp) => {
          this.completed = resp.data.spec.completed
        })
        .catch((err) => {
          common.DisplayFailureToast(
            'Failed to set list as completed' +
              '<br/>' +
              err.response.data.metadata.response
          )
        })
    },
    DeleteShoppingListTagItems (name) {
      Dialog.confirm({
        title: 'Delete shopping list tag items',
        message:
          `Are you sure that you wish to delete '${name}'?` +
          '<br/>' +
          'This action cannot be undone.',
        confirmText: 'Delete tag items',
        type: 'is-danger',
        hasIcon: true,
        onConfirm: () => {
          this.deleteLoading = true
          shoppinglist
            .DeleteShoppingListTagItems(this.id, name)
            .then((resp) => {
              this.deleteLoading = false
              common.DisplaySuccessToast('Deleted the shopping list tag')
              this.GetShoppingListItems()
            })
            .catch((err) => {
              this.deleteLoading = false
              common.DisplayFailureToast(
                'Failed to delete the shopping list' +
                  '<br/>' +
                  err.response.data.metadata.response
              )
            })
        }
      })
    },
    DeleteShoppingList (id) {
      Dialog.confirm({
        title: 'Delete shopping list',
        message:
          'Are you sure that you wish to delete this shopping list?' +
          '<br/>' +
          'This action cannot be undone.',
        confirmText: 'Delete shopping list',
        type: 'is-danger',
        hasIcon: true,
        onConfirm: () => {
          this.deleteLoading = true
          window.clearInterval(this.intervalLoop)
          shoppinglist
            .DeleteShoppingList(id)
            .then((resp) => {
              common.DisplaySuccessToast('Deleted the shopping list')
              shoppinglistCommon.DeleteShoppingListFromCache(id)
              setTimeout(() => {
                this.$router.push({ name: 'Shopping list' })
              }, 1 * 1000)
            })
            .catch((err) => {
              this.deleteLoading = false
              common.DisplayFailureToast(
                'Failed to delete the shopping list' +
                  '<br/>' +
                  err.response.data.metadata.response
              )
            })
        }
      })
    },
    GetShoppingListItems () {
      var obtained
      switch (this.itemDisplayState) {
        case 1:
          obtained = false
          break
        case 2:
          obtained = true
          break
      }

      shoppinglist
        .GetShoppingListItems(this.id, this.sortBy, undefined)
        .then((resp) => {
          var responseList = resp.data.list || []
          this.totalItems = responseList === null ? 0 : responseList.length
          if (this.list === null) {
            this.list = []
          }

          if (responseList !== this.list) {
            this.listFull = responseList
            this.list = responseList.filter(
              (item) =>
                item.obtained === obtained || typeof obtained === 'undefined'
            )
            shoppinglistCommon.WriteShoppingListToCache(this.id, this.list)
            this.listIsLoading = false
            this.hasInitialLoaded = true
          }
        })
    },
    UpdateShoppingListItemTag (tagName, tagNameNew) {
      shoppinglist
        .UpdateShoppingListItemTag(this.id, tagName, tagNameNew)
        .catch((err) => {
          common.DisplayFailureToast(
            'Failed to update the shopping list tag' +
              '<br/>' +
              err.response.data.metadata.response
          )
        })
    },
    TimestampToCalendar (timestamp) {
      return common.TimestampToCalendar(timestamp)
    },
    LoopStart () {
      if (shoppinglistCommon.GetShoppingListAutoRefresh() === 'false') {
        return
      }
      this.intervalLoop = window.setInterval(() => {
        if (this.editing === true) {
          return
        }
        this.GetShoppingList()
        this.GetShoppingListItems()

        var now = new Date()
        var timePassed =
          now.getTime() / 1000 - this.loopCreated.getTime() / 1000
        if (timePassed >= 3600 / 4) {
          window.clearInterval(this.intervalLoop)
        }
      }, 3 * 1000)
    },
    LoopStop () {
      window.clearInterval(this.intervalLoop)
    },
    CheckDeviceIsMobile () {
      this.deviceIsMobile = common.DeviceIsMobile()
    },
    ManageStickyHeader () {
      this.HeaderIsSticky =
        window.pageYOffset > document.getElementById('ListName').offsetTop + 30
    },
    ResetLoopTime () {
      this.loopCreated = new Date()
    },
    FocusEl (name) {
      this.$nextTick(() => {
        this.$refs[name].focus()
      })
    },
    FocusName () {
      this.FocusEl('name')
    },
    FocusNotes () {
      this.FocusEl('notes')
    },
    FocusSearch () {
      this.FocusEl('search')
    },
    TagIsExcluded (tag) {
      return this.totalTagExcludeList.includes(tag)
    },
    GetFlatmates () {
      flatmates.GetAllFlatmates().then((resp) => {
        if (resp.data.list === null) {
          this.flatmates = []
          return
        }
        this.flatmates = resp.data.list
      })
    }
  },
  watch: {
    id () {
      console.log('ID changed')
      window.location.reload(false)
    },
    sortBy () {
      shoppinglistCommon.WriteShoppingListSortBy(this.sortBy)
      this.listIsLoading = true
      this.ResetLoopTime()
      this.LoopStop()
      this.LoopStart()
    },
    itemDisplayState () {
      this.listIsLoading = true
      this.GetShoppingListItems()
      shoppinglistCommon.WriteShoppingListObtainedFilter(
        this.id,
        this.itemDisplayState
      )
    },
    itemSearch () {
      shoppinglistCommon.WriteShoppingListSearch(this.id, this.itemSearch)
    },
    hasInitialLoaded () {
      this.canAnimate = true
    },
    completed () {
      var enableAnimations = common.GetEnableAnimations()
      if (
        this.completed === true &&
        enableAnimations !== 'false' &&
        this.canAnimate === true
      ) {
        common.Hooray()
      }
    },
    author () {
      if (this.authorNames !== '') {
        return
      }
      flatmates.GetFlatmate(this.author).then((resp) => {
        this.authorNames = resp.data.spec.names
      })
    },
    authorLast () {
      if (this.authorLastNames !== '') {
        return
      }
      flatmates.GetFlatmate(this.author).then((resp) => {
        this.authorLastNames = resp.data.spec.names
      })
    },
    templateId () {
      if (typeof this.templateId === 'undefined' || this.templateId === '') {
        return
      }
      shoppinglist.GetShoppingList(this.templateId).then((resp) => {
        this.templateListName = resp.data.spec.name
      })
    },
    isNewItemModalActive () {
      if (this.isNewItemModalActive !== false) {
        return
      }
      this.GetShoppingListItems()
    },
    isEditItemModalActive () {
      if (this.isEditItemModalActive !== false) {
        return
      }
      this.GetShoppingListItems()
    }
  },
  async beforeMount () {
    this.GetShoppingList()
    this.GetShoppingListItems()
    this.GetFlatmates()
    if (window.innerWidth <= 330) {
      this.ratherSmallerScreen = true
    }
  },
  async created () {
    this.CheckDeviceIsMobile()
    window.addEventListener('resize', this.CheckDeviceIsMobile, true)
    window.addEventListener('scroll', this.ManageStickyHeader, true)
    this.LoopStart()
    window.addEventListener('focus', this.ResetLoopTime, true)
    // TODO better way to do this? why does this not pull in through the data state function?
    this.itemDisplayState = shoppinglistCommon.GetShoppingListObtainedFilter(
      this.id
    )
  },
  mounted () {
    if (typeof this.ItemId !== 'undefined') {
      var el = this.$refs[this.ItemId][0].$el
      window.scrollTo(0, el.offsetTop)
    }
  },
  beforeDestroy () {
    this.LoopStop()
    window.removeEventListener('resize', this.CheckDeviceIsMobile, true)
    window.removeEventListener('scroll', this.ManageStickyHeader, true)
    window.removeEventListener('focus', this.ResetLoopTime, true)
  }
}
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
</style>
