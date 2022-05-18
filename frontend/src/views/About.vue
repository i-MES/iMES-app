<template>
  <v-container>
    <v-row justify="space-around">
      <v-card width="80%">
        <img class="mx-auto" :src="logo" alt="Logo" />
        <ul class="ma-10">
          <li class="info-item">
            <div class="name">{{ t('aboutpage.project-repository') }}</div>
            <OpenLink class="link" href="https://github.com/i-MES/iMES-app">
              https://github.com/i-MES/iMES-app</OpenLink
            >
          </li>
        </ul>
      </v-card>
    </v-row>
    <v-row justify="space-around">
      <v-card width="80%">
        <v-card-text>
          <div class="font-weight-bold ml-1 mb-2">
            {{ t('aboutpage.history') }}
          </div>
          <v-timeline truncate-line="none" direction="vertical" side="before">
            <v-timeline-item fill-dot dot-color="orange">
              <template v-slot:icon>
                <span>^_^</span>
              </template>
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
  </v-container>
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n'
import logo from '@/assets/images/logo.svg'
import OpenLink from '../components/public/OpenLink.vue'

const { t } = useI18n({ useScope: 'global' })
</script>

<style lang="scss"></style>
