<template>
  <div class="flex flex-col divide-y">
    <div>
      <div class="flex flex-col md:flex-row justify-between items-start md:item-center py-2 space-y-2 lg:space-y-0">
        <label for="firebase_student_id">Firebase User ID</label>
        <input
          type="text" name="firebase_student_id" id="firebase_student_id"
          :value="$store.state.user.id"
          disabled class="w-full md:w-1/2 input input-bordered" />
      </div>
      <div class="flex flex-col md:flex-row justify-between items-start md:item-center py-2 space-y-2 lg:space-y-0">
        <label for="student_id">Student ID</label>
        <input
          type="text" name="student_id" id="student_id"
          :value="$store.state.user.associatedId"
          disabled class="w-full md:w-1/2 input input-bordered" />
      </div>
      <div class="flex flex-col md:flex-row justify-between items-start md:item-center py-2 space-y-2 lg:space-y-0">
        <label for="email">E-mail</label>
        <input
          type="text" name="email" id="email"
          :value="$store.state.user.email"
          disabled class="w-full md:w-1/2 input input-bordered" />
      </div>
      <div class="flex flex-col md:flex-row justify-between items-start md:item-center py-2 space-y-2 lg:space-y-0">
        <label for="sex">Sex</label>
        <select name="sex" class="select select-bordered">
          <option :selected="g.value == $store.state.user.sex" :value="g.value" :key="g.value" v-for="g in sexList">{{ g.label }}</option>
        </select>
      </div>
      <div class="flex flex-col md:flex-row justify-between items-start md:item-center py-2 space-y-2 lg:space-y-0">
        <label for="department">Department</label>
        <select name="department" class="select select-bordered">
          <option value="none" selected>None</option>
          <option :selected="dept.id == $store.state.user.department" :value="dept.id" :key="dept.id" v-for="dept in $store.state.departmentList">{{ dept.label }} ({{ dept.id }})</option>
        </select>
      </div>
    </div>
    <div class="py-6">
      <p class="font-bold">Connections</p>
      <div class="flex flex-col">
        <div 
          v-for="(conn, i) in $store.state.user.connections" 
          :key="'conn_' + i" class="flex justify-between items-center py-2">
          <label>{{ conn.provider }}</label>
          <button 
            :disabled="conn.provider == 'email'" 
            @click="disconnectUserConnection(conn.provider)" 
            class="btn btn-error btn-outline">Disconnect</button>
        </div>
        <p class="py-4 text-gray-600" v-if="$store.state.user.connections.length == 0">No connections found.</p>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { catchAndNotifyError } from '../../notify';
export default {
  mounted() {
    // TODO:
  },
  computed: {
    sexList() {
      return [
        {
          label: 'Male',
          value: 'male'
        },
        {
          label: 'Female',
          value: 'female'
        }
      ];
    }
  },
  methods: {
    async disconnectUserConnection(name: string) {
      try {
        const { data } = await this.$client.delete(`/user/connections/${name}`);
        this.$notify({ type: 'success', text: data['message'] });
      } catch (e) {
        catchAndNotifyError(this, e);
      } finally {
        try {
          await this.$store.dispatch('getUserInfo');
        } catch (e) {
          catchAndNotifyError(this, e);
        }
      }
    }
  }
}
</script>