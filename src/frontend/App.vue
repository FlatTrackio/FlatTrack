<template>
  <div id="app">
    <div>
      <topbar v-if="onMobile && !publicPages" />
      <navbar v-if="!onMobile && !publicPages && displayNavigationBar" />
    </div>
    <div class="pad-bottom full-height main-view-container" :class="{ 'pad-left': !publicPages && !onMobile && displayNavigationBar }">
      <router-view class="main-view" />
    </div>
    <div>
      <bottombar v-if="onMobile && !publicPages" />
    </div>
  </div>
</template>

<script>
import Vue from 'vue'
import Component from 'vue-class-component'
import themes from '@/frontend/common/theme'

export default {
  name: 'App',
  data () {
    return {
      onMobile: false,
      displayNavigationBar: true,
      publicPages: window.location.pathname === '/login' || window.location.pathname === '/setup' || (window.location.pathname.split('/')[1] === 'useraccountconfirm' && typeof window.location.pathname.split('/')[2] !== 'undefined') || window.location.pathname === '/forgot-password'
    }
  },
  async beforeMount () {
    themes.SetThemeDefault()
  },
  created () {
    this.adjustForMobile()
    window.addEventListener('resize', this.adjustForMobile.bind(this))
  },
  components: {
    navbar: () => import('@/frontend/components/common/navbar.vue'),
    topbar: () => import('@/frontend/components/common/topbar.vue'),
    bottombar: () => import('@/frontend/components/common/bottombar.vue')
  },
  methods: {
    adjustForMobile () {
      var vm = this
      vm.onMobile = window.innerWidth <= 870
      if (vm.onMobile) {
        vm.displayNavigationBar = false
      } else {
        vm.displayNavigationBar = true
      }
    }
  }
}

</script>

<style lang="scss">
@import "~bulma/sass/utilities/_all";

$material-icons-font-path: '~material-icons/iconfont/';

@import '~material-icons/iconfont/material-icons.scss';

@import '~@mdi/font/css/materialdesignicons.min.css';

$midgray: #c9c9c9;
$lightred: #b55c5c;
$primary: #00a7d6;
$primary-invert: findColorInvert($primary);

$colors: (
    "white": ($white, $black),
    "black": ($black, $white),
    "light": ($light, $light-invert),
    "midgray": ($midgray, #c9c9c9),
    "dark": ($dark, $dark-invert),
    "primary": ($primary, $primary-invert),
    "info": ($info, $info-invert),
    "success": ($success, $success-invert),
    "warning": ($warning, $warning-invert),
    "danger": ($danger, $danger-invert),
    "lightred": ($lightred, #b55c5c)
);

$link: hsl(217, 71%, 53%);
$link-invert: $black;
$link-focus-border: $primary;
$breadcrumb-item-color: $link;

@import "~bulma";
@import "~buefy/src/scss/buefy";

</style>

<style>
.darken a.is-disabled {
    color: black;
}

.card-margin {
    margin-top: 10px;
    margin-bottom: 10px;
}

.pad-left {
    margin-left: 260px;
}

.pad-top {
    margin-top: 20px;
}

.pad-bottom {
    margin-bottom: 40px;
}

.full-height {
    height: 100%;
}

html {
    background-color: #f3f3f3;
}

.pointer-cursor-on-hover {
    cursor: pointer;
}

.card {
    user-select: none;
    border: 1px dashed rgba(0, 0, 0, .0);
}

.card:hover {
    background-color: #fbfbfb;
    transition: box-shadow 0.4s;
    box-shadow: black 0 0px 45px -33px;
    transition: box-shadow 0.1s;
}

.card:active {
    background-color: #f1f0f0;
    transition: height 0.1s, width 0.1s, margin 0.1s;
    height: calc(100% - 8px);
    border: 1px dashed #000000;
    box-shadow: black 0 0px 66px -27px;
    top: 1px;
}

.form-width {
    width: 380px;
    margin: auto;
}

@media (max-width : 870px) {
    .form-width {
        width: calc(100% - 25px);
    }
}
</style>
