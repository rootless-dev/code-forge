<script lang="ts" setup>
import type {NavigationMenuItem} from "#ui/components/NavigationMenu.vue";
import {useMediaQuery} from "@vueuse/core";

const props = defineProps({
  sections: Array<NavigationMenuItem>,
})

const items = computed(() => [
  props.sections?.map(section => ({
    label: section.label,
    to: section.to,
  })) || [],

])

const isSmallScreen = useMediaQuery('(max-width: 767px)')

const orientation = computed(() => {
  return isSmallScreen.value ? 'vertical' : 'horizontal'
})
</script>

<template>
  <UContainer class="fixed top-0 left-0 min-w-full flex justify-around items-center p-4 bg-gray-900/40 backdrop-blur">
    <div>
      <NuxtImg
          height="40"
          src="https://static.vecteezy.com/system/resources/thumbnails/024/553/534/small_2x/lion-head-logo-mascot-wildlife-animal-illustration-generative-ai-png.png"
          width="40"
      />
    </div>
    <UNavigationMenu
        :items="items"
        :orientation="orientation"
        highlight
        highlight-color="primary"
        variant="link"
    />

    <div>
      <div class="flex gap-2">
        <UButton class="cursor-pointer" to="/login">Sign In</UButton>
        <UButton to="/register" variant="outline">Sign Up</UButton>
      </div>
    </div>
  </UContainer>

</template>

<style scoped>

</style>