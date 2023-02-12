<template>
  <slot :openDialog="openDialog"></slot>
</template>

<script lang="ts" setup>
import { popupCenter } from '../auth';

function openDialog() {
  const formLink = `https://docs.google.com/forms/d/e/1FAIpQLScz0BFFvJ40WhUEp4aDyY9cr4IeDQgTdECB-zZ2mv3J8x1JWg/viewform?usp=pp_url`;
  popupCenter({ url: formLink, title: 'Feedback form', w: 500, h: 800  });
}
</script>