<template>
  <div class="bottombar">
    <md-bottom-bar class="md-accent bottombar-background" md-sync-route>
      <md-bottom-bar-item to="/" exact md-label="Home" md-icon="home"></md-bottom-bar-item>
      <md-bottom-bar-item to="/apps" md-label="Apps" md-icon="apps"></md-bottom-bar-item>
      <md-bottom-bar-item to="/profile" md-label="Profile" md-icon="account_box"></md-bottom-bar-item>
      <md-bottom-bar-item to="/admin" md-label="Admin" md-icon="web"></md-bottom-bar-item>
    </md-bottom-bar>
  </div>
</template>

<script>
export default {
  name: 'bottombar'
}
</script>

<style>
.bottombar {
    position: fixed;
    width: 100%;
    bottom: 0;
    display: inline-flex;
    align-items: flex-end;
    background: rbga(#209cee, 0.8);
}

.bottombar-background {
    background-color: white;
}
</style>

<style lang="scss" scoped>
  .phone-viewport {
    width: 322px;
    height: 200px;
    display: inline-flex;
    align-items: flex-end;
    overflow: hidden;
    border: 1px solid rgba(#000, .26);
    background: rgba(#000, .06);
  }
</style>
