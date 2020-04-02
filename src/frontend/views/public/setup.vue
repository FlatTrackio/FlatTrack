<template>
<div>
  <headerDisplay/>
  <div class="container">
    <section class="section form-width">
      <h1 class="title is-1">Set up FlatTrack</h1>
      <p class="subtitle is-5">
        Let's get started
      </p>
      <br/>
      <div class="form-width">
        <!-- TODO add system icon -->
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
        <!-- TODO add account icon -->
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
      </div>
    </section>
  </div>
</div>
</template>

<script>
import headerDisplay from '@/frontend/components/header-display.vue'
import registration from '@/frontend/requests/public/registration'
import { ToastProgrammatic as Toast, LoadingProgrammatic as Loading } from 'buefy'

export default {
  name: 'setup',
  data () {
    return {
      language: 'English',
      timezone: 'Pacific/Auckland',
      flatName: null,
      names: null,
      email: null,
      password: null
    }
  },
  components: {
    headerDisplay
  },
  methods: {
    Register: (form) => {
      const loadingComponent = Loading.open({
        container: null
      })
      setTimeout(() => loadingComponent.close(), 20 * 1000)
      registration.PostAdminRegister(form).then(resp => {
        if (resp.data.data !== '' || typeof resp.data.data !== 'undefined') {
          localStorage.setItem('authToken', resp.data.data)
        } else {
          Error('failed to find authToken')
        }
        Toast.open({
          duration: 8 * 1000,
          message: `Welcome to FlatTrack!<br/>${resp.data.metadata.response}`,
          position: 'is-top',
          type: 'is-success',
          hasIcon: true
        })
        setTimeout(() => {
          loadingComponent.close()
          window.location.href = '/'
        }, 3 * 1000)
      }).catch(err => {
        loadingComponent.close()
        Toast.open({
          duration: 8 * 1000,
          message: `${err}`,
          position: 'is-top',
          type: 'is-danger',
          hasIcon: true
        })
      })
    }
  }
}
</script>

<style scoped>
.form-width {
    width: 380px;
    margin: auto;
}
</style>
