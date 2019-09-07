<template>
    <div>
        <section class="hero is-dark">
            <div class="hero-body">
                <p class="title">
                    About FlatTrack
                </p>
            </div>
        </section>
        <div class="container">
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
        </div>
    </div>
</template>

<script>
export default {
  name: 'About FlatTrack',
  data () {
    return {
      flattrackVersion: '0.0.1'
    }
  }
}
</script>

<style>

</style>
