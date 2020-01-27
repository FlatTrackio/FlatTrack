<template>
    <div>
        <headerDisplay/>
        <div class="container">
          <section class="section">
            <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
                <ul>
                <li><a href="/">Home</a></li>
                <li class="is-active"><a href="/aboutflat">About this flat</a></li>
                </ul>
            </nav>
            <h1 class="title">About this flat</h1>
            <div v-if="points && points.length">
              <h2 class="subtitle">Here's a few things you should know</h2>
              <div class="card-margin" v-for="point of points" v-bind:key="point">
                <div class="card">
                  <div class="card-content">
                    <div class="content">
                      {{ point.line }}
                    </div>
                    <div v-if="point.subPoints && point.subPoints.length">
                      <div class="card" v-for="subPoint in point.subPoints" v-bind:key="subPoint">
                        <div class="card-content">
                          <div class="content">
                            {{ subPoint }}
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            <div v-else>
              <h2 class="subtitle">Hmmm, it appears no information about your flat has been added yet. Check back later.</h2>
            </div>
          </section>
        </div>
    </div>
</template>

<script>
import headerDisplay from '@/components/header-display'
import { GetAPIflatInfo } from '@/requests/authenticated/flatinfo'

export default {
  name: 'About this flat',
  data () {
    return {
      points: [
      ]
    }
  },
  components: {
    headerDisplay
  },
  created () {
    GetAPIflatInfo()
      .then(resp => {
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
