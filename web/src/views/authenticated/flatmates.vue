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
        <breadcrumb back-link-name="Apps" :current-page-name="$route.name" />
        <h1 class="title is-1">Flatmates</h1>
        <p class="subtitle is-3">
          <span v-if="typeof GroupQuery === 'undefined'">
            Get to know your flatmates
          </span>
          <span v-else>
            Listing flatmates, filtering by the group {{ GroupQuery }}
          </span>
        </p>
        <b-loading
          v-model:active="pageLoading"
          :is-full-page="false"
          :can-cancel="false"
        />
        <div v-if="members && members.length > 0">
          <div v-for="member of members" :key="member" class="card-margin">
            <div class="card">
              <div class="card-content">
                <div class="media">
                  <div class="media-left">
                    <figure class="image is-48x48">
                      <img src="@/assets/96x96.png" alt="Placeholder image" />
                    </figure>
                  </div>
                  <div class="media-content">
                    <p class="title is-4">{{ member.names }}</p>
                    <p class="subtitle is-6">
                      Joined {{ TimestampToCalendar(member.creationTimestamp) }}
                    </p>
                  </div>
                </div>
                <div class="content">
                  <b-field grouped group-multiline>
                    <div
                      v-for="group in member.groups"
                      :key="group"
                      class="control"
                    >
                      <b-taglist attached>
                        <b-tag type="is-dark"> is </b-tag>
                        <b-tag type="is-info"> {{ group }} </b-tag>
                      </b-taglist>
                    </div>
                  </b-field>
                  <p v-if="member.phoneNumber" class="subtitle is-6">
                    Phone:
                    <a :href="`tel:${member.phoneNumber}`"
                      >{{ member.phoneNumber }}</a
                    ><br />
                  </p>
                  <p v-if="member.email" class="subtitle is-6">
                    Email:
                    <a :href="`mailto:${member.email}`">{{ member.email }}</a
                    ><br />
                  </p>
                  <a
                    v-if="member.birthday && member.birthday !== 0"
                    class="subtitle is-6"
                  >
                    Birthday: {{ TimestampToCalendar(member.birthday) }}<br />
                  </a>
                  <b-field
                    v-if="
                      member.registered !== true ||
                        member.disabled === true ||
                        member.deletionTimestamp !== 0
                    "
                    grouped
                    group-multiline
                  >
                    <div class="control">
                      <b-taglist v-if="member.registered !== true" attached>
                        <b-tag type="is-dark"> has </b-tag>
                        <b-tag type="is-danger"> not registered </b-tag>
                      </b-taglist>
                    </div>
                    <div class="control">
                      <b-taglist v-if="member.disabled === true" attached>
                        <b-tag type="is-dark"> has </b-tag>
                        <b-tag type="is-warning"> account disabled </b-tag>
                      </b-taglist>
                    </div>
                    <div class="control">
                      <b-taglist v-if="member.deletionTimestamp !== 0" attached>
                        <b-tag type="is-dark"> has </b-tag>
                        <b-tag type="is-danger"> account deactivatived </b-tag>
                      </b-taglist>
                    </div>
                  </b-field>
                </div>
              </div>
            </div>
          </div>
          <b-button
            v-if="$route.query.id || $route.query.group"
            type="is-warning"
            expanded
            @click="ClearFilter"
          >
            Clear filter
          </b-button>
          <div class="m-3">
            <p>
              {{ members.length }} {{ members.length === 1 ? "flatmate" :
              "flatmates" }}
            </p>
          </div>
        </div>
        <div v-if="!members || members.length === 0">
          <div class="card">
            <div class="card-content">
              <div class="media">
                <div class="media-left">
                  <b-icon icon="account-off" size="is-medium" />
                </div>
                <div class="media-content">
                  <p v-if="!pageLoading" class="subtitle is-4">
                    No flatmates found.
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
              <p class="content subtitle is-6">
                Either you haven't added any flatmates, or you are trying to
                search or filter for flatmates and none could be found.
              </p>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
  import * as emoji from "node-emoji";
  import flatmates from "@/requests/authenticated/flatmates";
  import common from "@/common/common";
  import breadcrumb from "@/components/common/breadcrumb.vue";

  export default {
    name: "FlatFlatmates",
    components: {
      breadcrumb,
    },
    data() {
      return {
        pageLoading: true,
        members: [],
        emojiSmile: emoji.get("smile"),
      };
    },
    computed: {
      GroupQuery() {
        return this.$route.query.group;
      },
      IdQuery() {
        return this.$route.query.id;
      },
    },
    watch: {
      GroupQuery() {
        this.FetchAllFlatmates();
      },
    },
    async beforeMount() {
      if (typeof this.IdQuery !== "undefined") {
        var id = this.IdQuery;
        flatmates
          .GetFlatmate(id)
          .then((resp) => {
            this.members = [resp.data.spec];
            this.pageLoading = false;
          })
          .catch((err) => {
            common.DisplayFailureToast(
              "Failed fetch flatmate info" + `<br/>${err}`
            );
          });
      } else {
        this.FetchAllFlatmates();
      }
    },
    methods: {
      CopyHrefToClipboard() {
        common.CopyHrefToClipboard();
      },
      FetchAllFlatmates() {
        if (typeof this.GroupQuery !== "undefined") {
          var group = this.GroupQuery;
        } else if (typeof this.IdQuery !== "undefined") {
          var id = this.IdQuery;
        }
        var notSelf = true;
        flatmates
          .GetAllFlatmates(id, notSelf, group)
          .then((resp) => {
            this.pageLoading = false;
            this.members = resp.data.list;
          })
          .catch((err) => {
            common.DisplayFailureToast(
              "Failed to list flatmates" + `<br/>${err}`
            );
          });
      },
      ClearFilter() {
        this.$router.replace({ name: "My Flatmates" });
        this.FetchAllFlatmates();
      },
      TimestampToCalendar(timestamp) {
        return common.TimestampToCalendar(timestamp);
      },
    },
  };
</script>
