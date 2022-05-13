<template>
  <v-container>
    <v-expansion-panels multiple>
      <v-expansion-panel v-for="ti in store.testitems" :key="ti.id" :title="ti.title"
        :text="ti.desc">
        <div>
          <v-progress-linear v-model="value" :buffer-value="bufferValue">
          </v-progress-linear>
        </div>
      </v-expansion-panel>
    </v-expansion-panels>
  </v-container>
</template>

<script lang="ts" setup>
import { ref, onMounted, watch, onBeforeUnmount } from "vue";
import { useBaseStore } from "../stores/index";
const store = useBaseStore();

const value = ref(10);
const bufferValue = ref(20);
const interval = setInterval(() => {
  value.value += Math.random() * (15 - 5) + 5;
  bufferValue.value += Math.random() * (15 - 5) + 6;
}, 2000);;

watch(
  value,
  (n) => {
    if (n < 100) return;
    value.value = 0;
    bufferValue.value = 10;
  }
)
onMounted(() => {
  store.loadTestItem()
})

onBeforeUnmount(() => {
  clearInterval(interval);
})
</script>
