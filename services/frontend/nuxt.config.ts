// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
    compatibilityDate: '2025-05-15',
    devtools: {enabled: true},

    modules: [
        '@nuxt/eslint',
        '@nuxt/fonts',
        '@nuxt/icon',
        '@nuxt/image',
        '@nuxt/test-utils',
        '@nuxt/ui',
        '@pinia/nuxt',
        'pinia-plugin-persistedstate',
        '@formkit/auto-animate',
        '@nuxtjs/tailwindcss'
    ],
    css: ['~/assets/css/main.css'],
})