<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><router-link to="/apps/shopping-list">Shopping list</router-link></li>
              <li class="is-active"><router-link to="/apps/shopping-list/tags">Shopping tags</router-link></li>
            </ul>
        </nav>
        <h1 class="title is-1">Shopping tags</h1>
        <p class="subtitle is-3">Manage tags used in your lists</p>
        <div>
          <label class="label">Search for tags</label>
          <b-field>
            <b-input
              icon="magnify"
              size="is-medium"
              placeholder="Enter a tag name"
              type="search"
              expanded
              v-model="tagSearch"
              ref="search">
            </b-input>
            <p class="control">
              <b-select
                placeholder="Sort by"
                icon="sort"
                v-model="sortBy"
                size="is-medium"
                expanded>
                <option value="recentlyAdded">Recently Added</option>
                <option value="lastAdded">Last Added</option>
                <option value="recentlyUpdated">Recently Updated</option>
                <option value="lastUpdated">Last Updated</option>
                <option value="alphabeticalDescending">A-z</option>
                <option value="alphabeticalAscending">z-A</option>
            </b-select>
          </p>
          </b-field>
          <b-loading :is-full-page="false" :active.sync="pageLoading" :can-cancel="false"></b-loading>
          <section>
            <div class="card pointer-cursor-on-hover" @click="AddNewTag">
              <div class="card-content">
                <div class="media">
                  <div class="media-left">
                    <b-icon
                      icon="tag"
                      size="is-medium">
                    </b-icon>
                  </div>
                  <div class="media-content">
                    <p class="title is-4">Add a new tag</p>
                  </div>
                  <div class="media-right">
                    <b-icon icon="chevron-right" size="is-medium" type="is-midgray"></b-icon>
                  </div>
                </div>
              </div>
            </div>
          </section>
        </div>
        <!-- TODO fix floating button -->
        <floatingAddButton :func="AddNewTag" v-if="displayFloatingAddButton" />
        <br/>
        <div v-if="tagsFiltered.length > 0">
          <!-- Card per-tag -->
          <div v-for="(tag, index) in tagsFiltered" v-bind:key="tag">
            <tagCard :tag="tag" :index="index" :tags="tags" />
          </div>
          <br/>
          <p>{{ tagsFiltered.length }} tag(s)</p>
        </div>
        <div v-else>
          <div class="card">
            <div class="card-content card-content-list">
              <div class="media">
                <div class="media-left">
                  <b-icon icon="tag" size="is-medium" type="is-midgray"></b-icon>
                </div>
                <div class="media-content">
                  <p class="subtitle is-4" v-if="tagSearch === '' && tags.length === 0">No tags added yet.</p>
                  <p class="subtitle is-4" v-else-if="tagSearch !== ''">No tags found.</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import common from '@/frontend/common/common'
import shoppinglist from '@/frontend/requests/authenticated/shoppinglist'
import { DialogProgrammatic as Dialog } from 'buefy'

export default {
  name: 'Shopping Tags',
  data () {
    return {
      displayFloatingAddButton: true,
      tags: [],
      deviceIsMobile: false,
      tagSearch: '',
      pageLoading: true,
      sortBy: 'alphabeticalDescending'
    }
  },
  components: {
    floatingAddButton: () => import('@/frontend/components/common/floating-add-button.vue'),
    tagCard: () => import('@/frontend/components/authenticated/shopping-list-tag-card.vue')
  },
  computed: {
    tagsFiltered () {
      return this.tags.filter((item) => {
        return this.TagDisplayState(item)
      })
    },
    newTag () {
      return this.$route.query.newtag
    }
  },
  methods: {
    goToRef (ref) {
      this.$router.push({ path: ref })
    },
    GetShoppingTags () {
      shoppinglist.GetShoppingTags(this.sortBy).then(resp => {
        this.tags = resp.data.list || []
        this.pageLoading = false
      }).catch(err => {
        common.DisplayFailureToast(`Hmmm seems somethings gone wrong loading the shopping tags; ${err.response.data.metadata.response}`)
      })
    },
    AddNewTag () {
      this.displayFloatingAddButton = false
      Dialog.prompt({
        title: 'New tag',
        message: `Enter the name of a tag to create.`,
        container: null,
        icon: 'tag',
        hasIcon: true,
        inputAttrs: {
          placeholder: 'e.g. Fruits and Veges',
          maxlength: 30
        },
        trapFocus: true,
        onConfirm: (value) => {
          shoppinglist.PostShoppingTag(value).then(() => {
            this.pageLoading = true
            this.GetShoppingTags()
            this.displayFloatingAddButton = true
          }).catch(err => {
            common.DisplayFailureToast(`Failed to create tag; ${err.response.data.metadata.response}`)
            this.displayFloatingAddButton = true
          })
        },
        onCancel: () => {
          this.displayFloatingAddButton = true
        }
      })
    },
    TagDisplayState (tag) {
      var vm = this
      return this.SearchTags(tag)
    },
    SearchTags (tag) {
      var vm = this
      return tag.name.toLowerCase().indexOf(vm.tagSearch.toLowerCase()) !== -1
    },
    CheckDeviceIsMobile () {
      this.deviceIsMobile = common.DeviceIsMobile()
    }
  },
  watch: {
    sortBy () {
      this.listIsLoading = true
      this.GetShoppingTags()
    },
    newTag () {
      if (this.newTag === 'prompt') {
        this.AddNewTag()
      }
    }
  },
  async beforeMount () {
    this.GetShoppingTags()
  },
  async created () {
    this.CheckDeviceIsMobile()
    window.addEventListener('resize', this.CheckDeviceIsMobile.bind(this))
  }
}
</script>

<style scoped>

</style>
