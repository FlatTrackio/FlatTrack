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
        <breadcrumb
          back-link-name="Admin accounts"
          :current-page-name="$route.name"
        />
        <h1 class="title is-1">Edit account</h1>
        <p class="subtitle is-4">Edit an existing user account</p>
        <b-loading
          v-model:active="pageLoading"
          :is-full-page="false"
          :can-cancel="false"
        />
        <div v-if="registered !== true">
          <div class="notification is-warning">
            <p class="subtitle is-6">
              <strong
                >This account has been created but doesn't appear to be
                registered.</strong
              >
            </p>
            <b-button
              :icon-left="
                showRegistrationCompletionDetails === false ? 'eye' : 'eye-off'
              "
              @click="
                showRegistrationCompletionDetails =
                  !showRegistrationCompletionDetails
              "
            >
              {{ showRegistrationCompletionDetails === false ? "Show" : "Hide"
              }} registration details
            </b-button>
            <div v-if="showRegistrationCompletionDetails === true">
              <br />
              <div class="notification">
                <div class="content">
                  <qrcode-vue
                    :value="
                      windowOrigin +
                        '/useraccountconfirm/' +
                        userAccountConfirmId +
                        '?secret=' +
                        userAccountConfirmSecret
                    "
                    :size="200"
                    level="H"
                  />
                  <br />
                  <p>
                    Have your flatmate scan the QR code above, or
                    <a type="is-text" @click="CopyRegistrationLink"
                      >click here</a
                    >
                    to copy the registration link for you to send to your
                    flatmate
                  </p>
                </div>
              </div>
            </div>
          </div>
          <br />
        </div>

        <b-field label="Name(s)">
          <b-input
            v-model="names"
            type="text"
            maxlength="60"
            placeholder="Enter your flatmate's name"
            icon="form-textbox"
            size="is-medium"
            icon-right="close-circle"
            icon-right-clickable
            required
            @icon-right-click="names = ''"
            @keyup.enter.native="PatchUserAccount"
          />
        </b-field>

        <b-field label="Email">
          <b-input
            v-model="email"
            type="email"
            maxlength="70"
            placeholder="Enter your flatmate's email"
            icon="email"
            size="is-medium"
            icon-right="close-circle"
            icon-right-clickable
            required
            @icon-right-click="email = ''"
            @keyup.enter.native="PatchUserAccount"
          />
        </b-field>

        <section>
          <b-field label="Groups">
            <b-taginput
              v-model="groupsFull"
              :data="availableGroups"
              field="name"
              autocomplete
              open-on-focus
              ellipsis
              icon="account-group"
              placeholder="Select groups"
              size="is-medium"
              @keyup.enter.native="PatchUserAccount"
              @typing="GetFilteredGroups"
            />
          </b-field>
        </section>
        <br />

        <b-field label="Phone number (optional)">
          <b-input
            v-model="phoneNumber"
            type="tel"
            placeholder="Enter your flatmate's phone number"
            icon="phone"
            size="is-medium"
            icon-right="close-circle"
            icon-right-clickable
            maxlength="30"
            @icon-right-click="phoneNumber = ''"
            @keyup.enter.native="PatchUserAccount"
          />
        </b-field>

        <b-field label="Birthday (optional)">
          <b-datepicker
            v-model="jsBirthday"
            :max-date="maxDate"
            :show-week-numbers="true"
            :focused-date="focusedDate"
            placeholder="Click to select birthday"
            icon="cake-variant"
            size="is-medium"
            trap-focus
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
            placeholder="Enter a password for your flatmate"
            icon="form-textbox-password"
            pattern="^([a-zA-Z]*).{10,}$"
            validation-message="password is invalid. Make sure that your password has: 10 or more characters, at least one lower case letter, at least one upper case letter, at least one number"
            icon-right="close-circle"
            icon-right-clickable
            size="is-medium"
            @icon-right-click="password = ''"
            @keyup.enter.native="PatchUserAccount"
          />
        </b-field>

        <b-field label="Confirm password">
          <b-input
            v-model="passwordConfirm"
            type="password"
            password-reveal
            maxlength="70"
            placeholder="Confirm a password for your flatmate"
            icon="form-textbox-password"
            pattern="^([a-zA-Z]*).{10,}$"
            validation-message="password is invalid. Make sure that your password has: 10 or more characters, at least one lower case letter, at least one upper case letter, at least one number"
            icon-right="close-circle"
            icon-right-clickable
            size="is-medium"
            @icon-right-click="passwordConfirm = ''"
            @keyup.enter.native="PatchUserAccount"
          />
        </b-field>
        <b-field>
          <b-button
            type="is-success"
            size="is-medium"
            icon-left="delta"
            native-type="submit"
            expanded
            @click="PatchUserAccount"
          >
            Update user account
          </b-button>
          <p v-if="myUserID !== id" class="control">
            <b-button
              type="is-danger"
              size="is-medium"
              icon-left="delete"
              native-type="submit"
              @click="DeactivateUserAccount(id)"
            />
          </p>
        </b-field>

        <b-collapse
          v-if="myUserID !== id"
          v-model="accountAdvancedOpen"
          class="card"
          animation="slide"
          aria-id="contentIdForA11y3"
        >
          <template #trigger="props">
            <div
              class="card-header"
              role="button"
              aria-controls="contentIdForA11y3"
            >
              <p class="card-header-title">Advanced options</p>
              <a class="card-header-icon">
                <b-icon :icon="props.open ? 'menu-down' : 'menu-up'" />
              </a>
            </div>
          </template>
          <div class="card-content">
            <div class="content">
              <b-button
                size="is-medium"
                :type="disabled ? 'is-danger' : 'is-warning'"
                :loading="disabledLoading"
                :icon-left="
                  disabled ? 'checkbox-outline' : 'checkbox-blank-outline'
                "
                @click="PatchUserAccountDisabled"
              >
                Disable
              </b-button>
            </div>
          </div>
        </b-collapse>
        <p class="subtitle is-6">
          Created {{ TimestampToCalendar(creationTimestamp) }} <br />
          <span v-if="creationTimestamp !== modificationTimestamp">
            Modified {{ TimestampToCalendar(modificationTimestamp) }}
          </span>
        </p>
      </section>
    </div>
  </div>
