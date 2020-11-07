<!--
     This program is free software: you can redistribute it and/or modify
     it under the terms of the Affero GNU General Public License as published by
     the Free Software Foundation, either version 3 of the License, or
     (at your option) any later version.

     This program is distributed in the hope that it will be useful,
     but WITHOUT ANY WARRANTY; without even the implied warranty of
     MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
     GNU General Public License for more details.

     You should have received a copy of the Affero GNU General Public License
     along with this program.  If not, see <https://www.gnu.org/licenses/>.
-->

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
