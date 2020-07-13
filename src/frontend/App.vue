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
import routerCommon from '@/frontend/router/common'

export default {
  name: 'App',
  data () {
    return {
      onMobile: false,
      displayNavigationBar: true
    }
  },
  computed: {
    publicPages () {
      return routerCommon.isPublicRoute(this.$route)
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
    "lightred": ($lightred, #b55c5c),
    "lightyellow": (#e2ab1f, #e2ab1f)
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
    background-color: #fbfbfb;
}

.pointer-cursor-on-hover {
    cursor: pointer;
}

.card {
    user-select: none;
    border: 1px dashed rgba(0, 0, 0, .0);
    transition: box-shadow 0.4s;
}

.card:hover {
    background-color: #fbfbfb;
    transition: box-shadow 0.4s, background-color 0.3s;
    box-shadow: black 0 0px 45px -33px;
}

.card:active {
    background-color: #f1f0f0;
    transition: height 0.1s, width 0.1s, margin 0.1s, box-shadow 0.4s;
    height: calc(100% - 8px);
    border: 1px dashed #000000;
    box-shadow: black 0 0px 66px -27px;
    top: 1px;
}

.hero-body, .section {
    padding: 3rem 1rem;
}

.form-width {
    width: 380px;
    margin: auto;
}

.button {
    box-shadow: #ffffff -1px -1px 4px 0px, #bababa 1px 1px 4px 0px;
}

.navbar-shadow {
    box-shadow: 0px -9px 23px 4px #292929;
}

@media (max-width : 870px) {
    .form-width {
        width: calc(100% - 25px);
    }
}
</style>
