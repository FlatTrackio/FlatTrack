<template>
  <div>
    <section class="hero is-dark">
      <div class="hero-body">
        <p class="title">
          FlatTrack - {{ deploymentName }}
        </p>
      </div>
    </section>
    <div class="container">
      <div id="menu-bar-items">
        <b-menu>
          <b-menu-list>
            <b-menu-item v-for="item in pages" v-bind:key="item" :href="item.url" class="menu-bar-item" :label="item.name"></b-menu-item>
          </b-menu-list>
        </b-menu>
      </div>
    </div>
    <footer class="footer">
      <div class="content has-text-centered">
        <div class="columns is-variable is-1-mobile is-0-tablet is-3-desktop is-8-widescreen is-2-fullhd">
          <div class="column">
            <a href="/#/aboutflat">About this flat</a>
          </div>
          <div class="column">
            <a href="/#/aboutflattrack">About this FlatTrack</a>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script>
import axios from 'axios'

export default {
  name: 'home',
  data () {
    return {
      deploymentName: '',
      pageErrors: [],
      pages: [
        {
          name: 'Tasks',
          url: '#/tasks',
          to: 'tasks'
        },
        {
          name: 'Shopping List',
          url: '#/shoppinglist',
          to: 'shoppinglist'
        },
        {
          name: 'Noticeboard',
          url: '#/noticeboard',
          to: 'noticeboard'
        },
        {
          name: 'Shared Calendar',
          url: '#/cal',
          to: 'cal'
        },
        {
          name: 'Recipes',
          url: '#/recipes',
          to: 'recipes'
        },
        {
          name: 'Flatmates',
          url: '#/members',
          to: 'members'
        },
        {
          name: 'Highfives',
          url: '#/highfives',
          to: 'highfives'
        }
      ]
    }
  },
  methods () {
  },
  created () {
    axios.get(`${location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : '')}/api/settings`)
      .then(response => {
        this.deploymentName = response.data.deploymentName
      })
      .catch(err => {
        this.pageErrors.push(err)
      })
  }

}
</script>

<style scoped>
#menu-bar-items {
  padding-left: 10px;
  padding-top: 10px;
  color: black;
}

.menu-bar-item {
  border-top: 1px solid rgba(111, 111, 111, 0.4);
}
</style>
