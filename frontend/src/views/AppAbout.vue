<template>
  <v-container>

    <v-row class="my-5" justify="space-around">
      <v-card width="80%">
        <app-logo :msg="t('homepage.welcome')" logoheight="250px" />
      </v-card>
    </v-row>

    <v-row justify="space-around">
      <v-card width="80%">
        <v-card-text>
          <div class="font-weight-bold ml-1 mb-2">
            {{ t('aboutpage.history') }}
          </div>
          <v-timeline direction="vertical" side="end">
            <v-timeline-item fill-dot dot-color="orange">
              <template v-slot:opposite>
                {{ t('aboutpage.h1-date') }}
              </template>
              <div>
                <div class="text-h6">{{ t('aboutpage.h1-title') }}</div>
                <p>{{ t('aboutpage.h1-content') }}</p>
              </div>
            </v-timeline-item>

            <v-timeline-item>
              <template v-slot:opposite>
                {{ t('aboutpage.h2-date') }}
              </template>
              <div>
                <div class="text-h6">{{ t('aboutpage.h2-title') }}</div>
                <p>{{ t('aboutpage.h2-content') }}</p>
              </div>
            </v-timeline-item>

            <v-timeline-item>
              <template v-slot:opposite>
                {{ t('aboutpage.h3-date') }}
              </template>
              <div>
                <div class="text-h6">{{ t('aboutpage.h3-title') }}</div>
                <p>{{ t('aboutpage.h3-content') }}</p>
              </div>
            </v-timeline-item>
          </v-timeline>
        </v-card-text>
      </v-card>
    </v-row>
    <v-row>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import AppLogo from '@/components/AppLogo.vue'
import { useI18n } from 'vue-i18n'
// import { useRouter } from 'vue-router'
// const router = useRouter()
const { t } = useI18n({ useScope: 'global' })

// 已结通过 app.use() 注册，此处不再需要单独 import
// import OpenLink from '../components/plugin/OpenLink.vue'


</script>

<style lang="scss">
</style>
