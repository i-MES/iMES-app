<template>
  <v-sheet class="ma-0 pa-0 overflow-y-auto "
    :height="availableHeight - store.toolbarheight - store.logHeight">
    <TI v-if="store.testPageViewModel" />
    <TE v-else />
  </v-sheet>
  <v-sheet class="ma-0 pa-0 overflow-y-auto "
    :height="store.toolbarheight + store.logHeight">
    <TL :availableHeight="availableHeight" />
  </v-sheet>
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
const availableHeight = display.height.value - store.appBarHeight
const windowSize = reactive({
  x: 0,
  y: 0,
})
const onResize = () => {
  windowSize.x = window.innerWidth
  windowSize.y = window.innerHeight
}
</script>

