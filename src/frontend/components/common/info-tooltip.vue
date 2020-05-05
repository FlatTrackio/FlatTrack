<template>
  <b-tooltip
    :label="message"
    type="is-white"
    position="is-top"
    multilined>
    <b-icon
      icon="information-outline"
      type="is-black"
      custom-size="mdi-24px">
    </b-icon>
  </b-tooltip>
</template>

<script>
export default {
  name: 'password-tooltip-icon',
  props: {
    message: String
  }
}
</script>
