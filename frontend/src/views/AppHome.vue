<template>
  <div>
    <app-logo :msg="t('homepage.welcome')" logoheight="250px" />
    <v-container>
      <v-row></v-row>
      <v-row>
        <v-col cols="4">
          <v-btn block flat variant="outlined" color="primary" height="80px"
            @click="onclickJump('test')">
            测试</v-btn>
        </v-col>
        <v-col cols="4">
          <v-btn block flat variant="outlined" color="primary" height="80px"
            @click="onclickJump('dashboard')">看板</v-btn>
        </v-col>
        <v-col cols="4">
          <v-btn block flat variant="outlined" color="primary" height="80px"
            @click="onclickJump('settings')">配置</v-btn>
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

const onclickJump = (target: string) => {
  router.push({
    name: target,
  })
}
</script>


<style lang="scss">
</style>
