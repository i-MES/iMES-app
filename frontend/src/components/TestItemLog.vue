<template>
  <v-hover v-slot="{ isHovering, props }" open-delay="600">
    <v-card class="mx-0 v-card--reveal" color="grey-lighten-4" v-bind="props"
      :height="(isHovering || sticky) ? '300px' : '48px'">
      <v-toolbar height="48">
        <v-toolbar-title>TestItem Log</v-toolbar-title>
        <v-spacer></v-spacer>
        <v-btn :icon="sticky ? 'mdi-pin' : 'mdi-pin-off'" @click="sticky = !sticky">
        </v-btn>
      </v-toolbar>

    </v-card>
  </v-hover>
</template>

<script lang="ts" setup>
import { ref } from "vue";
const sticky = ref(false);
</script>

<style>
.v-card--reveal {
  align-items: center;
  bottom: 0;
  justify-content: center;
  opacity: .9;
  position: absolute;
  width: 100%;
}
</style>