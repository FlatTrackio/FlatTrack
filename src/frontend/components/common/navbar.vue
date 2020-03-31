<template>
  <div id="navbar">
    <div id="navbar-inner">
      <div class="container is-fluid">
        <h1 class="title is-2"> My Flat </h1>
      </div>
      <br/>
      <br/>
      <br/>
      <b-menu>
        <b-menu-list label="General">
          <b-menu-item icon="home" label="Home" tag="router-link" to="/"></b-menu-item>
          <b-menu-item icon="account_box" label="Profile" tag="router-link" to="/profile"></b-menu-item>
        </b-menu-list>
        <b-menu-list label="Apps">
          <b-menu-item icon="information-outline" label="Shopping list" tag="router-link" to="/apps/shopping-list"></b-menu-item>
          <b-menu-item icon="information-outline" label="Flatmates" tag="router-link" to="/apps/flatmates"></b-menu-item>
        </b-menu-list>
        <b-menu-list label="Admin">
          <b-menu-item icon="information-outline" label="Flatmates"></b-menu-item>
          <b-menu-item icon="information-outline" label="Settings"></b-menu-item>
        </b-menu-list>
      </b-menu>
    </div>
  </div>
</template>

<script type="javascript">
export default {
  data () {
    return {
      isActive: true
    }
  }
}
</script>

<style>
#navbar {
    position: fixed;
    top: 0;
    bottom: 0;
    left: 0;
    width: 265px;
    z-index: 8;
    font-family: 'Raleway', sans-serif;
    color: rgba(0, 0, 0, 0.75);
    font-weight: 700;
    background-color: #ffffff;
    box-shadow: black -30px 0px 57px;
}

#navbar-inner {
  position: fixed;
  top: 15px;
  left: 15px;
}
</style>
