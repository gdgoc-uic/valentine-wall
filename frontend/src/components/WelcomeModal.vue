<template>
  <modal :open="store.isWelcomeModalOpen" class="w-full" modal-box-class="max-w-[50rem] w-full p-12">
    <div class="flex flex-col">
      <h2 class="text-2xl lg:text-4xl text-center font-bold mb-4">Welcome to UIC Valentine Wall!</h2>
      <div class="flex flex-col md:flex-row text-lg md:text-2xl lg:text-lg">
        <div class="p-3 w-full lg:w-1/3 flex flex-row md:flex-col items-center justify-center md:justify-start md:space-y-4 text-left md:text-center">
          <div class="w-1/3 md:w-full bg-rose-400 rounded-full p-4 md:p-8 mr-8 md:mr-0">
            <icon-welcome-feature-1 class="w-full h-full text-white" />
          </div>
          <p>Post, confess, or share your thoughts to your fellow ignacian anonymously!</p>
        </div>
        <div class="p-3 w-full lg:w-1/3 flex flex-row md:flex-col items-center justify-center md:justify-start md:space-y-4 text-left md:text-center">
          <div class="w-1/3 md:w-full bg-blue-400 rounded-full p-4 md:p-8 mr-8 md:mr-0">
            <icon-welcome-feature-2 class="w-full h-full text-rose-100" />
          </div>
          <p>No money? No problem! You can also send virtual gifts alongside your message!</p>
        </div>
        <div class="p-3 w-full lg:w-1/3 flex flex-row md:flex-col items-center justify-center md:justify-start md:space-y-4 text-left md:text-center">
          <div class="w-1/3 md:w-full bg-orange-300 rounded-full p-4 md:p-8 mr-8 md:mr-0">
            <icon-welcome-feature-3 class="w-full h-full text-white" />
          </div>
          <p>Discover many more messages sent by others through our public wall!</p>
        </div>
      </div>
      <button @click="store.toggleWelcomeModal()" class="btn btn-primary btn-lg self-center px-12 mt-8">
        Proceed
      </button>
    </div>
  </modal>
</template>

<script lang="ts" setup>
import Modal from './Modal.vue';
import IconWelcomeFeature1 from '~icons/home-icons/welcome_feature_1';
import IconWelcomeFeature2 from '~icons/home-icons/welcome_feature_2';
import IconWelcomeFeature3 from '~icons/home-icons/welcome_feature_3';
import { useStore } from '../store_new';

const store = useStore();
</script>