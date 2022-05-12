<template>
  <v-container>
    <v-expansion-panels>
      <v-expansion-panel
        v-for="ti in store.testitems"
        :key="ti.id"
        :title="ti.title"
        :text="ti.desc"
      >
      </v-expansion-panel>
    </v-expansion-panels>
  </v-container>
</template>

<script lang="ts" setup>
import { useI18n } from "vue-i18n";
import { useBaseStore } from "../store/index";

const { t } = useI18n({ useScope: "global" });
const store = useBaseStore();

[
  {
    id: "ti1",
    title: "Memory Test",
    desc: "aaaaaaaaaaaa",
  },
  {
    id: "ti2",
    title: "CPU Test",
    desc: "bbbbbbbbbb",
  },
  {
    id: "ti3",
    title: "Network Test",
    desc: "ccccccccccc",
  },
].map((ti) => {
  store.addTestItem(ti);
});

console.log(store.username);
</script>
