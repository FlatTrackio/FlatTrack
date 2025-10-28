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
  <div
    :class="deviceIsMobile ? 'floating-add-button-mobile' : ''"
    class="floating-add-button"
  >
    <div
      class="floating-add-button-colour"
      @click="goToRouterLinkOrRunFunc"
    >
      <b-icon
        icon="plus"
        size="is-medium"
        type="is-white"
      />
    </div>
  </div>
</template>

<script>
  import common from "@/common/common";

  export default {
    name: "FloatingAddButton",
    props: {
      routerLink: String,
      func: Function,
    },
    data() {
      return {
        deviceIsMobile: false,
      };
    },
    async beforeMount() {
      this.deviceIsMobile = common.DeviceIsMobile();
    },
    methods: {
      goToRouterLinkOrRunFunc() {
        var routerLink = this.routerLink;
        var func = this.func;
        if (routerLink !== "" && typeof routerLink !== "undefined") {
          this.$router.push(routerLink);
          return;
        }
        if (typeof func === "function") {
          func();
        }
      },
    },
  };
</script>

<style scope>
  .floating-add-button {
    display: block;
    position: fixed;
    bottom: 0;
    right: 0;
    z-index: 100;
    margin-bottom: 1rem;
    margin-right: 1rem;
  }

  .floating-add-button div {
    display: flex;
    border-radius: 50%;
    padding: 0.75rem;
    box-shadow: 0 0 21px -6px #000000b0;
  }

  .floating-add-button-mobile {
    margin-bottom: 5rem;
    margin-right: 1.5rem;
  }

  .floating-add-button-colour {
    background-color: #448aff;
  }
</style>
