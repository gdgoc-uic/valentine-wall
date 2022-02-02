<template>
  <div :style="{width: size, height: size}" class="loading-wrapper">
      <div
        :style="{width: size, height: `calc(${size} - 43px)`}"
        class="loading"></div>
  </div>
</template>

<script>
export default {
    props: {
        size: {
            default: '8rem'
        },
    }
}
</script>

<style lang="postcss" scoped>
@keyframes loadingAnimation {
    0% {
        transform: translateY(-40px);
        opacity: 1;
    }
    75% {
        transform: scaleX(70%);
    }
    100% {
        transform: translateY(0%) scaleY(60%);
        opacity: 0.7;
    }
}

.loading-wrapper {
    display: flex;
    flex-direction: column;
    justify-content: flex-end;
}

.loading {
    background-image: url(../assets/images/icon.png);
    background-size: contain;
    background-repeat: no-repeat;
    background-position: center;
    animation: loadingAnimation 700ms cubic-bezier(0.4, 0, 1, 1) alternate-reverse;
    animation-iteration-count: infinite;
}
</style>