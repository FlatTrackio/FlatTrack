<template>
    <div>
        <section class="hero is-info">
        <div class="hero-body">
          <p class="title">
            FlatTrack
          </p>
          <p class="subtitle">
            {{ deploymentName }}
          </p>
        </div>
        </section>
        <div class="container">
          <section class="section">
            <nav class="breadcrumb has-arrow-separator" aria-label="breadcrumbs">
                <ul>
                <li><a href="/#/">Home</a></li>
                <li class="is-active"><a href="/#/aboutflat">About this flat</a></li>
                </ul>
            </nav>
            <h1 class="title">About this flat</h1>
            <h2 class="subtitle">Here's a few things you should know</h2>
          </section>
          <section class="section">
            <ul class="menu-list">
                <li class="menu-label" v-for="point in points" v-bind:key="point">
                    {{ point.topPoint }}
                    <div v-if="point.subPoints && point.subPoints.length">
                        <ul class="menu-list">
                            <li class="menu-label" v-for="subPoint in point.subPoints" v-bind:key="subPoint">
                                {{ subPoint }}
                            </li>
                        </ul>
                    </div>
                </li>
            </ul>
          </section>
        </div>
    </div>
</template>

<script>
export default {
  name: 'About this flat',
  data () {
    return {
      points: [
        {
          topPoint: 'This is the first point'
        },
        {
          topPoint: 'Here\'s the second point'
        },
        {
          topPoint: 'Ah yes... the third point, how lovely -- something to admire'
        },
        {
          topPoint: 'fourth point, we can skip this one'
        },
        {
          topPoint: 'the fifth and final point',
          subPoints:
          [
            'this is a subpoint of the toppoint, to highlight a thing about it',
            'but also this point is relevant, so it should be included'
          ]
        }
      ],
      deploymentName: 'Keep track of your flat',
      pageLocation: location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : ''),
    }
  },
  created () {
    axios.get(`${location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : '')}/api/settings/deploymentName`)
      .then(response => {
        this.deploymentName = response.data.value
      })
      .catch(err => {
        this.pageErrors.push(err)
      })
  }
}
</script>

<style scoped>

</style>
