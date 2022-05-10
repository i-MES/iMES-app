<template>
  <div>User</div>
</template>

<script lang="ts" setup>
import { useI18n } from "vue-i18n";
import HelloWorld from "@/components/HelloWorld.vue";
import { useBaseStore } from "../store/index";

const { t } = useI18n({ useScope: "global" });
const store = useBaseStore();

console.log(store.username);
</script>
