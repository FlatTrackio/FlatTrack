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
        <breadcrumb back-link-name="Account" :current-page-name="$route.name" />
        <h1 class="title is-1">Security</h1>
        <p class="subtitle is-3">Manage your account's security</p>
        <div class="mb-5">
          <div class="field has-addons">
            <h1 class="title is-3">Password</h1>
            <p class="control">
              <infotooltip
                message="Make sure that your password has: 10 or more characters, at least one lower case letter, at least one upper case letter, at least one number"
              />
            </p>
          </div>
          <b-field label="Password">
            <b-input
              v-model="password"
              type="password"
              password-reveal
              placeholder="Enter to update your password"
              icon="form-textbox-password"
              size="is-medium"
              pattern="^([a-zA-Z]*).{10,}$"
              validation-message="password is invalid. Make sure that your password has: 10 or more characters, at least one lower case letter, at least one upper case letter, at least one number"
              maxlength="70"
              @keyup.enter.native="PatchProfile"
            />
          </b-field>
          <b-field label="Confirm password">
            <b-input
              v-model="passwordConfirm"
              type="password"
              password-reveal
              placeholder="Confirm to update your password"
              maxlength="70"
              size="is-medium"
              pattern="^([a-zA-Z]*).{10,}$"
              validation-message="password is invalid. Make sure that your password has: 10 or more characters, at least one lower case letter, at least one upper case letter, at least one number"
              icon="form-textbox-password"
              @keyup.enter.native="PatchProfile"
            />
          </b-field>
          <b-button
            type="is-success"
            size="is-medium"
            icon-left="delta"
            native-type="submit"
            expanded
            @click="PatchProfile"
          >
            Update password
          </b-button>
        </div>

        <h1 class="title is-3">Two-factor authentication</h1>
        <div class="field">
          <b-checkbox
            v-model="otpEnable"
            :type="otpIsEnabled === true ? 'is-success' : 'is-warning'"
            size="is-medium"
            :indeterminate="otpEnable"
            disabled
          >
            OTP
          </b-checkbox>
        </div>
        <div v-if="otpEnable">
          <qrcode-vue :value="''" :size="200" level="H" />
          <b-field label="Confirm your OTP code">
            <b-input
              type="number"
              placeholder="Enter your OTP code from your phone"
              maxlength="6"
              size="is-medium"
              icon="qrcode"
            />
          </b-field>
          <b-button
            type="is-success"
            size="is-medium"
            rounded
            expanded
            native-type="submit"
          >
            Enable
          </b-button>
        </div>

        <h1 class="title is-3">Sign out of all devices</h1>
        <div class="notification is-warning mb-4">
          <p class="subtitle is-6">
            <b>Please note:</b> revoking access is not unable and will require
            signing in again (including from this device).
          </p>
        </div>
        <b-button
          type="is-danger"
          size="is-medium"
          icon-left="close"
          expanded
          @click="ResetAuth"
        >
          Revoke access for all devices
        </b-button>
      </section>
    </div>
  </div>
</template>

<script>
  import common from "@/common/common";
  import profile from "@/requests/authenticated/profile";
  import infotooltip from "@/components/common/info-tooltip.vue";
  import breadcrumb from "@/components/common/breadcrumb.vue";

  export default {
    name: "AccountSecurity",
    components: {
      infotooltip,
      breadcrumb,
    },
    data() {
      const today = new Date();
      const maxDate = new Date(
        today.getFullYear() - 15,
        today.getMonth(),
        today.getDate()
      );
      const minDate = new Date(
        today.getFullYear() - 100,
        today.getMonth(),
        today.getDate()
      );

      return {
        maxDate: maxDate,
        minDate: minDate,
        focusedDate: maxDate,
        passwordConfirm: "",
        otpEnable: false,
        otpIsEnabled: false,
        jsBirthday: null,
        names: "",
        email: "",
        phoneNumber: "",
        groups: [],
        password: "",
        creationTimestamp: "",
      };
    },
    computed: {
      birthday() {
        return new Date(this.jsBirthday || 0).getTime() / 1000 || 0;
      },
    },
    methods: {
      CopyHrefToClipboard() {
        common.CopyHrefToClipboard();
      },
      GetProfile() {
        profile.GetProfile().then((resp) => {
          this.names = resp.data.spec.names;
          this.birthday = resp.data.spec.birthday || null;
          this.jsBirthday =
            this.birthday !== null ? new Date(this.birthday * 1000) : null;
          this.focusedDate = this.jsBirthday;
          this.phoneNumber = resp.data.spec.phoneNumber;
          this.groups = resp.data.spec.groups;
          this.email = resp.data.spec.email;
          this.creationTimestamp = resp.data.spec.creationTimestamp;
        });
      },
      PatchProfile() {
        if (this.password !== this.passwordConfirm) {
          common.DisplayFailureToast(
            this.$buefy,
            "Unable to use password as they either do not match"
          );
          return;
        }
        profile
          .PatchProfile(
            this.names,
            this.email,
            this.phoneNumber,
            this.birthday,
            this.password
          )
          .then((resp) => {
            if (resp.data.spec.id === "") {
              common.DisplayFailureToast(
                this.$buefy,
                "Failed to update profile"
              );
              return;
            }
            common.DisplaySuccessToast(
              this.$buefy,
              "Successfully updated your profile"
            );
            this.password = "";
            this.passwordConfirm = "";
          })
          .catch((err) => {
            common.DisplayFailureToast(
              this.$buefy,
              "Failed to update profile" +
                "<br/>" +
                err.response.data.metadata.response
            );
          });
      },
      ResetAuth() {
        this.$buefy.dialog.confirm({
          title: "Revoke access for all devices",
          message:
            "Are you sure that you wish to sign out of all devices?" +
            "<br/>" +
            "This action cannot be undone.",
          confirmText: "Sign out",
          type: "is-danger",
          hasIcon: true,
          onConfirm: () => {
            profile
              .PostAuthReset()
              .then((resp) => {
                common.DisplaySuccessToast(
                  this.$buefy,
                  "Successfully signed out of all devices"
                );
                window.location.href = "/login";
              })
              .catch((err) => {
                common.DisplayFailureToast(
                  this.$buefy,
                  "Failed to sign out of all devices" +
                    "<br/>" +
                    err.response.data.metadata.response
                );
              });
          },
        });
      },
      TimestampToCalendar(timestamp) {
        return common.TimestampToCalendar(timestamp);
      },
    },
  };
</script>
