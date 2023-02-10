<template>
  <modal 
    :open="store.state.isSetupModalOpen" 
    @update:open="store.state.isSetupModalOpen = $event"
    modal-box-class="max-w-[70rem] min-h-[80%] h-full">
    <div class="flex flex-col items-center h-full">
      <ul class="steps steps-horizontal text-lg">
        <li @click="step = 0" :class="{ 'step-primary': step >= 0 }" class="step">Welcome</li>
        <li @click="step = 1" :class="{ 'step-primary': step >= 1 }" class="step">Set Up</li>
        <li @click="step = 2" :class="{ 'step-primary': step >= 2 }" class="step text-lg">Terms &amp; Conditions</li>
      </ul>
      <div class="py-4 md:p-16 h-full w-full">
        <div class="flex flex-col h-full w-full">
          <div v-show="step == 0" class="h-full w-full flex flex-col">
            <div class="flex flex-col md:flex-row flex-1">
              <div class="flex flex-col text-center lg:text-right justify-center space-y-4 pr-24 md:w-2/3">
                <h2 class="text-3xl md:text-7xl font-bold">Welcome to UIC Valentine Wall!</h2>
                <p class="text-3xl">Please click "next" to get started</p>
              </div>
              <div class="md:w-1/3 flex justify-center">
                <icon-setup-welcome class="h-full w-1/2 md:w-full text-rose-400" />
              </div>
            </div>
            <button class="self-end px-12 btn bg-rose-500 hover:bg-rose-600 border-none mt-4" @click="step++">Next</button>
          </div>
          <basic-information-step
            v-show="step == 1"
            @success="onInfoFormSuccess"
            @error="onInfoFormError"
            @proceed="onHandleMove" />
          <terms-and-conditions-step
            v-show="step == 2"
            @status="onTCStatus" />
        </div>
      </div>
    </div>
  </modal>
</template>

<script lang="ts" setup>
import { catchAndNotifyError } from '../notify';
import Modal from './Modal.vue';
import BasicInformationStep from './SetupDialog/BasicInformationStep.vue';
import TermsAndConditionsStep from './SetupDialog/TermsAndConditionsStep.vue';
import IconSetupWelcome from '~icons/home-icons/setup_welcome';
import { pb } from '../client';
import { reactive, ref, watch } from 'vue';
import { useAuth, useStore } from '../store_new';
import { UserDetails } from '../types';

const store = useStore();
const { state: {user} } = useAuth();

const submitDetails = reactive<{
  student_id: string | null
  college_department: string | null
  sex: string | null
  terms_agreed: boolean
}>({
  student_id: null,
  college_department: null,
  sex: null,
  terms_agreed: false
});

const step = ref(0);

watch(step, (newVal, oldVal) => {
  if (newVal == null || newVal < 0) {
    step.value = 0;
  } else if (newVal > 2) {
    step.value = 2;
  }
});

function onHandleMove() {
  step.value++;
}

function onInfoFormError(e: unknown) {
  // catchAndNotifyError(this, e);
}

function onInfoFormSuccess(details: any) {
  if (details === null) {
    step.value--;
  }

  submitDetails.student_id = details.student_id;
  submitDetails.college_department = details.college_department;
  submitDetails.sex = details.sex;
}

function onTCStatus(status: boolean | null) {
  switch (status) {
    case true:
    case false:
      submitDetails.terms_agreed = status;
      submitSetupForm();
      break;
    case null:
      step.value--;
      break;
  }
}

async function submitSetupForm() {
  try {
    try {
      // TODO: include terms_agreed when submitting
      const { terms_agreed, ...newSubmitDetails } = submitDetails;
      const userDetails = await pb.collection('user_details').create({
        ...newSubmitDetails,
        user: user!.id
      });

      user!.expand.details = userDetails as UserDetails;

      // $notify({ type: 'success', text: 'Profile saved successfully.' });
      store.state.isSetupModalOpen = false;
      
    } catch (e) {
      // if (e instanceof APIResponseError && e.rawResponse.status == 403 && e.message == 'Access to the service is denied.') {
      //   this.$router.replace({ name: 'home-page' });
      //   await this.$store.state.dispatch('logout');
      // }
      throw e;
    }
  } catch(e) {
    // catchAndNotifyError(this, e);
  }
}
</script>