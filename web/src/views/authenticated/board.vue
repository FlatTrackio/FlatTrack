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
        <nav
          class="breadcrumb is-medium has-arrow-separator"
          aria-label="breadcrumbs"
        >
          <ul>
            <li><router-link :to="{ name: 'Apps' }">Apps</router-link></li>
            <li class="is-active">
              <router-link :to="{ name: 'Board' }">Board</router-link>
            </li>
          </ul>
        </nav>
        <h1 class="title is-1">Board</h1>
        <p class="subtitle is-3">Share ideas and updates</p>
        <br />
        <div>
          <label class="label">Search for posts</label>
          <b-field>
            <b-input
              icon="magnify"
              size="is-medium"
              placeholder="Enter text here"
              type="search"
              expanded
              v-model="listSearch"
              ref="search"
            >
            </b-input>
          </b-field>
          <b-loading
            :is-full-page="false"
            :active.sync="pageLoading"
            :can-cancel="false"
          ></b-loading>
          <section>
            <div
              class="card pointer-cursor-on-hover"
              @click="
                $router.push({
                  name: 'New board post',
                  query: { name: listSearch || undefined },
                })
              "
            >
              <div class="card-content">
                <div class="media">
                  <div class="media-left">
                    <b-icon icon="cart-plus" size="is-medium"> </b-icon>
                  </div>
                  <div class="media-content">
                    <p class="title is-4">Add a post</p>
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
        <floatingAddButton
          :routerLink="{
            name: 'New board post',
            query: { name: listSearch || undefined },
          }"
          v-if="displayFloatingAddButton"
        />
        <br />
        <div v-if="listsFiltered.length > 0">
          <boardListCardView
            :list="list"
            :authors="authors"
            :lists="lists"
            :index="index"
            v-for="(list, index) in listsFiltered"
            v-bind:key="list"
            :deviceIsMobile="deviceIsMobile"
            @lists="
              (l) => {
                lists = l;
              }
            "
          />
          <br />
          <p>{{ listsFiltered.length }} post(s)</p>
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
                      listSearch === '' && lists.length === 0 && !pageLoading
                    "
                  >
                    No posts yet. Try making one!
                  </p>
                  <p
                    class="subtitle is-4"
                    v-else-if="listSearch !== '' && !pageLoading"
                  >
                    No posts found.
                  </p>
                  <p class="subtitle is-4" v-else-if="pageLoading">
                    Loading lists...
                  </p>
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
import board from "@/requests/authenticated/board";
import flatmates from "@/requests/authenticated/flatmates";
import cani from "@/requests/authenticated/can-i";
import { DialogProgrammatic as Dialog } from "buefy";

export default {
  name: "board-list",
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
    };
  },
  components: {
    boardListCardView: () =>
      import("@/components/authenticated/board-list-card-view.vue"),
    floatingAddButton: () =>
      import("@/components/common/floating-add-button.vue"),
  },
  computed: {
    listsFiltered() {
      return this.lists.filter((item) => {
        return this.ListDisplayState(item);
      });
    },
  },
  methods: {
    GetBoardItems() {
      board
        .GetBoardItems(undefined, undefined, undefined)
        .then((resp) => {
          this.pageLoading = false;
          this.lists = resp.data.list || [];
        })
        .catch(() => {
          common.DisplayFailureToast(
            "Hmmm seems somethings gone wrong loading the posts"
          );
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
  watch: {
    sortBy() {
      this.listIsLoading = true;
      this.GetBoardItems();
    },
  },
  async beforeMount() {
    cani.GetCanIgroup("admin").then((resp) => {
      this.canUserAccountAdmin = resp.data.data;
    });
    this.GetBoardItems();
  },
  beforeDestroy() {
    window.removeEventListener("resize", this.CheckDeviceIsMobile, true);
  },
  async created() {
    this.CheckDeviceIsMobile();
    window.addEventListener("resize", this.CheckDeviceIsMobile, true);
  },
};
</script>

<style scoped></style>
