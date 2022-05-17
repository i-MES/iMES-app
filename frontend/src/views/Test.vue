<template>
  <v-tabs class="sticky" grow v-model="activeTab" background-color="red-lighten-2">
    <v-tab v-for="n in length" :key="n" :value="n">
      Item {{ n }}
    </v-tab>
  </v-tabs>

  <v-window v-model="activeTab">
    <v-window-item v-for="i in 13" :key="i" :value="i">
      <v-card>
        <TI />
      </v-card>
    </v-window-item>
  </v-window>

  <TILog />
</template>

<script lang="ts" setup>
import { ref, watch } from "vue"
import { useI18n } from "vue-i18n";
import { useBaseStore } from "../stores/index";
import TI from "../components/TestItem.vue";
import TILog from "../components/TestItemLog.vue";
const { t } = useI18n({ useScope: "global" });
const store = useBaseStore();

const length = ref(5)
const activeTab = ref(0)

watch(
  () => store.activeTestStepId,
  (nv) => {
    console.log('active teststep: ', nv)
  }
)
</script>

<style>
.sticky {
  position: stickey;
  top: 10;
  z-index: 1;
}
</style>