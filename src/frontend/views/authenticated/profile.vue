<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><a href="/">Home</a></li>
              <li class="is-active"><a href="/profile">Profile</a></li>
            </ul>
        </nav>
        <div class="card">
          <div class="card-content">
            <div class="media">
              <div class="media-left">
                <figure class="image is-48x48">
                  <img src="https://bulma.io/images/placeholders/96x96.png" alt="Placeholder image">
                </figure>
              </div>
              <div class="media-content">
                <p class="title is-4">{{ names }}</p>
                <p class="subtitle is-6">Subtitle</p>
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

        <b-field label="Email">
          <b-input type="email"
                   v-model="email"
                   maxlength="70"
                   required>
          </b-input>
        </b-field>
        <b-field label="Phone number">
          <b-input type="tel"
                   v-model="phoneNumber"
                   maxlength="30"
                   >
          </b-input>
        </b-field>
        <b-field label="Password">
          <b-input type="password"
                   v-model="password"
                   password-reveal
                   maxlength="70"
                   >
          </b-input>
        </b-field>
        <b-field label="Confirm password">
          <b-input type="password"
                   v-model="passwordConfirm"
                   password-reveal
                   maxlength="70"
                   >
          </b-input>
        </b-field>
        <!-- <b-button type="is-success" size="is-medium" rounded native-type="submit" @click="Register({ language, timezone, flatName, user: { names, email, password } })">Save</b-button> -->
      </section>
    </div>
  </div>
</template>

<script>
import profile from '@/frontend/requests/authenticated/profile'

export default {
  name: 'profile',
  data () {
    return {
      names: '',
      email: '',
      phoneNumber: '',
      groups: [],
      password: ''
    }
  },
  methods: {
    GetProfile () {
      profile.GetProfile().then(resp => {
        this.names = resp.data.spec.names
        this.groups = resp.data.spec.groups
        this.email = resp.data.spec.email
      })
    }
  },
  async created () {
    this.GetProfile()
  }
}
</script>
