<template>
  <div id="app">
    <router-view/>
  </div>
</template>

<script>
export default {
  name: 'App'
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

$link: hsl(217, 71%, 53%);;
$link-invert: $black;
$link-focus-border: $primary;
$breadcrumb-item-color: $link;

@import "~bulma";
@import "~buefy/src/scss/buefy";

</style>

<style scoped>
.container {
  margin-top: 20px;
}

.darken a.is-disabled {
  color: black;
}
</style>
