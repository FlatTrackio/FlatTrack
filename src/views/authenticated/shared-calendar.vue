<template>
  <div>
    <headerDisplay/>
    <div class="container">
      <section class="section">
        <nav class="breadcrumb is-medium has-arrow-separator" aria-label="breadcrumbs">
            <ul>
            <li><a href="/">Home</a></li>
            <li class="is-active"><a href="/shared-calendar">Shared Calendar</a></li>
            </ul>
        </nav>
        <h1 class="title">Shared Calendar</h1>
        <p>Page under construction, coming soon</p>
        <p>Select your platform to add it to your calendar</p>
      </section>
      <section class="section">
        <div class="card-margin" v-for="link of calendarLinks" v-bind:key="link">
          <div class="card">
            <div class="card-content">
              <div class="media">
                <div class="media-content">
                  <p class="title is-4">{{ link.text }}</p>
                </div>
              </div>
              <div class="content">
              </div>
            </div>
            <footer class="card-footer">
              <a :href="`/shared-calendar/${link.address}`" class="card-footer-item">Add</a>
            </footer>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script>
import headerDisplay from '@/components/header-display'

export default {
  name: 'Shopping List',
  data () {
    return {
      calendarLinks: [
        {
          text: 'Android',
          address: ''
        },
        {
          text: 'iOS',
          address: ''
        },
        {
          text: 'WebDav (most things)',
          address: ''
        },
        {
          text: 'CalDav (Macs)',
          address: ''
        }
      ]
    }
  },
  components: {
    headerDisplay
  }
}
</script>

<style scoped>

</style>
