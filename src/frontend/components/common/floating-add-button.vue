<template>
  <div :class="deviceIsMobile ? 'floating-add-button-mobile' : ''" class="floating-add-button">
    <md-speed-dial md-direction="bottom md-bottom-right">
      <md-speed-dial-target @click="goToRouterLinkOrRunFunc" class="floating-add-button-colour">
        <md-icon>add</md-icon>
      </md-speed-dial-target>
    </md-speed-dial>
  </div>
</template>

<script>
import common from '@/frontend/common/common'

export default {
  name: 'floating add button',
  data () {
    return {
      deviceIsMobile: false
    }
  },
  props: {
    routerLink: Object,
    func: Function
  },
  methods: {
    goToRouterLinkOrRunFunc () {
      // TODO remove print
      var routerLink = this.routerLink
      var func = this.func
      console.log({ routerLink, func })
      if (routerLink !== '' && typeof routerLink !== 'undefined') {
        this.$router.push(routerLink)
      } else if (typeof func === 'function') {
        func()
      }
    }
  },
  async beforeMount () {
    this.deviceIsMobile = common.DeviceIsMobile()
  }
}
</script>

<style scope>
.floating-add-button {
    display: block;
    position: fixed;
    bottom: 0;
    right: 0;
    z-index: 100;
}

.floating-add-button-mobile {
    margin-bottom: 50px;
}

.floating-add-button-colour {
    background-color: #448aff;
}
</style>
