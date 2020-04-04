<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><a href="/">Home</a></li>
              <li><a href="/apps">Apps</a></li>
              <li class="is-active"><a href="/apps/flatmates">Flatmates</a></li>
            </ul>
        </nav>
        <h1 class="title">Flatmates</h1>
        <div v-if="typeof groupQuery === 'undefined'">
          <p>These are your flatmates, make sure to get to know them {{ emojiSmile }}</p>
          <br>
        </div>
        <div v-else="">
          <p>Listing flatmates, filtering by the group {{ groupQuery }}</p>
          <br>
        </div>
        <div v-if="members && members.length">
          <div class="card-margin" v-for="member of members" v-bind:key="member">
            <div class="card">
              <div class="card-content">
                <div class="media">
                  <div class="media-content">
                    <p class="title is-4">{{ member.names }}</p>
                    <b-field grouped group-multiline>
                      <div class="control" v-for="group in member.groups" v-bind:key="group">
                        <b-taglist attached>
                          <b-tag type="is-dark">is</b-tag>
                          <b-tag type="is-info">{{ group }}</b-tag>
                        </b-taglist>
                      </div>
                    </b-field>
                  </div>
                </div>
                <div class="content">
                  <div v-if="member.phoneNumber">
                    Phone: <a :href="`tel:${member.phoneNumber}`">{{ member.phoneNumber }}</a><br/>
                  </div>
                  <div v-if="member.email">
                    Email: <a :href="`mailto:${member.email}`">{{ member.email }}</a><br/>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="section">
            <p>{{ members.length }} {{ members.length === 1 ? 'flatmate' : 'flatmates' }}</p>
          </div>
        </div>
        <div v-if="members && !members.length">
          <div class="card">
            <div class="card-content">
              <div class="media">
                <div class="media-content">
                  <p class="title is-4">No flatmates added</p>
                </div>
              </div>
              <div class="content">
                Hmm, it appears that you don't have an flatmates added.<br/>
                Please contact the administrator(s) to add your flatmates.
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import emoji from 'node-emoji'
import flatmates from '@/frontend/requests/authenticated/flatmates'
import common from '@/frontend/common/common'

export default {
  name: 'Flatmates',
  data () {
    return {
      members: [],
      groupQuery: undefined,
      emojiSmile: emoji.get('smile')
    }
  },
  created () {
    this.groupQuery = this.$route.query.group
    this.FetchAllFlatmates()
  },
  methods: {
    FetchAllFlatmates () {
      var params = {}
      if (typeof groupQuery !== 'undefined') {
        params.group = groupQuery
      }
      flatmates.GetAllFlatmates(params).then(resp => {
        this.members = resp.data.list
      }).catch(err => {
        common.DisplayFailureToast('Failed to list flatmates' + `<br/>${err}`)
      })
    }
  }
}
</script>

<style src="../../assets/style.css"></style>

<style scoped>

</style>
