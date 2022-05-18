<template>
  <div>
    <Logo :msg="t('homepage.welcome')" logoheight="250px" />
  </div>
</template>

<script lang="ts" setup>
import Logo from '@/components/Logo.vue'
import { useI18n } from 'vue-i18n'
const { t } = useI18n({ useScope: 'global' })
</script>

<style lang="scss"></style>
