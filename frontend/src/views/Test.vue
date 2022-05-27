<template>
  <v-sheet class="ma-0 pa-0 overflow-y-auto "
    :height="store.availableHeight - store.toolbarheight - store.logHeight">
    <TE v-if="store.TEorTI" />
    <TI v-else />
  </v-sheet>
  <v-sheet class="ma-0 pa-0 overflow-y-auto "
    :height="store.toolbarheight + store.logHeight">
    <TL />
  </v-sheet>
</template>

<script lang="ts" setup>
import { reactive, } from 'vue'
import { useI18n } from 'vue-i18n'
import { useBaseStore } from '../stores/index'
import TI from '../components/TestItem.vue'
import TE from '../components/TestEntity.vue'
import TL from '../components/TestLog.vue'
const { t } = useI18n({ useScope: 'global' })
const store = useBaseStore()
</script>

