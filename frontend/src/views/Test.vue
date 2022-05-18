<template>
  <div class="d-flex flex-row">
    <v-tabs class="sticky" direction="vertical" centered v-model="activeTab"
      color="deep-purple-accent-4">
      <v-tab key="all-entity" value="all-entity">All</v-tab>
      <v-tab v-for="n in tabLength" :key="n" :value="n">
        Item {{ n }}
      </v-tab>
      <v-tab key="add-entity" value="add-entity">Add</v-tab>
    </v-tabs>

    <v-window v-model="activeTab">
      <v-window-item key="all-entity" value="all-entity">
        <a>foo</a>
      </v-window-item>
      <v-window-item key="add-entity" value="add-entity">
        <a>bar</a>
      </v-window-item>
      <v-window-item v-for="i in tabLength" :key="i" :value="i">
        <TI :entityId="i" />
      </v-window-item>
    </v-window>
  </div>
  <TILog />
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref, watch } from "vue"
import { useI18n } from "vue-i18n";
import { useBaseStore, IAppStatusBar } from "../stores/index";
import TI from "../components/TestItem.vue";
import TILog from "../components/TestItemLog.vue";
const { t } = useI18n({ useScope: "global" });
const store = useBaseStore();
const tabLength = ref(5)
const activeTab = ref(0)

watch(
  () => store.activeTestStepId,
  (nv) => {
    console.log('active teststep: ', nv)
  }
)

interface IAppStatusBar {
  activeTab: number
}

onMounted(() => {
  store.appStatusBar.activeTab = activeTab.value
})

watch(
  () => activeTab.value,
  (nv) => {
    store.appStatusBar.activeTab = nv
  }
)
</script>

<style>
.sticky {
  position: stickey;
  top: 10;
}
</style>