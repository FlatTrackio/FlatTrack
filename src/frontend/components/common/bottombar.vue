<template>
  <div class="bottombar">
    <md-bottom-bar class="md-accent bottombar-background" md-sync-route>
      <md-bottom-bar-item to="/" exact md-label="Home" md-icon="home"></md-bottom-bar-item>
      <md-bottom-bar-item to="/apps" md-label="Apps" md-icon="apps"></md-bottom-bar-item>
      <md-bottom-bar-item to="/profile" md-label="Profile" md-icon="account_box"></md-bottom-bar-item>
      <md-bottom-bar-item to="/admin" md-label="Admin" md-icon="web" v-if="canUserAccountAdmin"></md-bottom-bar-item>
    </md-bottom-bar>
  </div>
</template>

<script>
import cani from '@/frontend/requests/authenticated/can-i'

export default {
  name: 'bottombar',
  data () {
    return {
      canUserAccountAdmin: false
    }
  },
  methods: {
    CanIadmin () {
      cani.GetCanIgroup('admin').then(resp => {
        this.canUserAccountAdmin = resp.data.data
      })
    }
  },
  async beforeMount () {
    this.CanIadmin()
  }
}
</script>

<style>
.bottombar {
    position: fixed;
    width: 100%;
    bottom: 0;
    display: inline-flex;
    align-items: flex-end;
    #background: rbga(#209cee, 0.8);
    z-index: 100;
    background: inherit;
}

.bottombar:before {
    content: '';
    box-shadow: inset 0 0 0 200px rgba(255,255,255,0.3);
    filter: blur(10px);
    background: inherit;
}

.bottombar-background {
    background-color: hsla(0, 0%, 100%, 0.73);
    backdrop-filter: blur(5px);
}
</style>
