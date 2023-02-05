<template>
  <div class="flex-col">
    <div>
      <h2 class="text-center text-2xl mb-5">Terms and Conditions</h2>
      <div class="bg-gray-200 p-8 rounded-md max-h-[30rem] overflow-y-scroll overflow-x-hidden whitespace-pre-wrap">
        {{ tc }}
      </div>
      <p class="mt-4">By clicking "Agree", you have agreed to the terms and conditions of this site. Should you violate any of the text above will result to account termination.</p>
    </div>
    <div class="w-full space-x-2 mt-4 flex">
      <button @click.prevent="$emit('status', null)" class="px-6 btn mr-auto">Go back</button>
      <button :disabled="!tc" @click="$emit('status', true)" class="px-6 btn bg-rose-500 hover:bg-rose-600 border-none" type="submit">Agree</button>
      <button :disabled="!tc" @click="$emit('status', false)" class="btn px-6">Disagree</button>
    </div>
  </div>
</template>

<script>
import { catchAndNotifyError } from '../../notify'
export default {
  emits: ['status'],
  mounted() {
    this.loadTC();
  },
  data() {
    return {
      tc: null,
    }
  },
  methods: {
    async loadTC() {
      try {
        // const { data } = await this.$client.get('/terms-and-conditions');
        // console.log(data);
        // this.tc = data;
      } catch (e) {
        catchAndNotifyError(this, e);
      }
    }
  }
}
</script>