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
              <li><router-link to="/account">My account</router-link></li>
              <li class="is-active"><router-link to="/account/profile">Profile</router-link></li>
            </ul>
        </nav>
        <h1 class="title is-1">Profile</h1>
        <p class="subtitle is-3">Manage your account</p>
        <b-loading :is-full-page="false" :active.sync="pageLoading" :can-cancel="false"></b-loading>
        <div class="card">
          <div class="card-content">
            <div class="media">
              <div class="media-content">
                <figure class="image is-128x128">
                  <img src="@/assets/256x256.png" alt="Placeholder image">
                </figure>
                <br/>
                <p class="title is-3">{{ names }}</p>
                <p class="subtitle is-5">Joined {{ TimestampToCalendar(creationTimestamp) }}</p>
              </div>
            </div>
          </div>
        </div>
        <br />

        <b-field grouped group-multiline>
          <div class="control" v-for="group in groups" v-bind:key="group">
            <b-taglist attached>
              <b-tag type="is-dark">is</b-tag>
              <b-tag type="is-info">{{ group }}</b-tag>
            </b-taglist>
          </div>
        </b-field>
        <br />

        <b-field label="Name(s)">
          <b-input
            type="text"
            v-model="names"
            maxlength="60"
            placeholder="Enter your name(s)"
            icon="textbox"
            size="is-medium"
            @keyup.enter.native="PatchProfile"
            required>
          </b-input>
        </b-field>

        <b-field label="Email">
          <b-input
            type="email"
            v-model="email"
            maxlength="70"
            icon="email"
            size="is-medium"
            placeholder="Enter your email address"
            @keyup.enter.native="PatchProfile"
            required>
          </b-input>
        </b-field>
        <b-field label="Phone number (optional)">
          <b-input
            type="tel"
            v-model="phoneNumber"
            placeholder="Enter your phone number"
            icon="phone"
            size="is-medium"
            @keyup.enter.native="PatchProfile"
            maxlength="30">
          </b-input>
        </b-field>

        <b-field label="Birthday (optional)">
          <b-datepicker
            v-model="jsBirthday"
            :max-date="maxDate"
            :min-date="minDate"
            :show-week-numbers="true"
            :focused-date="focusedDate"
            placeholder="Click to select birthday"
            icon="cake-variant"
            size="is-medium"
            trap-focus>
          </b-datepicker>
        </b-field>
        <br/>
        <b-button
          type="is-success"
          size="is-medium"
          icon-left="delta"
          native-type="submit"
          expanded
          @click="PatchProfile">
          Update profile
        </b-button>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/common/common'
import profile from '@/requests/authenticated/profile'

export default {
  name: 'profile',
  data () {
    const today = new Date()
    const maxDate = new Date(today.getFullYear() - 15, today.getMonth(), today.getDate())
    const minDate = new Date(today.getFullYear() - 100, today.getMonth(), today.getDate())

    return {
      maxDate: maxDate,
      minDate: minDate,
      focusedDate: maxDate,
      passwordConfirm: '',
      jsBirthday: null,
      pageLoading: true,
      names: '',
      email: '',
      phoneNumber: '',
      groups: [],
      password: '',
      birthday: 0,
      creationTimestamp: ''
    }
  },
  methods: {
    GetProfile () {
      profile.GetProfile().then(resp => {
        this.names = resp.data.spec.names
        this.birthday = resp.data.spec.birthday || null
        this.jsBirthday = this.birthday !== null ? new Date(this.birthday * 1000) : null
        this.focusedDate = this.jsBirthday
        this.phoneNumber = resp.data.spec.phoneNumber
        this.groups = resp.data.spec.groups
        this.email = resp.data.spec.email
        this.creationTimestamp = resp.data.spec.creationTimestamp
        this.pageLoading = false
      })
    },
    PatchProfile () {
      if (this.password !== this.passwordConfirm) {
        common.DisplayFailureToast('Unable to use password as they either do not match')
        return
      }
      var birthday = new Date(this.jsBirthday || 0).getTime() / 1000 || 0
      profile.PatchProfile(this.names, this.email, this.phoneNumber, this.birthday, this.password).then(resp => {
        if (resp.data.spec.id === '') {
          common.DisplayFailureToast('Failed to update profile')
          return
        }
        common.DisplaySuccessToast('Successfully updated your profile')
      }).catch(err => {
        common.DisplayFailureToast('Failed to update profile' + '<br/>' + err.response.data.metadata.response)
      })
    },
    TimestampToCalendar (timestamp) {
      return common.TimestampToCalendar(timestamp)
    }
  },
  async beforeMount () {
    this.GetProfile()
  }
}
</script>
