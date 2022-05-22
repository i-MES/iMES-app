<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <v-tabs class="sticky" centered v-model="activeTab" color="deep-purple-accent-4">
          <v-tab v-for="n in tabLength" :key="n" :value="n"> Entity {{ n }} </v-tab>
        </v-tabs>
        <v-window v-model="activeTab">
          <v-window-item v-for="i in tabLength" :key="i" :value="i">
            <TI :entityId="i" />
          </v-window-item>
        </v-window>
      </v-col>
    </v-row>
  </v-container>
  <TILog />
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useBaseStore } from '../stores/index'
import { useDisplay } from 'vuetify'
import TI from '../components/TestItem.vue'
import TIEntityOverall from '../components/TestEntityOverall.vue'
import TILog from '../components/TestItemLog.vue'
const { t } = useI18n({ useScope: 'global' })
const store = useBaseStore()
const tabLength = ref(5)
const activeTab = ref(0)
const display = useDisplay()

watch(
  () => store.activeTestStepId,
  (nv) => {
    console.log('active teststep: ', nv)
  }
)

onMounted(() => {
  store.tiPageAvilableHeight = display.height.value - store.toolbarheight * 2 - 38
})
</script>

<style>
.sticky {
  position: stickey;
}
</style>
