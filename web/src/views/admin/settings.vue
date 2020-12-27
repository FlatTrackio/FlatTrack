<!--
     This program is free software: you can redistribute it and/or modify
     it under the terms of the Affero GNU General Public License as published by
     the Free Software Foundation, either version 3 of the License, or
     (at your option) any later version.

     This program is distributed in the hope that it will be useful,
     but WITHOUT ANY WARRANTY; without even the implied warranty of
     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
     GNU General Public License for more details.

     You should have received a copy of the Affero GNU General Public License
     along with this program.  If not, see <https://www.gnu.org/licenses/>.
-->

<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><router-link to="/admin">Admin</router-link></li>
              <li class="is-active"><router-link to="/admin/settings">Settings</router-link></li>
            </ul>
        </nav>
        <h1 class="title is-1">Settings</h1>
        <p class="subtitle is-4">General FlatTrack settings</p>
        <b-loading :is-full-page="false" :active.sync="pageLoading" :can-cancel="false"></b-loading>
        <div>
          <label class="label">Flat name</label>
          <b-field>
            <b-input
              type="text"
              v-model="flatName"
              maxlength="20"
              placeholder="Enter your flat's name"
              icon="textbox"
              size="is-medium"
              icon-right="close-circle"
              icon-right-clickable
              @icon-right-click="flatName = ''"
              @keyup.enter.native="PostFlatName"
              expanded
              required>
            </b-input>
            <p class="control">
              <b-button
                type="is-primary"
                size="is-medium"
                icon-left="check"
                @click="PostFlatName">
              </b-button>
            </p>
          </b-field>
        </div>
        <div>
          <label class="label">Flat notes</label>
          <b-field>
            <b-input
              type="textarea"
              v-model="flatNotes"
              minlength="0"
              maxlength="500"
              placeholder="Enter notes about your flat. e.g: living space, rules, obligations, etc..."
              size="is-medium"
              icon-right="close-circle"
              icon-right-clickable
              @icon-right-click="flatNotes = ''"
              expanded>
            </b-input>
            <p class="control">
              <b-button
                type="is-primary"
                size="is-medium"
                icon-left="check"
                @click="PutFlatNotes">
              </b-button>
            </p>
          </b-field>
        </div>
        <br/>
      </section>
    </div>
  </div>
</template>

<script>
import flatInfo from '@/requests/authenticated/flatInfo'
import settings from '@/requests/admin/settings'
import common from '@/common/common'

export default {
  name: 'Flatmates',
  data () {
    return {
      pageLoading: true,
      flatName: '',
      flatNotes: ''
    }
  },
  async beforeMount () {
    flatInfo.GetFlatName().then(resp => {
      this.flatName = resp.data.spec
      return flatInfo.GetFlatNotes()
    }).then(resp => {
      this.flatNotes = resp.data.spec.notes
      this.pageLoading = false
    })
  },
  methods: {
    goToRef (ref) {
      this.$router.push({ path: ref })
    },
    PostFlatName () {
      if (this.flatName === '') {
        common.DisplayFailureToast('Error: Flat name must not be empty')
        return
      }
      settings.PostFlatName(this.flatName).then(resp => {
        common.DisplaySuccessToast('Set flat name')
      }).catch(err => {
        common.DisplayFailureToast('Failed set the flat name' + '<br/>' + (err.response.data.metadata.response || err))
      })
    },
    PutFlatNotes () {
      if (this.flatName === '') {
        common.DisplayFailureToast('Error: Flat notes must not be empty')
        return
      }
      settings.PutFlatNotes(this.flatNotes).then(resp => {
        common.DisplaySuccessToast('Set flat notes')
      }).catch(err => {
        common.DisplayFailureToast('Failed set the flat notes' + '<br/>' + (err))
      })
    },
    TimestampToCalendar (timestamp) {
      return common.TimestampToCalendar(timestamp)
    }
  }
}
</script>
