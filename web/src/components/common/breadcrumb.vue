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
  <b-breadcrumb separator="has-arrow-separator">
    <b-breadcrumb-item tag="router-link" :to="{ name: backLinkName }">
      {{ backLinkName }}
    </b-breadcrumb-item>
    <b-breadcrumb-item tag="router-link" :to="{ name: currentPageName }" active>
      {{ currentPageNameOverride || currentPageName }}
    </b-breadcrumb-item>
    <b-button
      icon-left="content-copy"
      size="is-small"
      @click="CopyHrefToClipboard()"
    />
  </b-breadcrumb>
</template>

<script>
  import common from "@/common/common";

  export default {
    name: "BreadcrumbNav",
    props: {
      backLinkName: String,
      currentPageName: String,
      currentPageNameOverride: String,
    },
    data() {},
    methods: {
      CopyHrefToClipboard() {
        common.CopyHrefToClipboard();
      },
    },
  };
</script>
