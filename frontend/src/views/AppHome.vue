<template>
  <div>
    <app-logo :msg="t('homepage.welcome')" logoheight="250px" />
  </div>
</template>

<script lang="ts" setup>
import AppLogo from '@/components/AppLogo.vue'
import { useI18n } from 'vue-i18n'
const { t } = useI18n({ useScope: 'global' })
</script>


<style lang="scss">
</style>
