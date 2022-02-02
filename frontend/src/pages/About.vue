<template>
  <main class="px-0 lg:px-2">
    <div class="max-w-7xl mx-auto flex rounded-lg shadow-lg h-[70vh]">
        <div class="bg-rose-50 rounded-l-lg w-1/2 p-8 lg:p-24 flex flex-col justify-center items-center">
            <img src="../assets/images/logo.png" alt="Valentine Wall" class="mb-20">
            <p class="font-bold text-center">Powered by:</p>
            <div class="flex w-full items-center justify-center">
                <img src="../assets/images/usg.png" class="w-1/3" alt="University Student Government - UIC">
                <img src="../assets/images/gdscuic.png" class="w-1/2" alt="Google Developer Student Clubs - UIC">
            </div>
        </div>
        <div class="bg-white rounded-r-lg w-1/2 p-16 overflow-y-auto credits">
            <img src="../assets/images/icon-big.png" class="w-1/2 text-center mx-auto mb-8" alt="Icon">
            <about-text />
        </div>
    </div>
  </main>
</template>

<script>
import { VueComponent as AboutText } from '../assets/texts/about.md';

export default {
    components: {
        AboutText
    }
}
</script>

<style lang="postcss">
.credits p {
    @apply mb-8;
}
</style>