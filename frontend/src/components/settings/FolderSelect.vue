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
              {{ desc }}<br />
              点击 X 恢复默认值
              <slot name="desc"></slot>
            </div>
          </div>
        </v-card-item>
        <v-card-actions class="pl-5">

          <v-text-field hide-details density="compact" v-model="curValue"
            append-icon="mdi-folder-search" clear-icon="mdi-close-circle" clearable
            @click:append="onclickselectfolder" @click:clear="clearMessage">
          </v-text-field>

        </v-card-actions>
      </v-card>
    </v-hover>
  </v-container>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue'
import { SelectFolder, GetUserCacheDefaultPath, GetStringSetting, SetStringSetting } from '../../../wailsjs/go/imes/Api'
import { useBaseStore } from '../../stores/index'
const store = useBaseStore()
const props = withDefaults(
  defineProps<{
    settingkey: string
    title: string
    desc?: string
    default?: string
    updatestore?: boolean
  }>(),
  {
    updatestore: false
  }
)
const curValue = ref()

onMounted(() => {
  // 从 config file 中加载默认值
  GetStringSetting(props.settingkey).then(
    (v: string) => {
      curValue.value = v
    }
  )
})

// 用户修改后存储到 config file，并刷新到 store.state
const onclickselectfolder = () => {
  SelectFolder('').then(
    (v) => {
      curValue.value = (v != '') ? v : curValue.value
      SetStringSetting(props.settingkey, curValue.value)
      if (props.updatestore) {
        store.LoadStringSetting(props.settingkey)
      }
    }
  )
}
// 其他辅助函数
const clearMessage = () => {
  GetUserCacheDefaultPath().then((v) => {
    if (v) {
      curValue.value = v
    }
  })
}
</script>