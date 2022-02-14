<template>
  <main class="px-0 lg:px-2">
    <div class="max-w-7xl mx-auto flex flex-col lg:flex-row rounded-lg lg:shadow-lg lg:h-[70vh]">
        <div class="lg:bg-rose-50 rounded-l-lg w-full lg:w-1/2 p-8 lg:p-24 flex flex-col justify-center items-center">
            <img src="../assets/images/logo.png" alt="Valentine Wall" class="mb-20">
            <div class="flex flex-row">
                <div class="flex-1">
                    <p class="font-bold text-center mb-8">Powered by:</p>
                    <div class="flex w-full items-center justify-center">
                        <img src="../assets/images/gdscuic.png" class="w-full" alt="Google Developer Student Clubs - UIC">
                    </div>
                </div>
                <div class="flex-1">
                    <p class="font-bold text-center">In partnership with:</p>
                    <div class="flex w-full items-center justify-center">
                        <img src="../assets/images/usg.png" class="w-2/3" alt="University Student Government - UIC">
                    </div>
                </div>
            </div>
        </div>
        <div class="bg-white rounded-r-lg w-full lg:w-1/2 py-16 px-8 lg:p-16 overflow-y-auto credits">
            <img src="../assets/images/icon-big.png" class="w-1/3 text-center mx-auto mb-8" alt="Icon">
            <about-text class="text-center mb-8" />
            <div class="flex flex-row justify-center flex-wrap">
                <div 
                    :key="'member_' + i" v-for="(m, i) in members" class="p-4 flex flex-col items-center w-1/2">
                    <img 
                        :alt="m.name" 
                        :src="'/about/' + m.firstName + '.jpg'" 
                        class="rounded-full mb-3 border-4 border-rose-300 hover:border-rose-500 transition-colors" />
                    <p class="text-rose-500 font-semibold text-xl mb-1">{{ m.name }}</p>
                    <span class="text-gray-600">{{ m.org }}</span>
                    <p class="text-center">{{ m.roles.join(', ') }}</p>
                </div>
            </div>
        </div>
    </div>
  </main>
</template>

<script>
import { VueComponent as AboutText } from '../assets/texts/about.md';
import { roles, org, members } from '../assets/about.json'

export default {
    components: {
        AboutText
    },
    computed: {
        members() {
            return members.map(m => {
                return {
                    name: m.name,
                    firstName: m.name.split(' ')[0],
                    org: org[m.org_id],
                    roles: m.role_ids.map(r => roles[r])
                }
            })
        }
    },
}
</script>