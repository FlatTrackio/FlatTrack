<template>
    <div>
        <headerDisplay genericMessage="true"/>
        <div class="container">
            <section class="section is-large">
                <h2 class="title is-2">Unknown page!</h2>
                <h4 class="title is-4">Hmmm maybe something went wrong :/</h4>
                <b-button tag="a" href="/#/" rounded>Go home</b-button>
            </section>
        </div>
        <footer class="footer">
            <div class="content has-text-centered">
                <div class="columns is-variable is-1-mobile is-0-tablet is-3-desktop is-8-widescreen is-2-fullhd">
                    <div class="column">
                        <a href="/#/">Go home</a>
                    </div>
                </div>
            </div>
        </footer>
    </div>
</template>

<script>
import headerDisplay from '../common/header-display'

export default {
  name: 'Unknown page',
  pageLocation: location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : ''),
  data () {
    return {
    }
  },
  components: {
    headerDisplay
  }
}
</script>

<style scoped>

</style>
