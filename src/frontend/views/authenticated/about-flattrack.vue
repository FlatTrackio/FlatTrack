<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><router-link :to="{ name: 'Home' }">Home</router-link></li>
              <li class="is-active"><router-link :to="{ name: 'About FlatTrack' }">About FlatTrack</router-link></li>
            </ul>
        </nav>
        <h1 class="title is-1">About FlatTrack</h1>
        <p class="subtitle is-3">What is FlatTrack?</p>
        <b-message type="is-primary" has-icon icon="help">
          <p class="is-size-5">
            FlatTrack is a
            <a href="https://simple.wikipedia.org/wiki/Free_and_open-source_software" target="_blank" rel="noreferrer">Free and Open Source</a>
            collaboration software for flats / community houses and homes
            with the goals of <b>easing common tasks</b> in living environments, <b>enabling closer collaboration</b>, and <b>empowering humans who live together</b>.
          </p>
        </b-message>
        <p class="subtitle is-3">This FlatTrack instance</p>
        <b-message type="is-warning" has-icon icon="information-outline">
          <p class="is-size-5">
            <b>Version</b>: {{ version || 'Unknown' }}
            <span v-if="version !== versionFrontend">(frontend {{ versionFrontend }})</span>
            <br/>
            <b>Commit hash</b>:
            <a v-if="commitHash !== '???' && typeof commitHash !== 'undefined'" :href="'https://gitlab.com/flattrack/flattrack/-/commit/' + commitHash" target="_blank" rel="noreferrer">{{ commitHash }}</a>
            <span v-else>
              Unknown
            </span>
            <a v-if="commitHash !== commitHashFrontend" :href="'https://gitlab.com/flattrack/flattrack/-/commit/' + commitHashFrontend" target="_blank" rel="noreferrer">{{ commitHashFrontend }}</a>
            <br/>
            <b>Mode</b>: {{ mode || 'Unknown' }}
            <span v-if="mode !== modeFrontend">(frontend {{ modeFrontend }})</span>
            <br/>
            <b>Date</b>: {{ date || 'Unknown' }}
            <span v-if="date !== dateFrontend">(frontend {{ dateFrontend }})</span>
            <br/>
            <b>Golang version</b>: {{ golangVersion || 'Unknown' }}
            <br/>
            <b>Vue.js Version</b>: {{ vuejsVersion || 'Unknown' }}
          </p>
        </b-message>
      </section>
    </div>
  </div>
</template>

<script>
import system from '@/frontend/requests/authenticated/system'
import common from '@/frontend/common/common'
import constants from '@/frontend/constants/constants'
import vue from 'vue'

export default {
  name: 'flat',
  data () {
    return {
      hasInitialLoaded: false,
      version: '',
      commitHash: '',
      mode: '',
      date: '',
      golangVersion: '',
      vuejsVersion: vue.version,
      versionFrontend: constants.appBuildVersion,
      commitHashFrontend: constants.appBuildHash,
      modeFrontend: constants.appBuildMode,
      dateFrontend: constants.appBuildDate
    }
  },
  methods: {
    GetVersion () {
      system.GetVersion().then(resp => {
        this.hasInitialLoaded = true
        this.version = resp.data.data.version
        this.commitHash = resp.data.data.commitHash
        this.mode = resp.data.data.mode
        this.date = resp.data.data.date
        this.golangVersion = resp.data.data.golangVersion
      })
    }
  },
  async beforeMount () {
    this.GetVersion()
  }
}
</script>
