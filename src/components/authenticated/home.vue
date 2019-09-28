<template>
  <div>
    <headerDisplay/>
    <div class="container">
      <section class="section">
        <h2 class="title is-2">Hey, {{ login.name }}!</h2>
        <h4 class="title is-4">Welcome to FlatTrack, where your flat or community house is organized</h4>
        <div id="menu-bar-items">
          <b-menu>
            <b-menu-list>
              <div v-for="item in pages" v-bind:key="item">
                <b-menu-item :href="item.url" class="menu-bar-item" :label="item.name" :disabled="item.disabled" v-if="determineItemDisplay(item, login)"></b-menu-item>
              </div>
            </b-menu-list>
          </b-menu>
          <br>
          <b-button rounded @click="memberSignOut">Sign Out</b-button>
        </div>
      </section>
    </div>
    <footer class="footer">
      <div class="content has-text-centered">
        <div class="columns is-variable is-1-mobile is-0-tablet is-3-desktop is-8-widescreen is-2-fullhd">
          <div class="column">
            <a href="/#/aboutflat">About this flat</a>
          </div>
          <div class="column">
            <a href="/#/aboutflattrack">About FlatTrack</a>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script>
import headerDisplay from '../common/header-display'
import { LoadingProgrammatic as Loading, DialogProgrammatic as Dialog } from 'buefy'

export default {
  name: 'home',
  data () {
    return {
      pageErrors: [],
      pageLocation: location.protocol + '//' + location.hostname + (location.port ? ':' + location.port : ''),
      pages: [
        {
          name: 'Tasks',
          url: '#/tasks',
          to: 'tasks',
          disabled: false
        },
        {
          name: 'Shopping List',
          url: '#/shopping-list',
          to: 'shoppinglist',
          disabled: true
        },
        {
          name: 'Noticeboard',
          url: '#/noticeboard',
          to: 'noticeboard',
          disabled: true
        },
        {
          name: 'Shared Calendar',
          url: '#/shared-calendar',
          to: 'cal',
          disabled: true
        },
        {
          name: 'Recipes',
          url: '#/recipes',
          to: 'recipes',
          disabled: true
        },
        {
          name: 'Flatmates',
          url: '#/members',
          to: 'members',
          disabled: false
        },
        {
          name: 'Highfives',
          url: '#/high-fives',
          to: 'highfives',
          disabled: true
        },
        {
          name: 'Account Settings',
          url: '#/account-settings',
          to: 'account-settings',
          disabled: true
        },
        {
          name: 'Admin',
          url: '#/admin',
          to: 'admin',
          groupRequires: ['admin']
        }
      ],
      login: {
        name: 'Person',
        group: 'flatmember'
      }
    }
  },
  methods: {
    memberSignOut: () => {
      Dialog.confirm({
        message: 'Are you sure you want to sign out?',
        confirmText: 'Sign out',
        type: 'is-warning',
        onConfirm: () => {
          const loadingComponent = Loading.open({
            container: null
          })
          setTimeout(() => {
            sessionStorage.clear()
            loadingComponent.close()
            location.href = '#/login'
            console.log('Signing out')
          }, 1.2 * 1000)
        }
      })
    },
    determineItemDisplay: (item, login) => {
      if (!item.groupRequires) return true
      return item.groupRequires.includes(login.group)
    }
  },
  components: {
    headerDisplay
  }
}
</script>

<style scoped>
#menu-bar-items {
  padding-left: 10px;
  padding-top: 10px;
  color: black;
}

#menu-bar-items .menu-bar-item {
  border-top: 1px solid rgba(111, 111, 111, 0.4);
}
</style>
