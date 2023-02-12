<template>
  <slot :openDialog="openDialog"></slot>
</template>

<script lang="ts" setup>
import { popupCenter } from '../auth';
import { useAuth } from '../store_new';

const { state: authState } = useAuth();
const props = defineProps({
  email: {
    type: String,
    default: ''
  },
  link: {
    type: String,
    required: true
  }
});

function openDialog() {
  const encodedEmail = encodeURIComponent(props.email);
  const encodedId = encodeURIComponent(authState.user.id);
  const encodedLink = encodeURIComponent(props.link);

  popupCenter({
    url: `https://docs.google.com/forms/d/e/1FAIpQLSdqGgcB42NflCWEJ4H-N5xE5_oKNQ_lNBYvYSeobeJGC9AWVQ/viewform?usp=pp_url&entry.700984842=${encodedEmail}&entry.1940307198=${encodedId}&entry.1779754379=${encodedLink}&entry.151988683=Option+1`,
    title: 'Report form',
    w: 800,
    h: 800
  });
}
</script>