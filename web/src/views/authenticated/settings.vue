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
  <div>
    <div class="container">
      <section class="section">
        <breadcrumb
          back-link-name="Account"
          :current-page-name="$route.name"
        />
        <h1 class="title is-1">
          Settings
        </h1>
        <p class="subtitle is-3">
          Manage settings for this device
        </p>
        <b-loading
          v-model:active="pageLoading"
          :is-full-page="false"
          :can-cancel="false"
        />

        <b-field label="Miscellaneous">
          <b-checkbox
            v-model="enableAnimations"
            size="is-medium"
          >
            Enable Animations
          </b-checkbox>
        </b-field>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/common/common'
  import breadcrumb from "@/components/common/breadcrumb.vue";

export default {
  name: 'SettingsHome',
    components: {
      breadcrumb,
    },
  data () {
    return {
      enableAnimations: common.GetEnableAnimations() !== 'false'
    }
  },
  watch: {
    enableAnimations () {
      common.WriteEnableAnimations(this.enableAnimations)
      if (this.enableAnimations !== 'false') {
        common.Hooray()
      }
    }
  },
  methods: {
    CopyHrefToClipboard () {
      common.CopyHrefToClipboard()
    }
  },
}
</script>
