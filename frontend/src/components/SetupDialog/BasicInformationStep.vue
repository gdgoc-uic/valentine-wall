<template>
  <div class="flex flex-col h-full">
    <div class="flex flex-col md:flex-row my-auto items-center flex-1">
      <div class="md:w-1/2 text-center md:text-left">
        <h2 class="text-5xl font-bold">Account Set-up</h2>
        <p class="text-lg mt-4">Set up important information such as your student ID, department, and more</p>
      </div>

      <form ref="form" @submit.prevent="shouldProceed" class=" w-1/2 flex flex-col">
        <div class="form-control">
          <label class="label">
            <span class="label-text">Enter your student ID</span>
          </label>
          <input class="input input-bordered" type="text" name="student_id" id="student_id_field" pattern="[0-9]{6,12}" placeholder="6 to 12-digit Student ID (e.g. 200xxxxxxxxx)">
        </div>
        <div class="form-control">
          <label class="label">
            <span class="label-text">College Deparment</span>
          </label>
          <select name="department" class="select select-bordered">
            <option value="none" selected>None</option>
            <option :value="dept.id"
              :key="dept.id"
              v-for="dept in store.state.departmentList">
              {{ dept.label }} ({{ dept.uid }})
            </option>
          </select>
        </div>
        <div class="form-control">
          <label class="label">
            <span class="label-text">Sex</span>
          </label>
          <select name="sex" class="select select-bordered">
            <!-- TODO: -->
            <!-- <option
              :value="g.value.toLowerCase()"
              :key="g.value"
              v-for="g in $store.state.getters.sexList">
              {{ g.label }}
            </option> -->
          </select>
        </div>
      </form>
    </div>
    <div class="flex flex-row w-full justify-between mt-4">
      <button @click="$emit('success', null)" class="px-12 btn">Go back</button>
      <button  
        @click="submitForm" 
        class="px-12 btn bg-rose-500 hover:bg-rose-600 border-none" type="submit">Next</button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import { pb } from '../../client';
import { useAuth, useStore } from '../../store_new';

const emailRegex = /^[a-z]+_([0-9]+)@uic.edu.ph$/;
const emit = defineEmits(['success', 'error', 'proceed']);
const form = ref<HTMLFormElement | null>();
const store = useStore();
const { state: {user} } = useAuth();

function getIdFromEmail(input: string): string {
  const matches = emailRegex.exec(input)
  return matches?.[1] ?? '';
}

function submitForm() {
  form.value?.dispatchEvent(new Event('submit', {'bubbles': true, 'cancelable': true }))
}

function shouldProceed(e: SubmitEvent) {
  if (!e.target || !(e.target instanceof HTMLFormElement)) return;
  const formData = new FormData(e.target);
  if (!formData.get("student_id")) {
    emit('error', new Error('Please input your ID.'));
    return;
  } else if (formData.get('student_id')?.toString() !== getIdFromEmail(user?.email)) {
    emit('error', new Error('Your ID from e-mail does not match with the one you have inputted.'));
    return;
  }

  if (!formData.get('department') || formData.get('department')?.toString() == 'none') {
    emit('error', new Error('Please select your department.'));
    return;
  }

  emit('success', {
    student_id: formData.get('student_id')?.toString(),
    college_department: formData.get('department')?.toString(),
    sex: formData.get('sex')?.toString(),
  });

  emit('proceed');
}

onMounted(() => {
  setTimeout(() => {
    const associatedIdField = document.getElementById('student_id_field');
    if (associatedIdField && associatedIdField instanceof HTMLInputElement) {
      const extractedId = getIdFromEmail(user?.email);
      associatedIdField.value = extractedId;
    }
  }, 500);
});
</script>