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
  <b-tooltip
    :label="message"
    type="is-white"
    :position="position ? position : 'is-top'"
    multilined
  >
    <b-icon
      icon="information-outline"
      type="is-lightyellow"
      custom-size="mdi-24px"
    >
    </b-icon>
  </b-tooltip>
</template>

<script>
export default {
  name: 'password-tooltip-icon',
  props: {
    message: String,
    position: String
  }
}
</script>
