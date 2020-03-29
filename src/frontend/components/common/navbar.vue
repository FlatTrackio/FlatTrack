<template>
  <div id="navbar">
    <div id="navbar-inner">
      <div class="container is-fluid">
        <h1 class="title is-2"> My Flat </h1>
      </div>
      <br/>
      <br/>
      <b-menu-list>
        <b-menu-item icon="information-outline" label="Home"></b-menu-item>
      </b-menu-list>
      <br/>
      <b-menu>
        <b-menu-list label="Apps">
          <b-menu-item icon="information-outline" label="Shopping list"></b-menu-item>
          <b-menu-item icon="information-outline" label="Flatmates"></b-menu-item>
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
}

#navbar-inner {
  position: fixed;
  top: 15px;
  left: 15px;
}
</style>
