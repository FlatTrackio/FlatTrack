<template>
    <div>
      <headerDisplay/>
      <div class="container">
        <section class="section">
          <h1 class="title is-1">Set up FlatTrack</h1>
          <p class="subtitle is-5">
              Let's get started
          </p>
          <br/>
          <h3 class="title is-4">System</h3>
          <b-field label="Language">
            <b-select
              placeholder="English"
              v-model="language"
              autofocus
              required
              expanded>
              <option value="English">English</option>
            </b-select>
          </b-field>
          <b-field label="Timezone">
            <b-select
              placeholder="Pacific/Auckland"
              v-model="timezone"
              autofocus
              required
              expanded>
              <option value="Pacific/Auckland">Pacific/Auckland</option>
            </b-select>
          </b-field>
          <b-field label="Flat name">
              <b-input type="text"
                  v-model="flatName"
                  maxlength="20"
                  required>
              </b-input>
          </b-field>
          <h3 class="title is-4">Your account</h3>
          <b-field label="Your name">
              <b-input type="text"
                  v-model="names"
                  maxlength="70"
                  required>
              </b-input>
          </b-field>
          <b-field label="Email">
              <b-input type="email"
                  v-model="email"
                  maxlength="70"
                  required>
              </b-input>
          </b-field>
          <b-field label="Password">
              <b-input type="password"
                  v-model="password"
                  password-reveal
                  maxlength="70"
                  @keyup.enter.native="Register({ language, timezone, flatName, user: { names, email, password } })"
                  required>
              </b-input>
          </b-field>
          <b-button type="is-success" size="is-medium" rounded native-type="submit" @click="Register({ language, timezone, flatName, user: { names, email, password } })">Setup</b-button>
        </section>
      </div>
    </div>
</template>

<script>
import headerDisplay from '@/frontend/components/header-display.vue'
import registration from '@/frontend/requests/public/registration'
import { NotificationProgrammatic as Notification } from 'buefy'

export default {
  name: 'setup',
  data () {
    return {
      language: 'English',
      timezone: 'Pacific/Auckland',
      flatName: '',
      names: '',
      email: '',
      password: ''
    }
  },
  components: {
    headerDisplay
  },
  methods: {
    Register: (form) => {
      registration.PostAdminRegister(form).then(resp => {
        console.log(resp)
        Notification.open({
          duration: 8 * 1000,
          message: `Welcome to FlatTrack! ${resp.data.metadata.message}`,
          position: 'is-bottom-right',
          type: 'is-light',
          hasIcon: true
        })
        setTimeout(() => {
          window.location.href = '/'
        }, 8 * 1000)
      }).catch(err => {
        Notification.open({
          duration: 8 * 1000,
          message: `${err}`,
          position: 'is-bottom-right',
          type: 'is-danger',
          hasIcon: true
        })
      })
    }
  }
}
</script>

<style scoped>

</style>
