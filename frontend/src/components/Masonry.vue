<template>
  <div>
    <slot></slot>
  </div>
</template>

<script lang="ts">
//@ts-expect-error
import Masonry from 'masonry-layout?client';
import { PropType } from '@vue/runtime-core';

interface Config extends Masonry.Options {
  destroyDelay?: string | number
}

export default {
  props: {
    config: {
      type: Object as PropType<Config>,
      default: {
        destroyDelay: 0
      }
    }
  },
  data() {
    return {
      masonry: null as unknown as Masonry,
      observer: null as unknown as MutationObserver
    }
  },
  mounted() {
    if (!import.meta.env.SSR) {
      this.observer = new MutationObserver(this.draw);
  
      const el: HTMLDivElement = <HTMLDivElement> this.$el;
      this.masonry = new Masonry(el, this.config);
  
      this.observer.observe(this.$el, {
        childList: true,
        subtree: true
      });
  
      this.$nextTick(this.draw);
    }
  },
  beforeUnmount() {
    if (!import.meta.env.SSR) {
      this.destroy();
      this.observer.disconnect();
    }
  },
  computed: {
    destroyDelay(): number {
      return typeof this.config.destroyDelay == 'string' 
        ? parseInt(this.config.destroyDelay, 10) 
        : typeof this.config.destroyDelay == 'number'
        ? this.config.destroyDelay
        : 0;
    },
  },
  methods: {
    draw() {
      if (this.masonry.reloadItems) {
        this.masonry.reloadItems();
      }

      if (this.masonry.layout) {
        this.masonry.layout();
      }
    },
    destroy() {
      setTimeout(this.masonry.destroy || function() {}, this.destroyDelay);
    }
  }
}
</script>