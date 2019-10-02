<template>
    <div>
        <headerDisplay/>
        <div class="container">
          <section class="section">
            <nav class="breadcrumb has-arrow-separator" aria-label="breadcrumbs">
                <ul>
                <li><a href="/#/">Home</a></li>
                <li class="is-active"><a href="/#/aboutflat">About this flat</a></li>
                </ul>
            </nav>
            <h1 class="title">About this flat</h1>
            <div v-if="points && points.length">
              <h2 class="subtitle">Here's a few things you should know</h2>
              <h2 class="title is-4" v-for="point in points" v-bind:key="point">
                - {{ point.line }}
                <div v-if="point.subPoints && point.subPoints.length">
                  <h3 class="subtitle is-5" v-for="subPoint in point.subPoints" v-bind:key="subPoint">
                    &nbsp;&nbsp;&nbsp;- {{ subPoint }}
                  </h3>
                </div>
              </h2>
            </div>
            <div v-else>
              <h2 class="subtitle">Hmmm, it appears no information about your flat has been added yet. Check back later.</h2>
            </div>
          </section>
        </div>
    </div>
</template>

<script>
import axios from 'axios'
import headerDisplay from '../common/header-display'

export default {
  name: 'About this flat',
  data () {
    return {
      points: [
      ],
      pageLocation: location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : '')
    }
  },
  components: {
    headerDisplay
  },
  created () {
    axios.get('/api/flatinfo',
      {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('authToken')}`
        }
      }).then(resp => {
      console.log(resp)

      // remap to add subpoints
      var transformRespData = {}
      resp.data.map(item => {
        if (typeof transformRespData[item.subpointOf] === 'undefined') {
          transformRespData[item.id] = {
            line: item.line,
            subPoints: []
          }
        } else {
          transformRespData[item.subpointOf].subPoints = [...transformRespData[item.subpointOf].subPoints, item.line]
        }
      })

      console.log(transformRespData)

      // remap back into an array
      transformRespData = Array.from(
        Object.keys(transformRespData),
        key => transformRespData[key]
      )

      console.log(transformRespData)

      this.points = transformRespData
    }).catch(err => {
      console.error(err)
    })
  }
}
</script>

<style scoped>

</style>
