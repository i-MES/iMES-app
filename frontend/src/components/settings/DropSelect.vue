/* eslint-disable @typescript-eslint/ban-types */
<template>
  <v-container class="mx-0 my-1 pa-0">
    <v-hover v-slot="{ isHovering, props }" open-delay="20" close-delay="20">
      <v-card class="ma-0 pa-0" :style="`opacity:` + (isHovering ? 1 : 0.8)"
        v-bind="props" :elevation="isHovering ? 5 : 0">

        <v-card-item>
          <div>
            <div class=" text-h6 mb-1">
              {{ title }}
            </div>
            <div class="text-caption">
              {{ desc }}
              <slot name="desc"></slot>
            </div>
          </div>
        </v-card-item>
        <v-card-actions class="pl-5">

          <!-- v-select == v-input + v-menu -->
          <v-select hide-details density="compact" variant="underlined" :items="items"
            v-model="curValue" @update:model-value="onselectupdate">
          </v-select>

        </v-card-actions>
      </v-card>
    </v-hover>
  </v-container>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { GetStringSetting, SetStringSetting } from '../../../wailsjs/go/imes/Api'
import { useBaseStore } from '../../stores/index'
const store = useBaseStore()
const props = withDefaults(
  defineProps<{
    settingkey: string
    title: string
    desc?: string
    default?: string
    items?: string[]
    updatestore?: boolean
  }>(),
  {
    updatestore: false
  }
)
const curValue = ref()
onMounted(() => {
  // 从 config file 中加载默认值，加载失败 focus 到 props.items[0]
  GetStringSetting(props.settingkey).then(
    (v: string) => {
      curValue.value = (v != '') ? v : (props.items ? props.items[0] : '')
    }
  )
})

// 用户修改后存储到 config file，并刷新到 store.state
const onselectupdate = (val: string) => {
  console.log(val)
  SetStringSetting(props.settingkey, val)
  if (props.updatestore) {
    store.LoadStringSetting(props.settingkey)
  }
}
</script>