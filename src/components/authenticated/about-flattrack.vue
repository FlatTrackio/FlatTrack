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
                <h1 class="title">About FlatTrack</h1>
                <div class="card">
                    <div class="content">
                        Version {{ flattrackVersion }}
                    </div>
                </div>
                <div class="card">
                    <div class="content">
                        FlatTrack is a <a href="https://simple.wikipedia.org/wiki/Free_and_open-source_software">FOSS</a> project which gives your flat some flow with things like tasks and shopping
                    </div>
                </div>
                <div class="card">
                    <div class="content">
                        You can collaborate on this project by visiting <a href="https://gitlab.com/flattrack">FlatTrack's GitLab page</a>.
                    </div>
                </div>
                <div class="card">
                    <div class="content">
                        This program comes with absolutely no warranty.
                        See the <a href="https://www.gnu.org/licenses/gpl-3.0.en.html">GNU General Public License, version 3 or later</a> for details.
                    </div>
                </div>
            </section>
        </div>
    </div>
</template>

<script>
import headerDisplay from '../common/header-display'

export default {
  name: 'About FlatTrack',
  data () {
    return {
      flattrackVersion: '0.0.1'
    }
  },
  components: {
    headerDisplay
  }
}
</script>

<style scoped>

</style>
