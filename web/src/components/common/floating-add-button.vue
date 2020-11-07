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
