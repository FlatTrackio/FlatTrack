<template>
  <div>
    <section>
      <div class="block">
        <b-sidebar
          type="is-light"
          :fullheight="fullheight"
          :fullwidth="fullwidth"
          :overlay="overlay"
          :right="right"
          can-cancel="false"
          :open.sync="open">
          <div class="p-1">
            <div class="navbar-text">
              <h1 class="title is-2">FlatTrack</h1>
            </div>
            <b-menu>
              <b-menu-list label="General">
                <b-menu-item icon="home" label="Home"></b-menu-item>
                <b-menu-item icon="account_box" label="Profile"></b-menu-item>
              </b-menu-list>
              <b-menu-list label="Apps">
                <b-menu-item icon="information-outline" label="Shopping list"></b-menu-item>
                <b-menu-item icon="people" label="Flatmates"></b-menu-item>
              </b-menu-list>
              <b-menu-list label="Admin">
                <b-menu-item icon="people" label="Flatmates"></b-menu-item>
                <b-menu-item icon="settings" label="Settings"></b-menu-item>
              </b-menu-list>
              <b-menu-list label="Help">
                <b-menu-item icon="web" label="FlatTrack help"></b-menu-item>
                <b-menu-item icon="phone" label="Contact admin"></b-menu-item>
              </b-menu-list>
            </b-menu>
          </div>
        </b-sidebar>
      </div>
    </section>
  </div>
</template>

<script>
export default {
  data () {
    return {
      isActive: true,
      open: true,
      overlay: false,
      fullheight: true,
      fullwidth: false,
      right: false
    }
  }
}
</script>

<style>
.p-1 {
  padding: 1em;
}

.p-1 .menu {
  margin-top: 23px;
}

.navbar-text {
  text-align: center;
}
</style>
