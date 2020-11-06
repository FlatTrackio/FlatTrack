<template>
  <div :class="deviceIsMobile ? 'floating-add-button-mobile' : ''" class="floating-add-button">
    <md-speed-dial md-direction="bottom md-bottom-right">
      <md-speed-dial-target @click="goToRefOrRunFunc(path, func)" class="floating-add-button-colour">
        <md-icon>add</md-icon>
      </md-speed-dial-target>
    </md-speed-dial>
  </div>
</template>

<script>
import common from '@/common/common'

export default {
  name: 'floating add button',
  data () {
    return {
      deviceIsMobile: false
    }
  },
  props: {
    path: String,
    func: Function
  },
  methods: {
    goToRef (ref) {
      this.$router.push({ path: ref })
    },
    goToRefOrRunFunc (ref, func) {
      if (ref !== '' && typeof ref !== 'undefined') {
        this.goToRef(ref)
        return
      }
      if (typeof func === 'function') {
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
