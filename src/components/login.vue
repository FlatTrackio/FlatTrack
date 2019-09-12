<template>
    <div>
        <headerDisplay genericMessage="true"/>
        <section class="section">
            <h1 class="title">Login</h1>
            <p class="subtitle">
                Weclome to FlatTrack, please login.
            </p>
            <b-field label="Email">
                <b-input type="email"
                    :value="form.email"
                    maxlength="70">
                </b-input>
            </b-field>
            <b-field label="Password">
                <b-input type="password"
                    :value="form.password"
                    password-reveal
                    maxlength="70">
                </b-input>
            </b-field>
            <b-button rounded>Login</b-button>
            <b-button rounded type="is-warning">Forgot Password</b-button>
        </section>
    </div>
</template>

<script>
import axios from 'axios'
import headerDisplay from './header-display'

export default {
  name: 'login',
  data () {
    return {
      form: {
        email: '',
        password: ''
      }
    }
  },
  components: {
    headerDisplay
  },
  methods: {
    postLogin: () => {
      axios.post('/api/login').then(resp => {
        console.log(resp)
      }).catch(err => {
        console.log(err)
      })
    }
  }
}
</script>

<style>

</style>
