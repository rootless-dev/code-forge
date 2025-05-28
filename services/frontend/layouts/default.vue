<script lang="ts" setup>

import IndexHeaderBar from "~/components/ui/headers/IndexHeaderBar.vue";
</script>

<template>
  <div>
    <header>
      <IndexHeaderBar/>
    </header>

    <main class="flex-1">
      <slot/>
    </main>

    <footer>
      <UContainer>
        all rights reserved
      </UContainer>
    </footer>
  </div>
</template>

<style scoped>

</style>