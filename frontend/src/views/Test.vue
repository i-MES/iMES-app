<template>
  <TI v-if="store.testPageViewModel" />
  <TE v-else />
  <TL />
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useBaseStore } from '../stores/index'
import { useDisplay } from 'vuetify'
import TI from '../components/TestItem.vue'
import TE from '../components/TestEntity.vue'
import TL from '../components/TestLog.vue'
const { t } = useI18n({ useScope: 'global' })
const display = useDisplay()
const store = useBaseStore()

onMounted(() => {
  store.tiPageAvilableHeight = display.height.value - store.toolbarheight * 2 - 38
})
</script>

<style>
.sticky {
  position: stickey;
}
</style>
