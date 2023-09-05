<script setup lang="ts">

</script>

<template>
  <header>
    <nav class="bg-blue-500 p-4">
      <div class="container mx-auto flex justify-between items-center">
        <div class="text-white text-xl font-semibold">
          CrossBoxFit
        </div>
        <ul class="flex space-x-4">
          <li><RouterLink to="/" class="text-white hover:text-gray-300 hover:underline">Home</RouterLink></li>
          <li><RouterLink to="/about" class="text-white hover:text-gray-300 hover:underline">About</RouterLink></li>
        </ul>
      </div>
    </nav>
  </header>
</template>

<style scoped>

</style>