<template>
  <div class="d-flex">
    <v-hover v-slot="{ isHovering, props }" open-delay="200" close-delay="200">
      <v-card :elevation="isHovering ? 12 : 2" class="mx-auto" v-bind="props">
        <img class="m-a-10" alt="Vue logo" src="@/assets/images/logo.svg" :style="`height:${logoheight}`" />
        <v-card-text class="my-4 text-center text-h6">
          {{ t('all.imes') }}
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn :icon="show ? 'mdi-chevron-up' : 'mdi-chevron-down'" @click="show = !show"> </v-btn>
        </v-card-actions>

        <v-expand-transition>
          <div v-show="show">
            <v-divider></v-divider>
            <v-card-text class="pa-10"> intelligent Manufacturing Execution System </v-card-text>
          </div>
        </v-expand-transition>
      </v-card>
    </v-hover>
  </div>
</template>

<script lang="ts" setup>
import { withDefaults, defineProps, ref } from 'vue'
import { useI18n } from 'vue-i18n'
const { t } = useI18n({ useScope: 'global' })

const show = ref(false)
const props = withDefaults(
  defineProps<{
    logoheight: string | number
    msg?: string
  }>(),
  {
    logoheight: '200px',
    msg: 'hello',
  }
)
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.logo {
  display: block;
  /* height: 50%; */
  text-align: center;
}
</style>
