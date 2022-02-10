<template>
  <main class="flex flex-col max-w-5xl mx-auto px-4">
    <h1 class="text-7xl font-bold text-rose-600 text-center self-center mb-12">Settings</h1>

    <article class="bg-white rounded-xl w-full shadow-lg">
        <section class="rounded-t-xl bg-rose-50 flex pt-2 md:pt-0">
            <ul class="menu menu-vertical w-full md:w-auto md:menu-horizontal md:self-center md:mx-auto">
                <router-link 
                    v-for="sectionRoute in settingSectionRoutes"
                    :key="sectionRoute.name"
                    :to="{ name: sectionRoute.name }" 
                    v-slot="{ href, route, navigate, isExactActive }"
                    custom>
                    <li :class="{ 'bordered': isExactActive }">
                        <a @click="navigate" :href="href">{{ sectionRoute.meta.label }}</a>
                    </li>
                </router-link>
            </ul>
        </section>

        <section class="px-8 lg:px-24 py-8">
            <router-view v-slot="{ Component }">
                <keep-alive>
                    <component :is="Component" />
                </keep-alive>
            </router-view>
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