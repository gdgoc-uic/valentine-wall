<template>
  <main class="flex flex-col max-w-5xl mx-auto px-4">
    <h1 class="text-7xl font-bold text-rose-600 text-center self-center mb-12">Settings</h1>

    <article class="bg-white rounded-xl w-full shadow-lg">
        <section class="rounded-t-xl bg-rose-50 flex pt-2">
            <div class="tabs self-center mx-auto">
                <router-link 
                    v-for="sectionRoute in settingSectionRoutes"
                    :key="sectionRoute.name"
                    :to="{ name: sectionRoute.name }" 
                    exact-active-class="tab-active" 
                    class="tab tab-lg tab-bordered"
                    exact>{{ sectionRoute.meta.label }}</router-link>
            </div>
        </section>

        <section class="p-8">
            <router-view></router-view>
        </section>
    </article>
  </main>
</template>

<script lang="ts">
import { RouteRecordRaw } from 'vue-router';
export default {
    computed: {
        settingSectionRoutes(): RouteRecordRaw[] {
            const routes = this.$router.getRoutes();
            const currentRouteData = routes.find(r => r.name === this.$route.matched[0].name);
            if (!currentRouteData) return [];
            return currentRouteData.children;
        }
    }
}
</script>