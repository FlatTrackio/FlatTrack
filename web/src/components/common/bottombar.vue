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
  <div :class="ratherSmallerScreen ? 'bottombar bottombar-fixed' : 'bottombar'">
    <b-loading
      v-model:active="pageLoading"
      :is-full-page="false"
      :can-cancel="false"
    />
    <div class="bottombar bbitems">
      <div v-for="item in routes" :key="item">
        <div class="bbitem" @click="$router.push({ name: item.routeName })">
          <div :class="$route.name === item.routeName ? 'active' : ''">
            <b-icon :icon="item.icon" size="is-small" />
            <span> {{ item.name }} </span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
  import cani from "@/requests/authenticated/can-i";

  export default {
    name: "BottomBar",
    data() {
      return {
        pageLoading: true,
        canUserAccountAdmin: false,
        ratherSmallScreen: false,
      };
    },
    computed: {
      routes() {
        return [
          {
            name: "Home",
            icon: "home",
            routeName: "Home",
            requiresAdmin: false,
          },
          {
            name: "Apps",
            icon: "apps",
            routeName: "Apps",
            requiresAdmin: false,
          },
          {
            name: "Account",
            icon: "account",
            routeName: "Account",
            requiresAdmin: false,
          },
          {
            name: "Admin",
            icon: "web",
            routeName: "Admin home",
            requiresAdmin: true,
          },
        ].filter(i => !i.requiresAdmin || i.requiresAdmin && this.canUserAccountAdmin === true);
      },
    },
    async beforeMount() {
      this.CanIadmin();
      if (window.innerWidth >= 330) {
        this.ratherSmallerScreen = true;
      }
    },
    methods: {
      CanIadmin() {
        cani.GetCanIgroup("admin").then((resp) => {
          this.canUserAccountAdmin = resp.data.data;
          this.pageLoading = false;
        });
      },
    },
  };
</script>

<style>
  .bottombar-fixed {
    position: fixed;
  }

  .bottombar {
    width: 100%;
    height: 3.5rem;
    bottom: 0;
    display: inline-flex;
    align-items: flex-end;
    z-index: 100;
    background-color: #ffffff52;
    box-shadow: 0 5px 5px -3px #0003, 0 8px 10px 1px #00000024,
      0 3px 14px 2px #0000001f;
    backdrop-filter: blur(10px);
  }

  .bottombar:before {
    content: "";
    box-shadow: inset 0 0 0 200px rgba(255, 255, 255, 0.3);
    filter: blur(10px);
  }

  .bbitems div {
    width: 100%;
    user-select: none;
  }
  .bbitems div .bbitem div {
    position: static;
    display: flex;
    flex-direction: column;
    align-items: center;
    width: 100%;
    padding: 0.5rem 0;
  }
  @keyframes bbactive {
    0%,
    40% {
      background-color: #81d1f02b;
      font-size: 1rem;
      height: 100%;
    }
    100% {
      background-color: #c1c1c13d;
      font-size: 1.09rem;
      height: calc(100% - 0.05px);
    }
  }
  .bbitems div .bbitem:hover:has(div.active) {
    background-color: #81d1f02b;
    transition: background-color 1s;
  }
  .bbitems div .bbitem:has(div.active) {
    background-color: #c1c1c13d;
    box-shadow: 0 0 6px -4px black inset;
    transition: background-color 1s;
    font-size: 1.09rem;
    animation-name: bbactive;
    animation-duration: 0.7s;
    height: calc(100% - 0.05px);
  }
  .bbitems div .bbitem:has(div.active) div {
    top: 0.09rem;
    position: relative;
  }

  @media (prefers-color-scheme: dark) {
    .bottombar {
      background-color: #2b2b2b4a;
    }
    .bbitems div .bbitem:has(div.active) {
      background-color: #d4d4d499;
      color: black;
    }
  }

  [data-theme="dark"] .bottombar {
    background-color: #2b2b2b4a;
  }
  [data-theme="dark"] .bbitems div .bbitem:has(div.active) {
    background-color: #d4d4d499;
    color: black;
  }
  @media (prefers-color-scheme: light) {
    .bottombar {
      background-color: #ffffff52;
    }
    .bbitems div .bbitem:has(div.active) {
      background-color: #c1c1c13d;
      color: unset;
    }
  }

  [data-theme="light"] .bottombar {
    background-color: #ffffff52;
  }
  [data-theme="light"] .bbitems div .bbitem:has(div.active) {
    background-color: #81d1f02b;
    color: black;
  }
</style>
