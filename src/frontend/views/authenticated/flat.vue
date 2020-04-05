<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><router-link to="/">Home</router-link></li>
              <li class="is-active"><router-link to="/flat">My flat</router-link></li>
            </ul>
        </nav>
        <h1 class="title is-1">{{ flatName }}</h1>
        <p class="subtitle is-3">About your flat</p>
      </section>
    </div>
  </div>
</template>

<script>
import flatInfo from '@/frontend/requests/authenticated/flatInfo'

export default {
  name: 'flat',
  data () {
    return {
      flatName: 'My flat'
    }
  },
  methods: {
    GetFlatName () {
      flatInfo.GetFlatName().then(resp => {
        this.flatName = resp.data.spec
      })
    }
  },
  async created () {
    this.GetFlatName()
  }
}
</script>
