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
    <div class="py-6">
      <h2 class="font-bold text-2xl">Connections</h2>
      <p>Third-party accounts for use on sharing content outside the site such as replies and etc.</p>

      <div class="flex flex-col">
        <!-- TODO: work on user connections -->
        <!-- <div v-if="user.connections.findIndex(c => c.provider === 'twitter') === -1" class="flex justify-between items-center py-2">
          <label>twitter</label>
          <button 
            @click="connectUserConnection('twitter')" 
            class="btn btn-success btn-outline">Connect</button>
        </div>
        <div v-if="user.connections.findIndex(c => c.provider === 'email') === -1" class="flex justify-between items-center py-2">
          <label>email</label>
          <button 
            @click="connectUserConnection('email')" 
            class="btn btn-success btn-outline">Connect</button>
        </div>
        <div 
          v-for="(conn, i) in user.connections" 
          :key="'conn_' + i" class="flex justify-between items-center py-2">
          <label>{{ conn.provider }}</label>
          <button 
            :disabled="conn.provider == 'email'" 
            @click="disconnectUserConnection(conn.provider)" 
            class="btn btn-error btn-outline">Disconnect</button>
        </div> -->
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useMutation } from '@tanstack/vue-query';
import { ref } from 'vue';
import { connectToEmail, connectToTwitter } from '../../auth';
import { pb } from '../../client';
import { catchAndNotifyError } from '../../notify';
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

const { mutateAsync: connectUserConnection } = useMutation((provider: string) => {
  switch (provider) {
    case "twitter":
      // connectToTwitter(store);
      break;
    case "email":
      // connectToEmail(store);
      break;
    default:
      throw new Error(`Unknown provider: ${provider}`);
  }

  // TODO:
  // return store.state.dispatch('getUserInfo');
  return Promise.resolve(0);
}, {
  onError(error) {
    // catchAndNotifyError(this, e);
  },
});

const { mutateAsync: disconnectUserConnection } = useMutation((name: string) => {
  // TODO:
  // return await this.$client.delete(`/user/connections/${name}`);
  return Promise.resolve(name);
}, {
  onSuccess() {
    // this.$notify({ type: 'success', text: data['message'] });
  },
  onError(error) {
    // catchAndNotifyError(this, e);
  },
  onSettled() {
    // TODO:
    // store.state.dispatch('getUserInfo');
  }
})
</script>