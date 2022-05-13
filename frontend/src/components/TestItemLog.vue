<template>
  <v-container>
    <v-toolbar>
      <v-toolbar-title>TestItem Log</v-toolbar-title>
      <v-spacer></v-spacer>
      <div>
        <v-switch v-model="sticky" label="Sticky Banner" hide-details></v-switch>
      </div>
    </v-toolbar>
    <v-banner single-line :sticky="sticky">
      We can't save your edits while you are in offline mode.

      <template v-slot:actions>
        <v-btn text color="deep-purple accent-4"> Get Online </v-btn>
      </template>
    </v-banner>
  </v-container>
</template>

<script lang="ts" setup>
import { ref } from "vue";
const sticky = ref(true);
</script>
