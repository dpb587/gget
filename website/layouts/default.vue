<template>
  <div>
    <Nuxt />
    <div class="max-w-sm border-t border-gray-300 mx-auto text-center text-gray-500 mt-12 py-4">
      <p>gget <fa icon="signature" fixed-width class="text-gray-400" /> danny berger</p>
      <p class="mt-2">
        <a href="https://dpb587.github.io/gget/" class="mx-1 text-gray-600 hover:text-gray-700 p-2">
          <span class="sr-only">Home</span>
          <fa icon="home" fixed-width />
        </a>
        <a href="https://github.com/dpb587/gget" class="mx-1 text-gray-600 hover:text-gray-700 p-2">
          <span class="sr-only">Repository</span>
          <fa icon="code" fixed-width />
        </a>
        <a href="https://github.com/dpb587" class="mx-1 text-gray-600 hover:text-gray-700 p-2">
          <span class="sr-only">Personal GitHub</span>
          <fa :icon="['fab', 'github']" fixed-width />
        </a>
        <a href="https://twitter.com/dpb587" class="mx-1 text-gray-600 hover:text-gray-700 p-2">
          <span class="sr-only">Personal Twitter</span>
          <fa :icon="['fab', 'twitter']" fixed-width />
        </a>
        <a href="https://dpb587.me/" class="mx-1 text-gray-600 hover:text-gray-700 p-2">
          <span class="sr-only">Personal Website</span>
          <fa icon="user-alt" fixed-width />
        </a>
      </p>
    </div>
  </div>
</template>

<style>
body {
  @apply bg-gray-100;
}
</style>