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
      <section class="section">
        <h1 class="title is-1">
          Confirm your account
        </h1>
        <p class="subtitle is-4">
          Final things to complete your sign up
        </p>

        <b-message
          v-if="
            (typeof id === 'undefined' || id === '' || idValid !== true) &&
              pageHasLoaded === true
          "
          type="is-danger"
          has-icon
        >
          Token Id is missing or is invalid.
          <br>
          <br>
          The link that you were provided appears to be damaged or invalid.
          <br>
          Please contact your system administrator or support.
        </b-message>
        <b-message
          v-if="typeof secret === 'undefined' || secret === ''"
          type="is-danger"
          has-icon
        >
          Missing a confirmation secret.
          <br>
          <br>
          The link that you were provided appears to be damaged or invalid.
          <br>
          Please contact your system administrator or support.
        </b-message>
        <div
          v-if="
            idValid === true && typeof secret !== 'undefined' && secret !== ''
          "
        >
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
              @keyup.enter.native="PostUserConfirm"
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
            />
          </b-field>
          <br>

          <div class="field has-addons is-marginless">
            <h1 class="title is-6 is-marginless">
              Password
            </h1>
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
              placeholder="Enter your password"
              icon="form-textbox-password"
              size="is-medium"
              pattern="^([a-zA-Z]*).{10,}$"
              validation-message="password is invalid. Make sure that your password has: 10 or more characters, at least one lower case letter, at least one upper case letter, at least one number"
              icon-right="close-circle"
              icon-right-clickable
              required
              @icon-right-click="password = ''"
              @keyup.enter.native="PostUserConfirm"
            />
          </b-field>

          <b-field label="Confirm password">
            <b-input
              v-model="passwordConfirm"
              type="password"
              password-reveal
              maxlength="70"
              placeholder="Confirm your password"
              icon="form-textbox-password"
              size="is-medium"
              pattern="^([a-zA-Z]*).{10,}$"
              validation-message="password is invalid. Make sure that your password has: 10 or more characters, at least one lower case letter, at least one upper case letter, at least one number"
              icon-right="close-circle"
              icon-right-clickable
              required
              @icon-right-click="passwordConfirm = ''"
              @keyup.enter.native="PostUserConfirm"
            />
          </b-field>

          <b-button
            type="is-success"
            size="is-medium"
            icon-left="check"
            native-type="submit"
            expanded
            @click="PostUserConfirm"
          >
            Complete
          </b-button>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
  import common from "@/common/common";
  import headerDisplay from "@/components/common/header-display.vue";
  import infotooltip from "@/components/common/info-tooltip.vue";
  import confirm from "@/requests/public/useraccountconfirm";

  export default {
    name: "AccountConfirmation",
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
        today.getDate()
      );

      return {
        minDate: minDate,
        maxDate: maxDate,
        focusedDate: maxDate,
        idValid: false,
        jsBirthday: null,
        pageHasLoaded: false,
        id: this.$route.params.id,
        secret: this.$route.query.secret,
        phoneNumber: null,
        password: null,
        passwordConfirm: null,
      };
    },
    computed: {
      birthday() {
        return new Date(this.jsBirthday || 0).getTime() / 1000 || 0;
      },
    },
    async beforeMount() {
      confirm.GetTokenValid(this.id).then((resp) => {
        this.idValid = resp.data.data;
        this.pageHasLoaded = true;
      });
    },
    methods: {
      PostUserConfirm() {
        if (this.password !== this.passwordConfirm && this.password !== "") {
          common.DisplayFailureToast(this.$buefy, "Passwords do not match");
          return;
        }
        const loadingComponent = this.$buefy.loading.open({
          container: null,
        });
        confirm
          .PostUserConfirm(
            this.id,
            this.secret,
            this.phoneNumber,
            this.birthday,
            this.password
          )
          .then((resp) => {
            if (resp.data.data === "") {
              common.DisplayFailureToast(
                this.$buefy,
                resp.data.metadata.response
              );
              return;
            }
            localStorage.setItem("authToken", resp.data.data);
            common.DisplaySuccessToast(
              this.$buefy,
              resp.data.metadata.response
            );
            setTimeout(() => {
              loadingComponent.close();
              window.location.href = "/";
            }, 2 * 1000);
          })
          .catch((err) => {
            loadingComponent.close();
            common.DisplayFailureToast(
              this.$buefy,
              err.response.data.metadata.response
            );
          });
      },
    },
  };
</script>

<style scoped></style>
