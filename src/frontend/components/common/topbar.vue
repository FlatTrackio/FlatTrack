<template>
  <div id="topbar">
    <b-navbar type="is-info" shadow="true" transparent="true">
      <template slot="brand">
        <b-navbar-item tag="router-link" :to="{ path: '/' }">
          <h1 class="title is-5" style="color: #fff;">FlatTrack (My Flat)</h1>
        </b-navbar-item>
      </template>
      <template slot="start">
        <b-navbar-item tag="router-link" :to="{ path: '/' }">
          Home
        </b-navbar-item>
        <b-navbar-item href="https://flattrack.io/help">
          Help
        </b-navbar-item>
      </template>
    </b-navbar>
  </div>
</template>

<script>
export default {
  name: 'topbar'
}
</script>

<style>
</style>
