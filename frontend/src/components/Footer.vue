<template>
  <footer class="text-center mt-24 pt-18 pb-24 flex flex-col items-center space-y-3 max-w-2xl mx-auto">    
    <div class="px-8 w-full">
      <img src="../assets/images/icon-big.png" class="w-1/3 text-center mx-auto mb-8" alt="Icon">
      <div class="mb-8">
        <about-text class="text-center" />
      </div>
      <div class="flex flex-row justify-center flex-wrap">
        <div class="p-2 flex flex-col items-center w-1/5 md:w-1/6" :key="'member_' + i" v-for="(m, i) in members">
          <tooltip>
            <div>
                <img 
                  :alt="m.name" 
                  :src="'/about/' + m.firstName + '.jpg'" 
                  class="rounded-full mb-3 border-4 border-rose-300 hover:border-rose-500 transition-colors" />
            </div>

            <template #popper>
              <div class="text-center flex flex-col">
                <span class="font-bold">{{  m.name  }}</span>
                <span>{{ m.org }}</span>
                <p class="text-center">{{ m.roles.join(', ') }}</p>
              </div>
            </template>
          </tooltip>
        </div>
      </div>
    </div>

    <div class="w-full p-4 md:hidden">
      <div class=" rounded-lg p-8 bg-rose-200 text-center space-y-4">
        <p class="text-2xl"><b>Let's not make your Valentine's ruined by bugs.</b> </p>
        
        <feedback-form v-slot="{ openDialog }">
          <div class="tooltip" data-tip="Problems, suggestions? Post it here!">
            <button @click="openDialog" class="btn btn-primary space-x-2">
              <icon-comment-add />
              <span>Add your Feedback</span>
            </button>
          </div>
        </feedback-form>
      </div>
    </div>

    <div class="mx-auto max-w-md w-full flex flex-col items-center space-y-4 pt-8">
      <img src="../assets/images/gdscuic.png" class="w-48" alt="Google Developer Student Clubs - University of the Immaculate Conception Chapter">
      <p>Copyright &copy; 2022, 2023 <br /> A project of Google Developer Student Clubs - University of the Immaculate Conception Chapter.</p>
      
      <ul class="flex self-center space-x-3">
        <li><a class="text-rose-500 hover:underline" href="https://www.uic.edu.ph/">UIC Website</a></li>
        <li><a class="text-rose-500 hover:underline" href="https://www.uic.edu.ph/privacy-policy/">UIC Privacy Policy</a></li>
        <li><a class="text-rose-500 hover:underline" href="https://facebook.com/dscuic">GDSC-UIC Facebook Page</a></li>
      </ul>
    </div>
  </footer>
</template>

<script lang="ts" setup>
import FeedbackForm from "./FeedbackForm.vue";
import { Tooltip } from 'floating-vue';
import { VueComponent as AboutText } from '../assets/texts/about.md';
import IconCommentAdd from "~icons/uil/comment-add";;
import { roles, org, members as rawMembers } from '../assets/about.json'
const members = rawMembers.map(m => ({
  name: m.name,
  firstName: m.name.split(' ')[0],
  org: org[m.org_id],
  roles: m.role_ids.map(r => roles[r])
}));
</script>