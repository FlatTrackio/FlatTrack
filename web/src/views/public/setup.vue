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
        <h2 class="title is-2">Welcome to FlatTrack!</h2>

        <h1 class="title is-1">Set up</h1>
        <p class="subtitle is-5">Let's get started and set up your instance.</p>
        <br />
        <div class="">
          <b-icon icon="cogs" size="is-medium" />
          <h3 class="title is-4">System</h3>
          <b-field label="Language">
            <b-select
              v-model="language"
              placeholder="English"
              autofocus
              required
              icon="web"
              size="is-medium"
              expanded
              @keyup.enter.native="Register"
            >
              <option value="English">English</option>
            </b-select>
          </b-field>
          <b-field label="Timezone">
            <b-select
              v-model="timezone"
              placeholder="Pacific/Auckland"
              autofocus
              required
              icon="map-clock"
              size="is-medium"
              expanded
              @keyup.enter.native="Register"
            >
              <option value="Pacific/Auckland">Pacific/Auckland</option>
            </b-select>
          </b-field>
          <br />
          <b-icon icon="home" size="is-medium" />

          <h3 class="title is-4">Your flat</h3>
          <b-field label="Flat name">
            <b-input
              v-model="flatName"
              type="text"
              maxlength="20"
              placeholder="Enter your flat's name"
              icon="form-textbox"
              size="is-medium"
              icon-right="close-circle"
              icon-right-clickable
              required
              @icon-right-click="flatName = ''"
              @keyup.enter.native="Register"
            />
          </b-field>
          <b-icon icon="account" size="is-medium" />

          <h3 class="title is-4">Your profile</h3>
          <p class="subtitle is-6">
            Note: your account profile will be set up as Administrator
          </p>
          <b-field label="Name(s)">
            <b-input
              v-model="names"
              type="text"
              maxlength="70"
              placeholder="Enter your name(s)"
              icon="form-textbox"
              size="is-medium"
              icon-right="close-circle"
              icon-right-clickable
              required
              @icon-right-click="names = ''"
              @keyup.enter.native="Register"
            />
          </b-field>
          <b-field label="Email">
            <b-input
              v-model="email"
              type="email"
              maxlength="70"
              placeholder="Enter your email address"
              icon="email"
              size="is-medium"
              icon-right="close-circle"
              icon-right-clickable
              required
              @icon-right-click="email = ''"
              @keyup.enter.native="Register"
            />
          </b-field>
          <b-field label="Phone number (optional)">
            <b-input
              v-model="phoneNumber"
              type="tel"
              placeholder="Enter your phone number"
              icon="phone"
              size="is-medium"
              icon-right="close-circle"
              icon-right-clickable
              maxlength="30"
              @icon-right-click="phoneNumber = ''"
              @keyup.enter.native="Register"
            />
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
              trap-focus
              @keyup.enter.native="Register"
            />
          </b-field>
          <br />
          <div class="field has-addons">
            <label class="label">Password</label>
            <p class="control">
              <infotooltip
                message="Make sure that your password has: 10 or more characters, at least one lower case letter, at least one upper case letter, at least one number"
              />
            </p>
          </div>
          <b-field>
            <b-input
              v-model="password"
              type="password"
              password-reveal
              maxlength="70"
              placeholder="Enter a password"
              icon="form-textbox-password"
              size="is-medium"
              pattern="^([a-zA-Z]*).{10,}$"
              validation-message="password is invalid. Make sure that your password has: 10 or more characters, at least one lower case letter, at least one upper case letter, at least one number"
              icon-right="close-circle"
              icon-right-clickable
              required
              @icon-right-click="password = ''"
              @keyup.enter.native="Register"
            />
          </b-field>
          <b-field label="Confirm password">
            <b-input
              v-model="passwordConfirm"
              type="password"
              password-reveal
              placeholder="Confirm your password"
              icon="form-textbox-password"
              size="is-medium"
              maxlength="70"
              pattern="^([a-zA-Z]*).{10,}$"
              validation-message="password is invalid. Make sure that your password has: 10 or more characters, at least one lower case letter, at least one upper case letter, at least one number"
              icon-right="close-circle"
              icon-right-clickable
              required
              @keyup.enter.native="Register"
              @icon-right-click="passwordConfirm = ''"
            />
          </b-field>

          <b-button
            type="is-success"
            size="is-medium"
            icon-left="check"
            native-type="submit"
            expanded
            @click="Register"
          >
            Setup
          </b-button>
          <div
            v-if="typeof message !== 'undefined' && message !== ''"
            class="notification is-warning mb-4 mt-2"
          >
            <p class="subtitle is-6">{{ message }}</p>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
  import registration from "@/requests/public/registration";
  import common from "@/common/common";
  import headerDisplay from "@/components/common/header-display.vue";
  import infotooltip from "@/components/common/info-tooltip.vue";

  export default {
    name: "SetupPage",
    components: {
      headerDisplay,
      infotooltip,
    },
    data() {
      const today = new Date();
      const maxDate = new Date(
        today.getFullYear() - 15,
        today.getMonth(),
        today.getDay()
      );
      const minDate = new Date(
        today.getFullYear() - 100,
        today.getMonth(),
        today.getDay()
      );
      const focusedDate = new Date(
        today.getFullYear() - 15,
        today.getMonth() - 5,
        today.getDay()
      );

      return {
        language: "English",
        timezone: "Pacific/Auckland",
        maxDate: maxDate,
        minDate: minDate,
        focusedDate: focusedDate,
        jsBirthday: null,
        passwordConfirm: "",
        flatName: null,
        message: common.GetSetupMessage(),
        names: null,
        email: null,
        phoneNumber: null,
        birthday: null,
        password: null,
      };
    },
    methods: {
      Register() {
        if (this.password !== this.passwordConfirm) {
          common.DisplayFailureToast(
            this.$buefy,
            "Error passwords do not match"
          );
          return;
        }

        this.birthday = new Date(this.jsBirthday || 0).getTime() / 1000 || 0;

        const loadingComponent = this.$buefy.loading.open({
          container: null,
        });
        setTimeout(() => loadingComponent.close(), 20 * 1000);
        var form = {
          language: this.language,
          timezone: this.timezone,
          flatName: this.flatName,
          user: {
            names: this.names,
            email: this.email,
            password: this.password,
            passwordConfirm: this.passwordConfirm,
            jsBirthday: this.jsBirthday,
            phoneNumber: this.phoneNumber,
          },
        };
        registration
          .PostAdminRegister(form, {
            secret: this.$route.query.secret,
          })
          .then((resp) => {
            common.DisplaySuccessToast(this.$buefy, "Welcome to FlatTrack!");
            setTimeout(() => {
              loadingComponent.close();
              window.location.href = "/";
            }, 3 * 1000);
          })
          .catch((err) => {
            loadingComponent.close();
            common.DisplayFailureToast(
              this.$buefy,
              err.response.data.metadata.response || err
            );
          });
      },
    },
  };
</script>

<style scoped></style>
