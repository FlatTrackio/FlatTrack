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
    <headerDisplay />
    <div class="container">
      <section class="section form-width">
        <breadcrumb
          back-link-name="Login"
          :current-page-name="$route.name"
        />
        <h1 class="title is-1">
          Forgot Password
        </h1>
        <p class="subtitle is-4">
          Enter you email to reset your password
        </p>
        <b-field
          label="Email"
          class="is-marginless"
        >
          <b-input
            v-model="email"
            type="email"
            maxlength="70"
            autofocus
            placeholder="Enter your email"
            size="is-medium"
            icon="email"
            required
            @keyup.enter.native="sendPasswordResetRequest(email)"
          />
        </b-field>
        <b-button
          expanded
          size="is-medium"
          type="is-primary"
          @click="sendPasswordResetRequest(email)"
        >
          Reset
        </b-button>
      </section>
    </div>
  </div>
</template>

<script>
import headerDisplay from '@/components/common/header-display.vue'
  import breadcrumb from "@/components/common/breadcrumb.vue";

export default {
  name: 'ForgotPassword',
  components: {
    headerDisplay,
      breadcrumb,
  },
  data () {
    return {
      email: ''
    }
  },
  methods: {
    sendPasswordResetRequest: (email) => {
      if (email === '') {
        this.$buefy.toast.open({
          message: 'Please provide a valid email address',
          position: 'is-bottom',
          type: 'is-danger'
        })
      }
      // post to api to request a password reset email
    }
  }
}
</script>

<style scoped>

</style>
