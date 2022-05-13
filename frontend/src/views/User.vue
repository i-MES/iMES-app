<template>
  <v-card>User: {{ $route.name }}</v-card>
  <v-card>{{ store.counter }}</v-card>
  <v-card>{{ store.anothercounter }}</v-card>
</template>

<script lang="ts" setup>
import { useBaseStore } from "../stores/index";

const store = useBaseStore();
</script>
