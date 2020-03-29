<template>
  <div>
    <!-- <headerDisplay/> -->
    <div class="container">
      <section class="section">
        <h2 class="title is-2">Home</h2>
        <!-- <div> -->
        <!--   <div v-for="item in pages" v-bind:key="item"> -->
        <!--     <b-button type="is-white" tag="a" :href="item.disabled === false ? `${item.url}` : ``" :disabled="item.disabled" v-if="determineItemDisplay(item, login)" size="is-large" expanded>{{ item.name }}</b-button> -->
        <!--   </div> -->
        <!-- </div> -->
        <!-- <br> -->
      </section>
    </div>
  </div>
</template>

<script>
import { LoadingProgrammatic as Loading, DialogProgrammatic as Dialog } from 'buefy'

export default {
  name: 'home',
  data () {
    return {
    }
  },
  methods: {
  },
  created () {
    // GetAPIprofile()
    //   .then(resp => {
    //     this.login = resp.data
    //   })
    //   .catch(err => {
    //     this.$buefy.notification.open({
    //       duration: 5000,
    //       message: `An error has occured: ${err}`,
    //       position: 'is-bottom-right',
    //       type: 'is-danger',
    //       hasIcon: true
    //     })
    //   })
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
