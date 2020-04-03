<template>
  <div>
    <headerDisplay />
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
          <ul>
            <li>
              <a href="/">Home</a>
            </li>
            <li class="is-active">
              <a href="/noticeboard">Noticeboard</a>
            </li>
          </ul>
        </nav>
        <h1 class="title">Noticeboard</h1>
        <h2 class="subtitle">Post to your flatmates</h2>
        <b-button type="is-light" tag="a" href="/noticeboard/p" rounded>Add new post</b-button>
        <br><br>
        <div v-if="posts">
          <div class="card" v-for="post in posts" v-bind:key="post">
            <div class="card-content">
              <div class="media">
                <div class="media-content">
                  <p class="title is-4">{{ post.title }}</p>
                  <p class="subtitle is-6">{{ members[post.author].names || 'Unknown' }}</p>
                </div>
              </div>
              <div class="content">
                {{ post.message }}
                <br>
                <time :datetime="post.timestamp"></time>
              </div>
            </div>
          </div>
          <div class="section">
            <p>{{ posts.length }} {{ posts.length === 1 ? 'post' : 'posts' }}</p>
          </div>
        </div>
    </div>
  </div>
</template>

<script>
import axios from 'axios'
import headerDisplay from '@/components/header-display'

export default {
  name: 'Shopping List',
  data () {
    return {
      posts: [],
      members: {}
    }
  },
  created () {
    axios
      .get(`/api/noticeboard`, {
        headers: {
          Authorization: `Bearer ${localStorage.getItem('authToken')}`
        }
      })
      .then(resp => {
        console.log(resp.data)
        this.posts = resp.data

        // get a list of the entries for this week
        return axios
          .get(`/api/members`, {
            params: {
              all: true
            },
            headers: {
              Authorization: `Bearer ${localStorage.getItem('authToken')}`
            }
          })
      })
      .then(resp => {
        var transformRespData = {}
        resp.data.map(item => {
          transformRespData[item.id] = item
          delete transformRespData[item.id].id
        })
        this.members = transformRespData
        console.log(transformRespData)
      })
      .catch(err => {
        this.$buefy.notification.open({
          duration: 5000,
          message: `An error has occured: ${err}`,
          position: 'is-bottom-right',
          type: 'is-danger',
          hasIcon: true
        })
      })
  },
  components: {
    headerDisplay
  }
}
</script>

<style scoped>

</style>
