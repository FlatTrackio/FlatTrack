<template>
  <div>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
              <li><a href="/">Home</a></li>
              <li class="is-active"><a href="/flat">My flat</a></li>
            </ul>
        </nav>
        <h1 class="title">{{ flatName }}</h1>
        <p>About your flat</p>
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
