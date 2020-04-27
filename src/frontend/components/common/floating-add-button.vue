<template>
  <div class="floating-add-button">
    <md-speed-dial md-direction="bottom md-bottom-right">
      <md-speed-dial-target @click="goToRef(path)" class="floating-add-button-colour">
        <md-icon>add</md-icon>
      </md-speed-dial-target>
    </md-speed-dial>
  </div>
</template>

<script>
export default {
  name: 'floating add button',
  props: {
    path: String
  },
  methods: {
    goToRef (ref) {
      this.$router.push({ path: ref })
    }
  }
}
</script>

<style scope>
.floating-add-button {
    display: block;
    position: fixed;
    bottom: 0;
    right: 0;
    margin-bottom: 50px;
    z-index: 100;
}

.floating-add-button-colour {
    background-color: #448aff;
}
</style>
