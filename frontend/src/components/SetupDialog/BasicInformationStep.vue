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
              v-for="dept in $store.state.departmentList">
              {{ dept.label }} ({{ dept.uid }})
            </option>
          </select>
        </div>
        <div class="form-control">
          <label class="label">
            <span class="label-text">Sex</span>
          </label>
          <select name="sex" class="select select-bordered">
            <option
              :value="g.value.toLowerCase()"
              :key="g.value"
              v-for="g in $store.getters.sexList">
              {{ g.label }}
            </option>
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

<script lang="ts">
import { pb } from '../../client';

const emailRegex = /^[a-z]+_([0-9]+)@uic.edu.ph$/;

export default {
  emits: ['success', 'error', 'proceed'],
  mounted() {
    setTimeout(() => {
      const associatedIdField = document.getElementById('student_id_field');
      if (associatedIdField && associatedIdField instanceof HTMLInputElement) {
        const extractedId = this.getIdFromEmail(pb.authStore.model?.email);
        associatedIdField.value = extractedId;
      }
    }, 500);
  },
  methods: {
    getIdFromEmail(input: string): string {
      const matches = emailRegex.exec(input)
      return matches?.[1] ?? '';
    },
    submitForm() {
      //@ts-ignore
      this.$refs.form.dispatchEvent(new Event('submit', {'bubbles': true, 'cancelable': true }))
    },
    shouldProceed(e: SubmitEvent) {
      if (!e.target || !(e.target instanceof HTMLFormElement)) return;
      const formData = new FormData(e.target);
      if (!formData.get("student_id")) {
        this.$emit('error', new Error('Please input your ID.'));
        return;
      } else if (formData.get('student_id')?.toString() !== this.getIdFromEmail(pb.authStore.model?.email)) {
        this.$emit('error', new Error('Your ID from e-mail does not match with the one you have inputted.'));
        return;
      }

      if (!formData.get('department') || formData.get('department')?.toString() == 'none') {
        this.$emit('error', new Error('Please select your department.'));
        return;
      }

      this.$emit('success', {
        student_id: formData.get('student_id')?.toString(),
        college_department: formData.get('department')?.toString(),
        sex: formData.get('sex')?.toString(),
      });

      this.$emit('proceed');
    },
  }
}
</script>