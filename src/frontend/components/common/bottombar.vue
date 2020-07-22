<template>
  <div :class="ratherSmallerScreen ? 'bottombar bottombar-fixed' : 'bottombar'">
    <b-loading :is-full-page="false" :active.sync="pageLoading" :can-cancel="false"></b-loading>
    <md-bottom-bar class="md-accent bottombar-background" md-sync-route>
      <md-bottom-bar-item :to="{ name: 'Home' }" exact md-label="Home" md-icon="home"></md-bottom-bar-item>
      <md-bottom-bar-item :to="{ name: 'Apps' }" md-label="Apps" md-icon="apps"></md-bottom-bar-item>
      <md-bottom-bar-item :to="{ name: 'Account' }" md-label="My Account" md-icon="account_box"></md-bottom-bar-item>
      <md-bottom-bar-item :to="{ name: 'Admin home' }" md-label="Admin" md-icon="web" v-if="canUserAccountAdmin"></md-bottom-bar-item>
    </md-bottom-bar>
  </div>
</template>

<script>
import cani from '@/frontend/requests/authenticated/can-i'
import common from '@/frontend/common/common'

export default {
  name: 'bottombar',
  data () {
    return {
      pageLoading: true,
      canUserAccountAdmin: false,
      ratherSmallScreen: false
    }
  },
  methods: {
    CanIadmin () {
      cani.GetCanIgroup('admin').then(resp => {
        this.canUserAccountAdmin = resp.data.data
        this.pageLoading = false
      })
    }
  },
  async beforeMount () {
    this.CanIadmin()
    if (window.innerWidth >= 330) {
      this.ratherSmallerScreen = true
    }
  }
}
</script>

<style>
.bottombar-fixed {
    position: fixed;
}

.bottombar {
    #position: fixed;
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
