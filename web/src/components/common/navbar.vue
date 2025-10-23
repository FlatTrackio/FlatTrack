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
  <div class="navbar-z">
    <b-loading
      v-model:active="pageLoading"
      :is-full-page="false"
      :can-cancel="false"
    />
    <section>
      <div class="block">
        <b-sidebar :fullheight="true" :open="true">
          <div class="hero is-bold navbar-text navbar-shadow ft">
            <!-- <div class="block"> -->
            <!--   <img -->
            <!--     src="" -->
            <!--     alt="FlatTrack logo" -->
            <!--     /> -->
            <!-- </div> -->
            <h1 class="title is-2">FlatTrack</h1>
          </div>
          <div class="p-1">
            <b-menu class="m-3">
              <b-menu-list label="General">
                <b-menu-item
                  icon="home"
                  label="Home"
                  tag="router-link"
                  :to="{ name: 'Home' }"
                />
                <b-menu-item
                  icon="information-outline"
                  label="My flat"
                  tag="router-link"
                  :to="{ name: 'My Flat' }"
                />
                <b-menu-item
                  icon="information-outline"
                  label="About FlatTrack"
                  tag="router-link"
                  :to="{ name: 'About FlatTrack' }"
                />
              </b-menu-list>
              <b-menu-list label="Apps">
                <b-menu-item
                  icon="format-list-checks"
                  label="Shopping list"
                  tag="router-link"
                  :to="{ name: 'Shopping list' }"
                />
                <b-menu-item
                  icon="account-group"
                  label="Flatmates"
                  tag="router-link"
                  :to="{ name: 'My Flatmates' }"
                />
                <b-menu-item
                  icon="apps"
                  label="Apps"
                  tag="router-link"
                  :to="{ name: 'Apps' }"
                />
              </b-menu-list>
              <b-menu-list v-if="canUserAccountAdmin" label="Admin">
                <b-menu-item
                  icon="account-group"
                  label="Flatmates"
                  tag="router-link"
                  :to="{ name: 'Admin accounts' }"
                />
                <b-menu-item
                  icon="cogs"
                  label="Admin apps"
                  tag="router-link"
                  :to="{ name: 'Admin home' }"
                />
              </b-menu-list>
              <b-menu-list label="Help">
                <b-menu-item
                  icon="open-in-new"
                  label="FlatTrack site"
                  tag="a"
                  target="_blank"
                  href="https://flattrack.io/"
                />
                <b-menu-item
                  icon="phone"
                  label="Contact admin"
                  tag="router-link"
                  :to="{ name: 'My Flatmates', query: { group: 'admin' } }"
                />
              </b-menu-list>
              <b-menu-list label="Account">
                <b-menu-item
                  icon="account-circle"
                  label="Profile"
                  tag="router-link"
                  :to="{ name: 'Account Profile' }"
                />
                <b-menu-item
                  icon="account"
                  label="My Account"
                  tag="router-link"
                  :to="{ name: 'Account' }"
                />
                <b-menu-item
                  icon="exit-to-app"
                  label="Sign out"
                  @click="signOut"
                />
              </b-menu-list>
            </b-menu>
          </div>
        </b-sidebar>
      </div>
    </section>
  </div>
</template>

<script>
  import common from "@/common/common";
  import cani from "@/requests/authenticated/can-i";

  export default {
    name: "NavBar",
    data() {
      return {
        isActive: true,
        overlay: false,
        fullheight: true,
        fullwidth: false,
        pageLoading: true,
        flatName: "My flat",
        canUserAccountAdmin: false,
      };
    },
    async beforeMount() {
      cani.GetCanIgroup("admin").then((resp) => {
        this.canUserAccountAdmin = resp.data.data;
        this.pageLoading = false;
      });
    },
    methods: {
      signOut() {
        common.SignoutDialog(this.$buefy);
      },
    },
  };
</script>

<style>
  .sidebar-content {
    /* NOTE: this is a hack */
    display: unset !important;
  }
  .navbar-z {
    z-index: 100;
  }

  .p-1 {
    padding-left: 1em !important;
    padding-bottom: 1em !important;
    padding-right: 1em !important;
  }

  .p-1 .menu {
    margin-top: 23px;
  }

  .navbar-text {
    text-align: center;
    padding: 15px;
  }

  .navbar-text .title {
    color: white;
  }

  .align-bottom {
    bottom: 0;
    left: 0;
    width: 100%;
    height: 40px;
    background-color: #496c8a40;
    position: absolute;
  }
  ul.menu-list > li > a > span:nth-child(2) {
    padding: 0 0.5rem;
  }

  div.hero.ft,
  section.hero.ft {
    background-image: linear-gradient(
      141deg,
      #208fbc 0%,
      hsl(207, 61%, 53%) 71%,
      #4d83db 100%
    ) !important;
    color: white;
  }
  section.hero.ft div p {
    color: white;
  }
</style>
