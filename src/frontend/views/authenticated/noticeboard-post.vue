<template>
    <div>
      <headerDisplay/>
      <div class="container">
        <section class="section">
          <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><a href="/">Home</a></li>
              <li><a href="/noticeboard">Noticeboard post</a></li>
              <li class="is-active"><a>{{ returnNameforID(id, name) || 'Add a post' }}</a></li>
            </ul>
          </nav>
          <h1 class="title">{{ returnNameforID(id, name) || 'Add a post' }}</h1>
          <div class="card">
            <div class="card-content">
              <div class="content">
                <b-field label="Title" v-if="!id">
                  <b-input placeholder="Title" v-model="title" maxlength="30" rounded required></b-input>
                </b-field>
                <b-field label="Message">
                  <b-input placeholder="What would you like to say?" v-model="message" maxlength="300" type="textarea" rounded required></b-input>
                </b-field>
               </div>
              <div v-if="returnNameforID(id, name)">
                <b-button type="is-success" native-type="submit" @click="updatePost(id, title, message)">Update</b-button>
                <b-button type="is-danger" @click="deletePost(id)">Delete</b-button>
              </div>
              <div v-else>
                <b-button type="is-success" native-type="submit" @click="addNewPost(title, message)">Post</b-button>
              </div>
            </div>
          </div>
        </section>
      </div>
    </div>
</template>

<script>
import { ToastProgrammatic as Toast, DialogProgrammatic as Dialog } from 'buefy'
import headerDisplay from '@/frontend/components/header-display'
import { GetAPInoticeboardEntryById, PostAPInoticeboardEntry, PutAPInoticeboardEntry, DeleteAPInoticeboardEntry } from '@/frontend/requests/authenticated/noticeboard'

export default {
  name: 'Noticeboard post',
  data () {
    return {
      id: this.$route.query.id,
      title: '',
      message: '',
      pageErrors: []
    }
  },
  created () {
    var id = this.$route.query.id
    GetAPInoticeboardEntryById(id)
      .then(resp => {
        this.post = resp.data
      })
      .catch(err => {
        this.pageErrors = [...this.pageErrors, err]
      })
  },
  methods: {
    addNewPost: (title, message) => {
      var post = {
        title,
        message
      }
      PostAPInoticeboardEntry(post)
        .then(resp => {
          Toast.open({
            message: 'Posted successfully',
            position: 'is-bottom',
            type: 'is-success'
          })
          setTimeout(() => {
            window.window.location.href = '/noticeboard'
          }, 1000)
        })
        .catch(err => {
          Toast.open({
            message: `Failed to add post: ${err}`,
            position: 'is-bottom',
            type: 'is-danger'
          })
        })
    },
    updatePost: (id, title, message) => {
      var post = {
        title,
        message
      }
      PutAPInoticeboardEntry(id, post)
        .then(resp => {
          Toast.open({
            message: 'Posted successfully',
            position: 'is-bottom',
            type: 'is-success'
          })
          setTimeout(() => {
            window.location.reload()
          }, 1000)
        })
        .catch(err => {
          Toast.open({
            message: `Failed to update post: ${err}`,
            position: 'is-bottom',
            type: 'is-danger'
          })
        })
    },
    deletePost: (id) => {
      Dialog.confirm({
        message: 'Are you sure you want to remove this post? (this cannot be undone)',
        confirmText: 'Remove post',
        type: 'is-danger',
        onConfirm: () => {
          DeleteAPInoticeboardEntry(id)
            .then(resp => {
              Toast.open({
                message: 'Post removed successfully',
                position: 'is-bottom',
                type: 'is-success'
              })
              setTimeout(() => {
                window.window.location.href = '/noticeboard'
              }, 1.2 * 1000)
            })
            .catch(err => {
              Toast.open({
                message: `Failed to remove post: ${err}`,
                position: 'is-bottom',
                type: 'is-danger'
              })
            })
        }
      })
    },
    returnNameforID: (id, name) => {
      if (typeof id === 'undefined') {
        return null
      }
      return name
    }
  },
  components: {
    headerDisplay
  }
}
</script>

<style scoped>

</style>
