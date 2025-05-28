<script lang="ts" setup>
import {definePageMeta} from "#imports";
import IndexHeaderBar from "~/components/ui/headers/IndexHeaderBar/index.vue";

const sections = [
  {label: 'Hero', to: '#hero'},
  {label: 'How it Works', to: '#how-it-works'},
  {label: 'Features', to: '#features'},
  {label: 'Templates', to: '#show-case'}
]

useHead({
  title: 'Code Forge | Welcome'
})

definePageMeta({
  layout: false
})
</script>

<template>
  <div class="flex flex-col">
    <header>
      <IndexHeaderBar :sections="sections"/>
    </header>

    <main class="flex flex-col overflow-auto">
      <section id="hero" class="flex flex-col items-center justify-center min-h-screen text-center px-4 bg-gray-800">
        <h1 class="text-5xl font-bold mb-4 text-primary">Automate Your Development Stack</h1>
        <p class="text-lg mb-8 text-gray-600 max-w-2xl text-primary/70">From templates to deploy, CodeForge helps
          developers and teams create and launch full-stack applications faster than ever.</p>
        <div class="flex gap-4">
          <button class="px-6 py-2 bg-primary text-white rounded hover:bg-primary-dark transition">Get Started</button>
          <button class="px-6 py-2 border border-primary text-primary rounded hover:bg-primary/10 transition">
            Documentation
          </button>
        </div>
      </section>

      <section id="how-it-works" class="flex flex-col items-center justify-center min-h-screen px-4 py-16">
        <h2 class="text-3xl font-bold mb-8">How It Works</h2>
        <ol class="grid grid-cols-1 md:grid-cols-4 gap-8 text-center max-w-5xl">
          <li>
            <h3 class="font-semibold text-xl mb-2">1. Select a Template</h3>
            <p class="text-gray-600">Choose from ready-to-go templates hosted on GitHub, GitLab, or Gitea.</p>
          </li>
          <li>
            <h3 class="font-semibold text-xl mb-2">2. Configure Stack</h3>
            <p class="text-gray-600">Add languages, services, dependencies and structure your stack easily.</p>
          </li>
          <li>
            <h3 class="font-semibold text-xl mb-2">3. Setup CI/CD</h3>
            <p class="text-gray-600">Define your pipeline strategy for local or cloud execution.</p>
          </li>
          <li>
            <h3 class="font-semibold text-xl mb-2">4. Deploy Anywhere</h3>
            <p class="text-gray-600">Launch your app on Kubernetes, Lambda, EC2 and more with one click.</p>
          </li>
        </ol>
      </section>

      <section id="features" class="flex flex-col items-center justify-center min-h-screen px-4 py-16">
        <h2 class="text-3xl font-bold mb-8">Platform Features</h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-8 max-w-6xl">
          <div class="p-4 rounded shadow bg-white">
            <h3 class="font-semibold text-lg">Multi-language Support</h3>
            <p class="text-sm text-gray-600">Create projects in Go, Python, Node.js, Vue, React and more.</p>
          </div>
          <div class="p-4 rounded shadow bg-white">
            <h3 class="font-semibold text-lg">Cloud Integration</h3>
            <p class="text-sm text-gray-600">Provision AWS, GCP, Azure resources with pre-built modules.</p>
          </div>
          <div class="p-4 rounded shadow bg-white">
            <h3 class="font-semibold text-lg">CI/CD Ready</h3>
            <p class="text-sm text-gray-600">Generate GitHub Actions, GitLab CI or custom pipelines.</p>
          </div>
          <div class="p-4 rounded shadow bg-white">
            <h3 class="font-semibold text-lg">Team Access Control</h3>
            <p class="text-sm text-gray-600">Support for developers, tech leads, DevOps and managers.</p>
          </div>
          <div class="p-4 rounded shadow bg-white">
            <h3 class="font-semibold text-lg">Kubernetes Support</h3>
            <p class="text-sm text-gray-600">Deploy your projects using best practices for K8s.</p>
          </div>
          <div class="p-4 rounded shadow bg-white">
            <h3 class="font-semibold text-lg">Secure Secrets Management</h3>
            <p class="text-sm text-gray-600">Handle environment variables, tokens and keys securely.</p>
          </div>
        </div>
      </section>

      <section id="show-case" class="flex flex-col items-center justify-center min-h-screen px-4 py-16">
        <h2 class="text-3xl font-bold mb-8">Explore Ready-to-Use Templates</h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 gap-8 max-w-6xl">
          <div class="p-4 bg-gray-100 rounded shadow">
            <h3 class="text-lg font-semibold">Go + PostgreSQL + Redis</h3>
            <p class="text-sm text-gray-600">Backend stack optimized for APIs and caching.</p>
          </div>
          <div class="p-4 bg-gray-100 rounded shadow">
            <h3 class="text-lg font-semibold">Vue.js + Node + MongoDB</h3>
            <p class="text-sm text-gray-600">Fullstack template for fast prototyping and dashboards.</p>
          </div>
          <div class="p-4 bg-gray-100 rounded shadow">
            <h3 class="text-lg font-semibold">Python + FastAPI + S3</h3>
            <p class="text-sm text-gray-600">Modern async backend with cloud storage support.</p>
          </div>
        </div>
      </section>
    </main>

    <footer>

    </footer>
  </div>

</template>

<style scoped>

</style>