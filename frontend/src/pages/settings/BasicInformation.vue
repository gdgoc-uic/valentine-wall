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
        <select :disabled="isSaving" name="sex" class="select select-bordered" v-model="sex">
          <option :value="g.value" :key="g.value" v-for="g in $store.getters.sexList">{{ g.label }}</option>
        </select>
      </div>
      <div class="flex flex-col md:flex-row justify-between items-start md:item-center py-2 space-y-2 lg:space-y-0">
        <label for="department">Department</label>
        <select :disabled="isSaving" name="department" class="select select-bordered" v-model="department">
          <option :value="dept.id" :key="dept.id" v-for="dept in $store.state.departmentList">{{ dept.label }} ({{ dept.id }})</option>
        </select>
      </div>

      <div class="flex items-center justify-end mb-8">
        <button :disabled="isSaving" @click="saveUserInfo" class="btn btn-success px-12">Save</button>
      </div>
    </div>
    <div class="py-6">
      <h2 class="font-bold text-2xl">Connections</h2>
      <p>Third-party accounts for use on sharing content outside the site such as replies and etc.</p>

      <div class="flex flex-col">
        <div v-if="$store.state.user.connections.findIndex(c => c.provider === 'twitter') === -1" class="flex justify-between items-center py-2">
          <label>twitter</label>
          <button 
            @click="connectUserConnection('twitter')" 
            class="btn btn-success btn-outline">Connect</button>
        </div>
        <div v-if="$store.state.user.connections.findIndex(c => c.provider === 'email') === -1" class="flex justify-between items-center py-2">
          <label>email</label>
          <button 
            @click="connectUserConnection('email')" 
            class="btn btn-success btn-outline">Connect</button>
        </div>
        <div 
          v-for="(conn, i) in $store.state.user.connections" 
          :key="'conn_' + i" class="flex justify-between items-center py-2">
          <label>{{ conn.provider }}</label>
          <button 
            :disabled="conn.provider == 'email'" 
            @click="disconnectUserConnection(conn.provider)" 
            class="btn btn-error btn-outline">Disconnect</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { connectToEmail, connectToTwitter } from '../../auth';
import { catchAndNotifyError } from '../../notify';
export default {
  mounted() {
    // TODO:
    this.sex = this.$store.state.user.sex;
    this.department = this.$store.state.user.department;
  },
  data() {
    return {
      sex: 'unknown',
      department: 'none',
      isSaving: false
    }
  },
  methods: {
    async saveUserInfo() {
      try {
        this.isSaving = true;
        // const { data } = await this.$client.patchJson('/user/info', {
        //   sex: this.sex,
        //   department: this.department
        // });

        // this.$notify({ type: 'success', text: data['message'] });
        // await this.$store.dispatch('getUserInfo');
        // this.sex = this.$store.state.user.sex;
        // this.department = this.$store.state.user.department;
      } catch (e) {
        catchAndNotifyError(this, e);
      } finally {
        this.isSaving = false;
      }
    },
    async connectUserConnection(provider: string) {
      try {
        switch (provider) {
          case "twitter":
            connectToTwitter(this.$store);
            break;
          case "email":
            connectToEmail(this.$store);
            break;
          default:
            throw new Error(`Unknown provider: ${provider}`);
        }
        await this.$store.dispatch('getUserInfo');
      } catch (e) {
        catchAndNotifyError(this, e);
      }
    },
    async disconnectUserConnection(name: string) {
      try {
        // const { data } = await this.$client.delete(`/user/connections/${name}`);
        // this.$notify({ type: 'success', text: data['message'] });
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