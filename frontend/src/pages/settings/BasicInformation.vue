<template>
  <div class="flex flex-col divide-y">
    <div>
      <div class="flex flex-col md:flex-row justify-between items-start md:item-center py-2 space-y-2 lg:space-y-0">
        <label for="firebase_student_id">User ID</label>
        <input
          type="text" name="firebase_student_id" id="firebase_student_id"
          :value="authState.user!.id"
          disabled class="w-full md:w-1/2 input input-bordered" />
      </div>
      <div class="flex flex-col md:flex-row justify-between items-start md:item-center py-2 space-y-2 lg:space-y-0">
        <label for="student_id">Student ID</label>
        <input
          type="text" name="student_id" id="student_id"
          :value="authState.user!.expand.details?.student_id"
          disabled class="w-full md:w-1/2 input input-bordered" />
      </div>
      <div class="flex flex-col md:flex-row justify-between items-start md:item-center py-2 space-y-2 lg:space-y-0">
        <label for="email">E-mail</label>
        <input
          type="text" name="email" id="email"
          :value="authState.user!.email"
          disabled class="w-full md:w-1/2 input input-bordered" />
      </div>
      <div class="flex flex-col md:flex-row justify-between items-start md:item-center py-2 space-y-2 lg:space-y-0">
        <label for="sex">Sex</label>
        <select :disabled="isSaving" name="sex" class="select select-bordered" v-model="sex">
          <option :value="g.value" :key="g.value" v-for="g in store.state.sexList">{{ g.label }}</option>
        </select>
      </div>
      <div class="flex flex-col md:flex-row justify-between items-start md:item-center py-2 space-y-2 lg:space-y-0">
        <label for="department">Department</label>
        <select :disabled="isSaving" name="department" class="select select-bordered" v-model="department">
          <option :value="dept.id" :key="dept.id" v-for="dept in store.state.departmentList">{{ dept.label }} ({{ dept.uid }})</option>
        </select>
      </div>

      <div class="flex items-center justify-end mb-8">
        <button 
          :disabled="isSaving" 
          @click="() => saveUserInfo({ sex, college_department: department })" 
          class="btn btn-success px-12">Save</button>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useMutation } from '@tanstack/vue-query';
import { ref } from 'vue';
import { pb } from '../../client';
import { useAuth, useStore } from '../../store_new';
import { UserDetails } from '../../types';

const store = useStore();
const { state: authState } = useAuth();

const sex = ref(pb.authStore!.model!.expand.details?.sex ?? 'unknown');
const department = ref(pb.authStore!.model!.expand.details?.department ?? 'none');
const { mutate: saveUserInfo, isLoading: isSaving } = useMutation((newDetails: { sex?: string, college_department?: string }) => {
  return pb.collection('user_details').update(authState.user!.details, newDetails);
}, {
  onSuccess(data) {
    // this.$notify({ type: 'success', text: data['message'] });
    authState.user!.expand.details = data as UserDetails;
  },
  onError(err) {
    // catchAndNotifyError(this, e);
  }
});
</script>