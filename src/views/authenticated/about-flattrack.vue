<template>
    <div>
        <headerDisplay/>
        <div class="container">
            <section class="section">
                <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
                    <ul>
                    <li><a href="/">Home</a></li>
                    <li class="is-active"><a href="/aboutflattrack">About FlatTrack</a></li>
                    </ul>
                </nav>
                <h1 class="title is-2">About FlatTrack</h1>
                <p class="subtitle is-4">What's to know?</p>
                <div class="card">
                    <div class="card-content">
                        <p class="subtitle">
                            Version
                        </p>
                        <p class="title">
                            {{ flattrackVersion }}
                        </p>
                    </div>
                </div>
                <div class="card">
                    <div class="card-content">
                        <p class="subtitle">
                            What is <a href="https://flattrack.io">FlatTrack</a>?
                        </p>
                        <p class="title">
                            A solution to keep many areas of your flat or community house organised
                        </p>
                    </div>
                </div>
                <div class="card">
                    <div class="card-content">
                        <p class="subtitle">
                            Collaborate
                        </p>
                        <p class="title">
                            You can collaborate on this project by visiting <a href="https://gitlab.com/flattrack">FlatTrack's GitLab page</a>.
                        </p>
                    </div>
                </div>
                <div class="card">
                    <div class="card-content">
                        <p class="subtitle">
                            Free and Open Source
                        </p>
                        <p class="title">
                            FlatTrack is a <a href="https://simple.wikipedia.org/wiki/Free_and_open-source_software">FOSS</a> project which means that you can freely study it, modify it, and share the changes.
                        </p>
                    </div>
                </div>
                <div class="card">
                    <div class="card-content">
                        <p class="subtitle">
                            License
                        </p>
                        <p class="title">
                            This program comes with absolutely no warranty.
                            See the <a href="https://www.gnu.org/licenses/gpl-3.0.en.html">GNU General Public License, version 3 or later</a> for details.
                        </p>
                    </div>
                </div>
            </section>
        </div>
    </div>
</template>

<script>
import axios from 'axios'
import { Service } from 'axios-middleware'
import headerDisplay from '@/components/header-display'

const service = new Service(axios)
service.register({
  onResponse (response) {
    if (response.status === 403) {
      localStorage.removeItem('authToken')
      location.href = '/'
    }
    return response
  }
})

export default {
  name: 'About FlatTrack',
  data () {
    return {
      flattrackVersion: ''
    }
  },
  components: {
    headerDisplay
  },
  created () {
    axios({
      method: 'get',
      url: `/api/meta`,
      headers: {
        Authorization: `Bearer ${localStorage.getItem('authToken')}`
      }}).then(resp => {
      this.flattrackVersion = resp.data.version
    }).catch(err => {
      this.pageErrors.push(err)
    })
  }
}
</script>

<style scoped>

</style>
