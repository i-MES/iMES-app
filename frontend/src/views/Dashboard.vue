<template>
  <v-card>
    <v-tabs v-model="activeTab" background-color="red-lighten-2">
      <v-tab v-for="n in length" :key="n" :value="n">
        Item {{ n }}
      </v-tab>
    </v-tabs>
    <v-card-text class="text-center">
      <v-btn :disabled="!length" text @click="length--">
        Remove Tab
      </v-btn>
      <v-divider class="mx-4" vertical></v-divider>
      <v-btn text @click="length++">
        Add Tab
      </v-btn>
    </v-card-text>

    <v-window v-model="activeTab">
      <v-window-item v-for="i in 3" :key="i" :value="i">
        <v-card>
          <TI />
        </v-card>
      </v-window-item>
    </v-window>
  </v-card>

  <TILog />
</template>

<script lang="ts" setup>
import { ref, watch } from "vue"
import { useI18n } from "vue-i18n";
import TI from "../components/TestItem.vue";
import TILog from "../components/TestItemLog.vue";
const { t } = useI18n({ useScope: "global" });

const length = ref(5)
const activeTab = ref(0)

watch(
  length,
  (nl) => {
    activeTab.value = nl - 1
  }
)
</script>
