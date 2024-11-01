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
        <nav
          class="breadcrumb is-medium has-arrow-separator"
          aria-label="breadcrumbs"
        >
          <ul>
            <li>
              <router-link :to="{ name: 'Admin home' }">
                Admin
              </router-link>
            </li>
            <li class="is-active">
              <router-link :to="{ name: 'Admin accounts' }">
                Accounts
              </router-link>
            </li>
          </ul>
          <b-button
            icon-left="content-copy"
            size="is-small"
            @click="CopyHrefToClipboard()"
          />
        </nav>
        <h1 class="title is-1">
          Accounts
        </h1>
        <p class="subtitle is-4">
          Manage the account of your flatmates
        </p>
        <b-loading
          v-model:active="pageLoading"
          :is-full-page="false"
          :can-cancel="false"
        />
        <div>
          <section>
            <div
              class="card pointer-cursor-on-hover"
              @click="$router.push({ name: 'Admin new account' })"
            >
              <div class="card-content">
                <div class="media">
                  <div class="media-left">
                    <b-icon
                      icon="account-plus"
                      size="is-medium"
                    />
                  </div>
                  <div class="media-content">
                    <p class="title is-4">
                      Add a new flatmate
                    </p>
                  </div>
                  <div class="media-right">
                    <b-icon
                      icon="chevron-right"
                      size="is-medium"
                      type="is-midgray"
                    />
                  </div>
                </div>
              </div>
            </div>
          </section>
        </div>
        <br>
        <div v-if="members && members.length">
          <div
            v-for="member of members"
            :key="member"
            class="card-margin pointer-cursor-on-hover"
            @click="
              $router.push({
                name: 'View user account',
                params: { id: member.id },
              })
            "
          >
            <div class="card">
              <div class="card-content">
                <div class="media">
                  <div class="media-left">
                    <figure class="image is-48x48">
                      <img
                        src="@/assets/96x96.png"
                        alt="Placeholder image"
                      >
                    </figure>
                  </div>
                  <div class="media-content">
                    <p class="title is-4">
                      {{ member.names }}
                    </p>
                    <p class="subtitle is-6">
                      Joined {{ TimestampToCalendar(member.creationTimestamp) }}
                    </p>
                  </div>
                  <div class="media-right">
                    <b-icon
                      icon="chevron-right"
                      size="is-medium"
                      type="is-midgray"
                    />
                  </div>
                </div>
                <div class="content">
                  <b-field
                    grouped
                    group-multiline
                  >
                    <div
                      v-for="group in member.groups"
                      :key="group"
                      class="control"
                    >
                      <b-taglist attached>
                        <b-tag type="is-dark">
                          is
                        </b-tag>
                        <b-tag type="is-info">
                          {{ group }}
                        </b-tag>
                      </b-taglist>
                    </div>
                  </b-field>
                  <p
                    v-if="member.phoneNumber"
                    class="subtitle is-6"
                  >
                    Phone:
                    <a :href="`tel:${member.phoneNumber}`">{{
                      member.phoneNumber
                    }}</a><br>
                  </p>
                  <p
                    v-if="member.email"
                    class="subtitle is-6"
                  >
                    Email:
                    <a :href="`mailto:${member.email}`">{{ member.email }}</a><br>
                  </p>
                  <a
                    v-if="member.birthday && member.birthday !== 0"
                    class="subtitle is-6"
                  >
                    Birthday: {{ TimestampToCalendar(member.birthday) }}<br>
                  </a>
                  <b-field
                    v-if="
                      member.registered !== true || member.disabled === true
                    "
                    grouped
                    group-multiline
                  >
                    <div class="control">
                      <b-taglist
                        v-if="member.registered !== true"
                        attached
                      >
                        <b-tag type="is-dark">
                          has
                        </b-tag>
                        <b-tag type="is-danger">
                          not registered
                        </b-tag>
                      </b-taglist>
                    </div>
                    <div class="control">
                      <b-taglist
                        v-if="member.disabled === true"
                        attached
                      >
                        <b-tag type="is-dark">
                          has
                        </b-tag>
                        <b-tag type="is-warning">
                          account disabled
                        </b-tag>
                      </b-taglist>
                    </div>
                  </b-field>
                </div>
              </div>
            </div>
          </div>
          <div class="section">
            <p>
              {{ members.length }}
              {{ members.length === 1 ? "flatmate" : "flatmates" }}
            </p>
          </div>
        </div>
        <div v-if="members && !members.length">
          <div class="card">
            <div class="card-content">
              <div class="media">
                <div class="media-left">
                  <b-icon
                    icon="account-off"
                    size="is-medium"
                  />
                </div>
                <div class="media-content">
                  <p
                    v-if="!pageLoading"
                    class="subtitle is-4"
                  >
                    No flatmates found.
                  </p>
                  <p
                    v-else-if="pageLoading"
                    class="subtitle is-4"
                  >
                    Loading flatmates...
                  </p>
                </div>
              </div>
              <p class="content subtitle is-5">
                Hmmm, it appears that you don't have an flatmates added.<br>
              </p>
            </div>
          </div>
        </div>
        <floatingAddButton :router-link="{ name: 'Admin new account' }" />
      </section>
    </div>
  </div>
</template>

<script>
import emoji from 'node-emoji'
import flatmates from '@/requests/authenticated/flatmates'
import common from '@/common/common'

export default {
  name: 'FlatmatesAccounts',
  components: {
    floatingAddButton: () =>
      import('@/components/common/floating-add-button.vue')
  },
  data () {
    return {
      members: [],
      groupQuery: undefined,
      emojiSmile: emoji.get('smile'),
      pageLoading: true
    }
  },
  async beforeMount () {
    this.groupQuery = this.$route.query.group
    this.FetchAllFlatmates()
  },
  methods: {
    CopyHrefToClipboard () {
      common.CopyHrefToClipboard()
    },
    FetchAllFlatmates () {
      if (typeof this.groupQuery !== 'undefined') {
        var group = this.groupQuery
      }
      flatmates
        .GetAllFlatmates(undefined, undefined, group)
        .then((resp) => {
          this.pageLoading = false
          this.members = resp.data.list
        })
        .catch((err) => {
          common.DisplayFailureToast(
            'Failed to list flatmates' +
              `<br/>${err.response.data.metadata.response}`
          )
        })
    },
    TimestampToCalendar (timestamp) {
      return common.TimestampToCalendar(timestamp)
    }
  }
}
</script>

<style src="../../assets/style.css"></style>

<style scoped></style>
