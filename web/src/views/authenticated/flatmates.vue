<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><router-link to="/apps">Apps</router-link></li>
              <li class="is-active"><router-link to="/apps/flatmates">Flatmates</router-link></li>
            </ul>
        </nav>
        <h1 class="title is-1">Flatmates</h1>
        <p class="subtitle is-3">
          <span v-if="typeof GroupQuery === 'undefined'">
            Get to know your flatmates
          </span>
          <span v-else>
            Listing flatmates, filtering by the group {{ GroupQuery }}
          </span>
        </p>
        <b-loading :is-full-page="false" :active.sync="pageLoading" :can-cancel="false"></b-loading>
        <div v-if="members && members.length > 0">
          <div class="card-margin" v-for="member of members" v-bind:key="member">
            <div class="card">
              <div class="card-content">
                <div class="media">
                  <div class="media-left">
                    <figure class="image is-48x48">
                      <img src="@/assets/96x96.png" alt="Placeholder image">
                    </figure>
                  </div>
                  <div class="media-content">
                    <p class="title is-4">{{ member.names }}</p>
                    <p class="subtitle is-6">Joined {{ TimestampToCalendar(member.creationTimestamp) }}</p>
                  </div>
                </div>
                <div class="content">
                  <b-field grouped group-multiline>
                    <div class="control" v-for="group in member.groups" v-bind:key="group">
                      <b-taglist attached>
                        <b-tag type="is-dark">is</b-tag>
                        <b-tag type="is-info">{{ group }}</b-tag>
                      </b-taglist>
                    </div>
                  </b-field>
                  <p class="subtitle is-6" v-if="member.phoneNumber">
                    Phone: <a :href="`tel:${member.phoneNumber}`">{{ member.phoneNumber }}</a><br/>
                  </p>
                  <p class="subtitle is-6" v-if="member.email">
                    Email: <a :href="`mailto:${member.email}`">{{ member.email }}</a><br/>
                  </p>
                  <a class="subtitle is-6" v-if="member.birthday && member.birthday !== 0">
                    Birthday: {{ TimestampToCalendar(member.birthday) }}<br/>
                  </a>
                  <b-tag type="is-danger" v-if="member.registered !== true">Has not registered</b-tag>
                </div>
              </div>
            </div>
          </div>
          <b-button
            type="is-warning"
            @click="ClearFilter"
            v-if="$route.query.id || $route.query.group"
            expanded>
            Clear filter
          </b-button>
          <div class="section">
            <p>{{ members.length }} {{ members.length === 1 ? 'flatmate' : 'flatmates' }}</p>
          </div>
        </div>
        <div v-if="!members || members.length === 0">
          <div class="card">
            <div class="card-content">
              <div class="media">
                <div class="media-left">
                  <b-icon
                    icon="account-off"
                    size="is-medium">
                  </b-icon>
                </div>
                <div class="media-content">
                  <p class="subtitle is-4" v-if="!pageLoading">No flatmates found.</p>
                  <p class="subtitle is-4" v-else-if="pageLoading">Loading flatmates...</p>
                </div>
              </div>
              <p class="content subtitle is-6">
                Either you haven't added any flatmates, or you are trying to search or filter for flatmates and none could be found.
              </p>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import emoji from 'node-emoji'
import flatmates from '@/requests/authenticated/flatmates'
import common from '@/common/common'

export default {
  name: 'Flatmates',
  data () {
    return {
      pageLoading: true,
      members: [],
      emojiSmile: emoji.get('smile')
    }
  },
  async beforeMount () {
    this.FetchAllFlatmates()
  },
  computed: {
    GroupQuery () {
      return this.$route.query.group
    },
    IdQuery () {
      return this.$route.query.id
    }
  },
  methods: {
    FetchAllFlatmates () {
      if (typeof this.GroupQuery !== 'undefined') {
        var group = this.GroupQuery
      } else if (typeof this.IdQuery !== 'undefined') {
        var id = this.IdQuery
      }
      var notSelf = true
      flatmates.GetAllFlatmates(id, notSelf, group).then(resp => {
        this.pageLoading = false
        this.members = resp.data.list
      }).catch(err => {
        common.DisplayFailureToast('Failed to list flatmates' + `<br/>${err}`)
      })
    },
    ClearFilter () {
      this.$router.replace({ name: 'My Flatmates' })
    },
    TimestampToCalendar (timestamp) {
      return common.TimestampToCalendar(timestamp)
    }
  },
  watch: {
    GroupQuery () {
      this.FetchAllFlatmates()
    }
  }
}
</script>

<style src="../../assets/style.css"></style>

<style scoped>

</style>
