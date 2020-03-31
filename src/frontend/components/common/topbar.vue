<template>
  <div class="topbar">
    <b-navbar type="is-info" shadow="true" fixed-top="true" transparent="true">
      <template slot="brand">
        <b-navbar-item tag="router-link" :to="{ path: '/' }">
          <h1 class="title is-5" style="color: #fff;">FlatTrack (My Flat)</h1>
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
