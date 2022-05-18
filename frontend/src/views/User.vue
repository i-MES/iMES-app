<template>
  <v-card>User: {{ $route.name }}</v-card>
  <v-card>{{ store.counter }}</v-card>

</template>

<script lang="ts" setup>
import { onMounted, ref } from "vue";
import { useBaseStore, IAppStatusBar } from "../stores/index";
const sheet = ref(true)
interface IAppStatusBar {
  userCount: number
}

onMounted(() => {
  store.appStatusBar.userCount = 3
})
const store = useBaseStore();
</script>
