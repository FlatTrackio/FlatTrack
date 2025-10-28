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
        <h1 class="title is-1">Profile</h1>
        <p class="subtitle is-3">Manage your account</p>
        <b-loading
          v-model:active="pageLoading"
          :is-full-page="false"
          :can-cancel="false"
        />
        <div class="card">
          <div class="card-content">
            <div class="media">
              <div class="media-content">
                <figure class="image is-128x128">
                  <img src="@/assets/256x256.png" alt="Placeholder image" />
                </figure>
                <p v-if="!pageLoading" class="title is-3">{{ names }}</p>
                <b-skeleton
                  class="mb-5"
                  v-else
                  size="is-medium"
                  width="35%"
                  :animated="true"
                />
                <p v-if="!pageLoading" class="subtitle is-5">
                  Joined {{ TimestampToCalendar(creationTimestamp) }}
                </p>
                <b-skeleton
                  class="mb-5"
                  v-else
                  size="is-medium"
                  width="35%"
                  :animated="true"
                />
              </div>
            </div>
          </div>
        </div>

        <b-field grouped group-multiline>
          <div v-for="group in groups" :key="group" class="control">
            <b-taglist attached>
              <b-tag type="is-dark"> is </b-tag>
              <b-tag type="is-info"> {{ group }} </b-tag>
            </b-taglist>
          </div>
        </b-field>

        <b-field label="Name(s)">
          <b-input
            v-model="names"
            type="text"
            maxlength="60"
            placeholder="Enter your name(s)"
            icon="form-textbox"
            size="is-medium"
            icon-right="close-circle"
            icon-right-clickable
            required
            @icon-right-click="names = ''"
            @keyup.enter.native="PatchProfile"
          />
        </b-field>

        <b-field label="Email">
          <b-input
            v-model="email"
            type="email"
            maxlength="70"
            icon="email"
            size="is-medium"
            placeholder="Enter your email address"
            icon-right="close-circle"
            icon-right-clickable
            required
            @icon-right-click="email = ''"
            @keyup.enter.native="PatchProfile"
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
            @keyup.enter.native="PatchProfile"
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
        <b-button
          type="is-success"
          size="is-medium"
          icon-left="delta"
          native-type="submit"
          expanded
          @click="PatchProfile"
        >
          Update profile
        </b-button>
      </section>
    </div>
  </div>
</template>

<script>
  import common from "@/common/common";
  import profile from "@/requests/authenticated/profile";
  import breadcrumb from "@/components/common/breadcrumb.vue";

  export default {
    name: "AccountProfile",
    components: {
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
        birthday: null,
        jsBirthday: null,
        pageLoading: true,
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
        return new Date(this.jsBirthday || 0).getTime() / 1000 || null;
      },
    },
    async beforeMount() {
      this.GetProfile();
    },
    methods: {
      CopyHrefToClipboard() {
        common.CopyHrefToClipboard();
      },
      GetProfile() {
        profile.GetProfile().then((resp) => {
          this.names = resp.data.spec.names;
          this.birthday = resp.data.spec.birthday;
          this.jsBirthday =
            this.birthday !== null ? new Date(this.birthday * 1000) : null;
          this.focusedDate = this.jsBirthday;
          this.phoneNumber = resp.data.spec.phoneNumber;
          this.groups = resp.data.spec.groups;
          this.email = resp.data.spec.email;
          this.creationTimestamp = resp.data.spec.creationTimestamp;
          this.pageLoading = false;
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
      TimestampToCalendar(timestamp) {
        return common.TimestampToCalendar(timestamp);
      },
    },
  };
</script>
