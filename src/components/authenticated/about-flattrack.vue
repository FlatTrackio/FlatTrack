<template>
    <div>
        <headerDisplay/>
        <div class="container">
            <section class="section">
                <nav class="breadcrumb has-arrow-separator" aria-label="breadcrumbs">
                    <ul>
                    <li><a href="/#/">Home</a></li>
                    <li class="is-active"><a href="/#/aboutflattrack">About FlatTrack</a></li>
                    </ul>
                </nav>
                <h1 class="title is-2">About FlatTrack</h1>
                <p class="subtitle is-4">What's to know?</p>
                <div class="card-margin">
                    <div class="card">
                        <div class="card-content">
                            <div class="content">
                                Version {{ flattrackVersion }}
                            </div>
                        </div>
                    </div>
                </div>
                <div class="card-margin">
                    <div class="card">
                        <div class="card-content">
                            <div class="content">
                                FlatTrack is a <a href="https://simple.wikipedia.org/wiki/Free_and_open-source_software">FOSS</a> project which gives your flat some flow with things like tasks and shopping
                            </div>
                        </div>
                    </div>
                </div>
                <div class="card-margin">
                    <div class="card">
                        <div class="card-content">
                            <div class="content">
                                You can collaborate on this project by visiting <a href="https://gitlab.com/flattrack">FlatTrack's GitLab page</a>.
                            </div>
                        </div>
                    </div>
                </div>
                <div class="card-margin">
                    <div class="card">
                        <div class="card-content">
                            <div class="content">
                                This program comes with absolutely no warranty.
                                See the <a href="https://www.gnu.org/licenses/gpl-3.0.en.html">GNU General Public License, version 3 or later</a> for details.
                            </div>
                        </div>
                    </div>
                </div>
            </section>
        </div>
    </div>
</template>

<script>
import axios from 'axios'
import headerDisplay from '../common/header-display'

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
      }}).then(response => {
      this.flattrackVersion = response.data.version
    }).catch(err => {
      this.pageErrors.push(err)
    })
  }
}
</script>

<style scoped>

</style>
