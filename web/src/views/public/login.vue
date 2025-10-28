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
        <h1 class="title is-1">
          Login
        </h1>
        <p class="subtitle is-4">
          Welcome to FlatTrack, please login.
        </p>
        <b-field
          label="Email"
          class="is-marginless"
        >
          <b-input
            v-model="email"
            name="email"
            type="email"
            maxlength="70"
            autofocus
            placeholder="Enter your email"
            size="is-medium"
            icon="email"
            icon-right="close-circle"
            icon-right-clickable
            required
            @keyup.enter.native="postLogin"
            @icon-right-click="email = ''"
          />
        </b-field>
        <b-field
          label="Password"
          class="is-marginless"
        >
          <b-input
            v-model="password"
            name="password"
            type="password"
            password-reveal
            maxlength="70"
            placeholder="Enter your password"
            size="is-medium"
            icon="form-textbox-password"
            pattern="^([a-zA-Z]*).{10,}$"
            icon-right="close-circle"
            icon-right-clickable
            required
            @keyup.enter.native="postLogin"
            @icon-right-click="password = ''"
          />
        </b-field>
        <div class="field">
          <p class="control">
            <b-button
              icon-left="login"
              native-type="submit"
              size="is-medium"
              type="is-primary"
              expanded
              @click="postLogin"
            >
              Login
            </b-button>
            <b-button
              tag="a"
              href="forgot-password"
              icon-left="lifebuoy"
              size="is-medium"
              expanded
              disabled
              type="is-text"
            >
              Forgot Password
            </b-button>
          </p>
          <div
            v-if="typeof message !== 'undefined' && message !== ''"
            class="notification is-warning mb-4 mt-2"
          >
            <p class="subtitle is-6">
              {{ message }}
            </p>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
  import login from "@/requests/public/login";
  import headerDisplay from "@/components/common/header-display.vue";
  import common from "@/common/common";

  export default {
    name: "LoginPage",
    components: {
      headerDisplay,
    },
    data() {
      return {
        redirect: this.$route.query.redirect || undefined,
        authToken: this.$route.query.authToken || undefined,
        message: common.GetLoginMessage() || undefined,
        email: "",
        password: "",
      };
    },
    mounted() {
      if (typeof this.$route.query.authToken !== "undefined") {
        common.SetAuthToken(this.$route.query.authToken);
        this.$router.push({ name: "Home" });
        return;
      }
      this.checkForLoginToken();
    },
    methods: {
      postLogin() {
        const loadingComponent = this.$buefy.loading.open({
          container: null,
        });
        setTimeout(() => loadingComponent.close(), 20 * 1000);
        login
          .PostUserAuth(this.email, this.password)
          .then((resp) => {
            common.SetAuthToken(resp.data.data);
            setTimeout(() => {
              loadingComponent.close();
              if (typeof this.redirect !== "undefined" && this.redirect) {
                this.$router.push({ path: this.redirect });
                return;
              }
              window.location.href = "/";
            }, 1 * 1000);
          })
          .catch((err) => {
            loadingComponent.close();
            common.DisplayFailureToast(
              this.$buefy,
              err.response.data.metadata.response || err
            );
          });
      },
      checkForLoginToken() {
        var authToken = common.GetAuthToken();
        if (
          !(
            typeof authToken === "undefined" ||
            authToken === null ||
            authToken === ""
          )
        ) {
          login.GetUserAuth(false).then((res) => {
            // verify token via request or something
            const loadingComponent = this.$buefy.loading.open({
              container: null,
            });
            common.DisplaySuccessToast(
              this.$buefy,
              "You are still signed in, going to the home page..."
            );
            setTimeout(() => {
              loadingComponent.close();
              if (
                this.redirect !== null &&
                typeof this.redirect !== "undefined"
              ) {
                this.$router.push({ path: this.redirect });
                return;
              } else {
                this.$router.push({ name: "Home" });
              }
            }, 1 * 1000);
          });
        }
      },
    },
  };
</script>

<style scoped></style>
