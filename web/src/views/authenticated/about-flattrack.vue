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
        <nav
          class="breadcrumb is-medium has-arrow-separator"
          aria-label="breadcrumbs"
        >
          <ul>
            <li>
              <router-link :to="{ name: 'Home' }">
                Home
              </router-link>
            </li>
            <li class="is-active">
              <router-link :to="{ name: 'About FlatTrack' }">
                About FlatTrack
              </router-link>
            </li>
          </ul>
          <b-button
            icon-left="content-copy"
            size="is-small"
            @click="CopyHrefToClipboard()"
          />
        </nav>
        <h1 class="title is-1">
          About FlatTrack
        </h1>
        <p class="subtitle is-3">
          What is FlatTrack?
        </p>
        <b-message
          type="is-primary"
          has-icon
          icon="help"
        >
          <p class="is-size-5">
            <a
              href="https://flattrack.io"
              target="_blank"
              rel="noreferrer"
            >FlatTrack</a>
            is a
            <a
              href="https://simple.wikipedia.org/wiki/Free_and_open-source_software"
              target="_blank"
              rel="noreferrer"
            >Free and Open Source</a>
            collaboration software for flats / community houses and homes with
            the goals of <b>easing common tasks</b> in living environments,
            <b>enabling closer collaboration</b>, and
            <b>empowering humans who live together</b>.
          </p>
        </b-message>
        <p class="subtitle is-3">
          Collaborate
        </p>
        <b-message
          type="is-primary"
          has-icon
          icon="help"
        >
          <p class="is-size-5">
            FlatTrack is all about community.<br>
            Together, we can define the standards on technology's assistance in
            flat life.<br>
            Contribution is welcome and you are too! Please join the community
            over on
            <a
              href="https://gitlab.com/flattrack/flattrack"
              target="_blank"
              rel="noreferrer"
            >GitLab</a>.
          </p>
        </b-message>
        <p class="subtitle is-3">
          This FlatTrack instance
        </p>
        <b-message
          type="is-warning"
          has-icon
          icon="information-outline"
        >
          <p class="is-size-5">
            <b>Version</b>: {{ version || "Unknown" }}
            <span v-if="version !== versionFrontend">(frontend {{ versionFrontend }})</span>
            <br>
            <b>Commit hash</b>:
            <a
              v-if="commitHash !== '???' && typeof commitHash !== 'undefined'"
              :href="
                'https://gitlab.com/flattrack/flattrack/-/commit/' + commitHash
              "
              target="_blank"
              rel="noreferrer"
            >{{ commitHash }}</a>
            <span v-else> Unknown </span>
            <a
              v-if="commitHash !== commitHashFrontend"
              :href="
                'https://gitlab.com/flattrack/flattrack/-/commit/' +
                  commitHashFrontend
              "
              target="_blank"
              rel="noreferrer"
            >{{ commitHashFrontend }}</a>
            <br>
            <b>Mode</b>: {{ mode || "Unknown" }}
            <span v-if="mode !== modeFrontend">(frontend {{ modeFrontend }})</span>
            <br>
            <b>Date</b>: {{ date || "Unknown" }}
            <span v-if="date !== dateFrontend">(frontend {{ dateFrontend }})</span>
            <br>
            <b>Golang version</b>: {{ golangVersion || "Unknown" }}
            <br>
            <b>Vue.js Version</b>: {{ vuejsVersion || "Unknown" }}
            <br>
            <b>Operating System</b>: {{ osType || "Unknown" }}
            <br>
            <b>Architecture</b>: {{ osArch || "Unknown" }}
          </p>
        </b-message>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/common/common'
import system from '@/requests/authenticated/system'
import constants from '@/constants/constants'
import vue from 'vue'

export default {
  name: 'AboutFlattrack',
  data () {
    return {
      hasInitialLoaded: false,
      version: '',
      commitHash: '',
      mode: '',
      date: '',
      golangVersion: '',
      osType: '',
      osArch: '',
      vuejsVersion: vue.version,
      versionFrontend: constants.appBuildVersion,
      commitHashFrontend: constants.appBuildHash,
      modeFrontend: constants.appBuildMode,
      dateFrontend: constants.appBuildDate
    }
  },
  async beforeMount () {
    this.GetVersion()
  },
  methods: {
    CopyHrefToClipboard () {
      common.CopyHrefToClipboard()
    },
    GetVersion () {
      system.GetVersion().then((resp) => {
        this.hasInitialLoaded = true
        this.version = resp.data.data.version
        this.commitHash = resp.data.data.commitHash
        this.mode = resp.data.data.mode
        this.date = resp.data.data.date
        this.golangVersion = resp.data.data.golangVersion
        this.osType = resp.data.data.osType
        this.osArch = resp.data.data.osArch
      })
    }
  }
}
</script>
