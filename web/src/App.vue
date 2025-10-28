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
  <div id="app">
    <div>
      <topbar v-if="onMobile && !publicPages" />
      <navbar v-if="!onMobile && !publicPages && displayNavigationBar" />
    </div>
    <div
      class="pad-bottom full-height main-view-container"
      :class="{ 'pad-left': !publicPages && !onMobile && displayNavigationBar }"
    >
      <router-view
        :key="$route.fullPath"
        class="main-view"
      />
    </div>
    <div>
      <bottombar v-if="onMobile && !publicPages" />
    </div>
  </div>
</template>

<script>
  import theme from "@/common/theme";
  import routerCommon from "@/router/common";
  import navbar from "@/components/common/navbar.vue";
  import topbar from "@/components/common/topbar.vue";
  import bottombar from "@/components/common/bottombar.vue";

  export default {
    name: "App",
    components: {
      navbar: navbar,
      topbar: topbar,
      bottombar: bottombar,
    },
    data() {
      return {
        onMobile: false,
        displayNavigationBar: true,
      };
    },
    computed: {
      publicPages() {
        return routerCommon.isPublicRoute(this.$route);
      },
    },
    created() {
      this.adjustForMobile();
      window.addEventListener("resize", this.adjustForMobile.bind(this));
      document.documentElement.setAttribute('data-theme', theme.GetTheme().name)
    },
    methods: {
      adjustForMobile() {
        var vm = this;
        vm.onMobile = window.innerWidth <= 870;
        if (vm.onMobile) {
          vm.displayNavigationBar = false;
        } else {
          vm.displayNavigationBar = true;
        }
      },
    },
  };
</script>

<style lang="scss">
  $material-design-icons-font-path: "~@mdi/font/";
  @use "@mdi/font";
  @use "bulma/sass" with (
    $primary: #00a7d6,
    $link: #5487db,
    $info: #3e8ed0,
    $success: #48c78e,
    $warning: #ffe08a,
    $danger: #f14668,

    $family-sans-serif: (
      "Nunito",
      BlinkMacSystemFont,
      -apple-system,
      "Segoe UI",
      "Roboto",
      "Oxygen",
      "Ubuntu",
      "Cantarell",
      "Fira Sans",
      "Droid Sans",
      "Helvetica Neue",
      "Helvetica",
      "Arial",
      sans-serif,
    ),
    $size-1: 3rem,
    $size-2: 2.5rem,
    $size-3: 2rem,
    $size-4: 1.5rem,
    $size-5: 1.25rem,
    $size-6: 1rem,
    $size-7: 0.75rem,

    $radius-small: 2px,
    $radius: 4px,
    $radius-large: 6px,
    $radius-rounded: 9999px
  );

  @use "buefy/src/scss/buefy";
</style>

<style>
  body {
    font-family: BlinkMacSystemFont, -apple-system, Segoe UI, Roboto, Oxygen,
      Ubuntu, Cantarell, Fira Sans, Droid Sans, Helvetica Neue, Helvetica, Arial,
      sans-serif !important;
  }

  .darken a.is-disabled {
    color: black;
  }

  .card-margin {
    margin-top: 10px;
    margin-bottom: 10px;
  }

  .pad-left {
    margin-left: 260px;
  }

  .pad-top {
    margin-top: 20px;
  }

  .pad-bottom {
    margin-bottom: 100px;
  }

  .full-height {
    height: 100%;
  }

  html {
    background-color: #fbfbfb;
  }

  .pointer-cursor-on-hover {
    cursor: pointer;
  }

  .card {
    user-select: none;
    border: 1px dashed rgba(0, 0, 0, 0);
    transition: box-shadow 0.4s, border 0.2s;
  }

  .card:hover {
    transition: box-shadow 0.4s, background-color 0.3s, border 0.2s, filter 0.3s;
    box-shadow: black 0 0px 45px -30px;
    border: 1px solid darkgray;
  }

  .card:active {
    transition: height 0.1s, width 0.1s, margin 0.1s, box-shadow 0.4s, filter 0.3s;
    height: calc(100% - 8px);
    border: 1px dashed #000000;
    box-shadow: black 0 0px 66px -27px;
    top: 1px;
  }

  .hero-body,
  .section {
    padding-top: 2rem;
    padding-right: 1rem;
    padding-bottom: 4rem;
    padding-left: 1rem;
  }

  .form-width {
    width: 30rem;
    margin: auto;
  }

  .button {
    box-shadow: #ffffff -1px -1px 4px 0px, #bababa 1px 1px 4px 0px;
  }

  .navbar-shadow {
    box-shadow: 0px -9px 23px 4px #292929;
  }

  .remove-shadow {
    box-shadow: none;
  }

  .notes-highlight {
    background-color: #fff;
    padding: 10px;
    box-shadow: 0 0 30px -35px #000;
    border-radius: 5px;
  }

  *::-webkit-scrollbar {
    width: 12px;
  }
  *::-webkit-scrollbar-track {
    background: white;
  }
  *::-webkit-scrollbar-thumb {
    background-color: #c3c3c3;
    border-radius: 20px;
    border: 3px solid white;
    transition: background-color 0.3s, border 0.3s;
  }
  *::-webkit-scrollbar-thumb:hover {
    background-color: #b1b1b1;
    border-radius: 20px;
    border: 2px solid white;
    transition: background-color 0.3s, border 0.3s;
  }

  nav.breadcrumb {
    display: flex;
  }

  nav.breadcrumb ul {
    align-items: center;
  }

  div.item-page section.modal-card-body {
    padding-bottom: 5rem;
  }

  @media (max-width: 870px) {
    .form-width {
      width: calc(100% - 25px);
    }
  }

  @media (prefers-color-scheme: dark) {
    .card:active,
    .card:hover {
      filter: brightness(1.5);
      transition: filter 0.3s;
    }
    .card:active {
      border: 1px dashed #ffffff;
    }
    .notes-highlight {
      background-color: #000;
    }
  }
  [data-theme="dark"] {
    .card:active,
    .card:hover {
      filter: brightness(1.5);
      transition: filter 0.3s;
    }
    .card:active {
      border: 1px dashed #ffffff;
    }
    .notes-highlight {
      background-color: #000;
    }
  }
  @media (prefers-color-scheme: light) {
    .card:active,
    .card:hover {
      filter: unset;
      transition: unset;
    }
    .card:active {
      border: 1px dashed darkgray;
    }
    .notes-highlight {
      background-color: #fff;
    }
  }
  [data-theme="light"] {
    .card:active,
    .card:hover {
      filter: unset;
      transition: unset;
    }
    .card:active {
      border: 1px dashed darkgray;
    }
    .notes-highlight {
      background-color: #fff;
    }
  }
</style>
