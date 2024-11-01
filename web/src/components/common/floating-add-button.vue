<!--
     This program is free software: you can redistribute it and/or modify
     it under the terms of the Affero GNU General Public License as published by
     the Free Software Foundation, either version 3 of the License, or
     (at your option) any later version.

     This program is distributed in the hope that it will be useful,
     but WITHOUT ANY WARRANTY; without even the implied warranty of
     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
     GNU General Public License for more details.

     You should have received a copy of the Affero GNU General Public License
     along with this program.  If not, see <https://www.gnu.org/licenses/>.
-->

<template>
  <div
    :class="deviceIsMobile ? 'floating-add-button-mobile' : ''"
    class="floating-add-button"
  >
    <md-speed-dial md-direction="bottom md-bottom-right">
      <md-speed-dial-target
        class="floating-add-button-colour"
        @click="goToRouterLinkOrRunFunc"
      >
        <md-icon>add</md-icon>
      </md-speed-dial-target>
    </md-speed-dial>
  </div>
</template>

<script>
import common from '@/common/common'

export default {
  name: 'FloatingAddButton',
  props: {
    routerLink: {
      default: "/",
      type: String,
    },
    func: {
      type: Function,
      default: () => {},
    }
  },
  data () {
    return {
      deviceIsMobile: false
    }
  },
  async beforeMount () {
    this.deviceIsMobile = common.DeviceIsMobile()
  },
  methods: {
    goToRouterLinkOrRunFunc () {
      var routerLink = this.routerLink
      var func = this.func
      if (routerLink !== '' && typeof routerLink !== 'undefined') {
        this.$router.push(routerLink)
        return
      }
      if (typeof func === 'function') {
        func()
      }
    }
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