</template>

<script>
  import common from "@/common/common";
  import groups from "@/requests/authenticated/groups";
  import flatmates from "@/requests/authenticated/flatmates";
  import adminFlatmates from "@/requests/admin/flatmates";
  import QrcodeVue from "qrcode.vue";
  import infotooltip from "@/components/common/info-tooltip.vue";
  import breadcrumb from "@/components/common/breadcrumb.vue";

  export default {
    name: "EditAccount",
    components: {
      QrcodeVue,
      infotooltip,
      breadcrumb,
    },
    data() {
      var today = new Date();
      var maxDate = new Date(
        today.getFullYear() - 15,
        today.getMonth(),
        today.getYear()
      );
      var minDate = new Date(
        today.getFullYear() - 100,
        today.getMonth(),
        today.getYear()
      );
      var windowOrigin = window.location.origin;

      return {
        windowOrigin: windowOrigin,
        maxDate: maxDate,
        minDate: minDate,
        showRegistrationCompletionDetails: false,
        userAccountConfirmId: null,
        userAccountConfirmSecret: null,
        pageLoading: true,
        disabledLoading: false,
        accountAdvancedOpen: false,
        myUserID: null,
        id: this.$route.params.id,
        names: null,
        email: null,
        phoneNumber: null,
        birthday: null,
        groups: [],
        registered: null,
        disabled: false,
        password: null,
        passwordConfirm: null,
        availableGroups: [],
        jsBirthday: null,
        groupsFull: [],
      };
    },
    async beforeMount() {
      this.GetAvailableGroups();
      this.GetUserAccount();
      this.myUserID = common.GetUserIDFromJWT();
      if (this.registered !== true) {
        adminFlatmates.GetUserAccountConfirms(this.id).then((resp) => {
          var confirmsList = resp.data.list;
          for (var confirmItem in confirmsList) {
            if (confirmsList[confirmItem].userId !== this.id) {
              continue;
            }
            this.userAccountConfirmId = confirmsList[confirmItem].id;
            this.userAccountConfirmSecret = confirmsList[confirmItem].secret;
          }
        });
      }
    },
    methods: {
      CopyHrefToClipboard() {
        common.CopyHrefToClipboard();
      },
      TimestampToCalendar(timestamp) {
        return common.TimestampToCalendar(timestamp);
      },
      GetAvailableGroups() {
        groups
          .GetGroups()
          .then((resp) => {
            this.availableGroups = resp.data.list;
            resp.data.list.map((group) => {
              if (group.defaultGroup === true) {
                this.groupsFull.push(group);
              }
            });
          })
          .catch((err) => {
            common.DisplayFailureToast(
              this.$buefy,
              "Failed to list groups" +
                "<br/>" +
                err.response.data.metadata.response
            );
          });
      },
      GetFilteredGroups(text) {
        this.groups = this.availableGroups.filter((group) => {
          return (
            group.name.toString().toLowerCase().indexOf(text.toLowerCase()) >= 0
          );
        });
      },
      GetUserAccount() {
        flatmates
          .GetFlatmate(this.id)
          .then((resp) => {
            var user = resp.data.spec;
            this.names = user.names;
            this.email = user.email;
            this.phoneNumber = user.phoneNumber;
            this.birthday = user.birthday;
            this.registered = user.registered;
            this.disabled = user.disabled;
            this.groups = user.groups;
            this.groupsFull = [];
            this.creationTimestamp = user.creationTimestamp;
            this.modificationTimestamp = user.modificationTimestamp;
            this.availableGroups.map((group) => {
              if (this.groups.includes(group.name)) {
                this.groupsFull.push(group);
              }
            });
            this.jsBirthday =
              typeof this.birthday !== "undefined"
                ? new Date(this.birthday * 1000)
                : null;
            this.pageLoading = false;
          })
          .catch((err) => {
            common.DisplayFailureToast(
              this.$buefy,
              this.$buefy,
              "Failed to fetch user account" +
                "<br/>" +
                (err.response.data.metadata.response || err)
            );
            this.$router.push({ name: "Admin accounts" });
          });
      },
      PatchUserAccount() {
        if (
          this.password !== this.passwordConfirm &&
          this.password !== null &&
          typeof this.password !== "undefined"
        ) {
          common.DisplayFailureToast(this.$buefy, "Passwords do not match");
          return;
        }
        this.birthday = new Date(this.jsBirthday || 0).getTime() / 1000 || null;

        var groups = [];
        this.groupsFull.map((group) => {
          if (group === "" || group.name === "") {
            return;
          }
          groups.push(group.name);
        });
        adminFlatmates
          .PatchFlatmate(
            this.id,
            this.names,
            this.email,
            this.phoneNumber,
            this.birthday,
            groups,
            this.password
          )
          .then((resp) => {
            common.DisplaySuccessToast(this.$buefy, "Updated user account");
            setTimeout(() => {
              this.$router.push({ name: "Admin accounts" });
            }, 1 * 1000);
          })
          .catch((err) => {
            common.DisplayFailureToast(
              this.$buefy,
              "Failed to patch user account" +
                "<br/>" +
                (err.response.data.metadata.response || err)
            );
          });
      },
      PatchUserAccountDisabled() {
        this.disabledLoading = true;
        adminFlatmates
          .PatchFlatmateDisabled(this.id, !this.disabled)
          .then((resp) => {
            this.disabled = resp.data.spec.disabled;
            this.disabledLoading = false;
            common.DisplaySuccessToast(
              this.$buefy,
              `${this.disabled ? "Disabled" : "Enabled"} user account`
            );
          })
          .catch((err) => {
            this.disabledLoading = false;
            common.DisplayFailureToast(
              this.$buefy,
              "Failed to patch user account disabled field" +
                "<br/>" +
                (err.response.data.metadata.response || err)
            );
          });
      },
      DeactivateUserAccount(id) {
        this.$buefy.dialog.confirm({
          title: "Deactivate user account",
          message:
            "Are you sure that you wish to permanently deactivate this account?" +
            "<br/>" +
            "This action cannot be undone.",
          confirmText: "Deactivate account",
          type: "is-danger",
          hasIcon: true,
          onConfirm: () => {
            adminFlatmates
              .DeleteFlatmate(id)
              .then((resp) => {
                common.DisplaySuccessToast(
                  this.$buefy,
                  "Permanently deactivated user account"
                );
                setTimeout(() => {
                  this.$router.push({ name: "Admin accounts" });
                }, 1 * 1000);
              })
              .catch((err) => {
                common.DisplayFailureToast(
                  this.$buefy,
                  "Failed to permanently deactive user account" +
                    "<br/>" +
                    (err.response.data.metadata.response || err)
                );
              });
          },
        });
      },
      CopyRegistrationLink() {
        var registrationLink = `${window.location.origin}/useraccountconfirm/${this.userAccountConfirmId}?secret=${this.userAccountConfirmSecret}`;
        window.prompt("Copy the following link", registrationLink);
      },
    },
  };
</script>
