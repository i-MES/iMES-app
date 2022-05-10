<template>
  <v-container>
    <v-row justify="space-between">
      <v-btn icon="mdi-heart" color="primary"></v-btn>

      <v-btn icon="mdi-star" color="secondary"></v-btn>

      <v-btn icon="mdi-cached" color="info"></v-btn>

      <v-btn icon="mdi-thumb-up" color="success"></v-btn>
    </v-row>
    <v-row class="text-center">
      <v-col cols="12">
        <v-img :src="logo" class="my-3" contain height="200" />
      </v-col>
    </v-row>
    <v-dialog v-model="dialog">
      <template v-slot:activator="{ props }">
        <v-btn color="primary" v-bind="props"> Open Dialog </v-btn>
      </template>

      <v-card>
        <v-card-text>
          Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do
          eiusmod tempor incididunt ut labore et dolore magna aliqua.
        </v-card-text>
        <v-card-actions>
          <v-btn color="primary" block @click="dialog = false">Close Dialog
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import logo from "../assets/images/logo.svg";
import { useBaseStore } from "../store/index";

const store = useBaseStore();
const dialog = ref(false);
</script>
