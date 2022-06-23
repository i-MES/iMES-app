<template>
  <div>
    <app-logo :msg="t('homepage.welcome')" logoheight="250px" />
    <v-container>
      <v-row justify="center">
        <v-col cols="4" xl="2">
          <v-btn block flat variant="outlined" color="primary" height="80px"
            @click="router.push({ name: 'test' })">
            {{ t("nav.test") }}</v-btn>
        </v-col>
        <v-col cols="4" xl="2">
          <v-btn block flat variant="outlined" color="primary" height="80px"
            @click="router.push({ name: 'dashboard' })">{{ t('nav.dashboard') }}</v-btn>
        </v-col>
        <v-col cols="4" xl="2">
          <v-btn block flat variant="outlined" color="primary" height="80px"
            @click="router.push({ name: 'settings' })">{{ t('nav.settings') }}</v-btn>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script lang="ts" setup>
import AppLogo from '@/components/AppLogo.vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
const router = useRouter()
const { t } = useI18n({ useScope: 'global' })
</script>


<style lang="scss">
</style>
