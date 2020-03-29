<template>
  <div id="app">
    <div>
      <navbar/>
    </div>
    <div class="pad-left full-height main-view-container">
      <router-view class="main-view" />
    </div>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import Component from 'vue-class-component'

export default {
  name: 'App',
  components: {
    navbar: () => import('@/frontend/components/common/navbar.vue')
  }
}

</script>

<style lang="scss">
@import "~bulma/sass/utilities/_all";

$primary: #00a7d6;
$primary-invert: findColorInvert($primary);

$colors: (
    "white": ($white, $black),
    "black": ($black, $white),
    "light": ($light, $light-invert),
    "dark": ($dark, $dark-invert),
    "primary": ($primary, $primary-invert),
    "info": ($info, $info-invert),
    "success": ($success, $success-invert),
    "warning": ($warning, $warning-invert),
    "danger": ($danger, $danger-invert)
);

$link: hsl(217, 71%, 53%);
$link-invert: $black;
$link-focus-border: $primary;
$breadcrumb-item-color: $link;
$footer-padding: 10px;

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
    margin-left: 265px;
}

.pad-top {
    margin-top: 20px;
}

.full-height {
    height: 100%;
}

html {
    background-color: #f3f3f3;
}
</style>
