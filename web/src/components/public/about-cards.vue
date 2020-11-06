<template>
    <section>
        <div style="padding: 1.5rem">
        <b-carousel
            v-model="carousel"
            :animated="animated"
            :has-drag="drag"
            :autoplay="autoPlay"
            :pause-hover="pauseHover"
            :pause-info="pauseInfo"
            :pause-info-type="pauseType"
            :interval="interval"
            :repeat="repeat"
            @change="info($event)">
            <b-carousel-item v-for="(carousel, i) in carousels" :key="i">
                <section :class="`hero is-medium is-${carousel.color} is-bold`">
                    <div class="hero-body has-text-centered">
                        <h1 class="title">{{carousel.title}}</h1>
                    </div>
                </section>
            </b-carousel-item>
        </b-carousel>
    </section>
</template>

<script>
export default {
  data () {
    return {
      carousel: 0,
      animated: 'fade',
      drag: false,
      autoPlay: false,
      pauseHover: false,
      pauseInfo: false,
      repeat: false,
      pauseType: 'is-primary',
      interval: 3000,
      carousels: [
        { title: 'Slide 1', color: 'dark' },
        { title: 'Slide 2', color: 'primary' },
        { title: 'Slide 3', color: 'info' },
        { title: 'Slide 4', color: 'success' },
        { title: 'Slide 5', color: 'warning' },
        { title: 'Slide 6', color: 'danger' }
      ]
    }
  },
  methods: {
    info (value) {
      this.carousel = value
      this.$buefy.toast.open({
        message: `This Slide ${value} !`,
        type: 'is-info'
      })
    }
  }
}
</script>
