<template>
  <div>
    <div id="bottombar">
      <b-tabs type="is-toggle" size="is-medium" class="block" position="is-centered" expanded>
        <b-tab-item label="Home" icon="google-photos"></b-tab-item>
        <b-tab-item label="Menu" icon="library-music"></b-tab-item>
        <b-tab-item label="Profile" icon="video"></b-tab-item>
      </b-tabs>
    </div>
  </div>
</template>

<script>
export default {
  name: 'bottombar'
}
</script>

<style>
#bottombar .b-tabs .tab-content {
  display: none;
}
</style>
