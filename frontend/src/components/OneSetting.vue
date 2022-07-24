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
          <v-select v-if="type == `select`" hide-details density="compact"
            variant="underlined" :items="items" v-model="curValue"
            @update:model-value="onselectupdate">
          </v-select>

          <v-text-field v-if="type == `input`" hide-details density="compact"
            v-model="curValue" @update:model-value="onselectupdate">
          </v-text-field>

          <v-text-field v-if="type == `selectfolder`" hide-details density="compact"
            v-model="curValue" append-icon="mdi-folder-search"
            clear-icon="mdi-close-circle" clearable @click:append="onclickselectfolder"
            @click:clear="clearMessage">
          </v-text-field>

          <v-text-field v-if="type == `button`" hide-details density="compact"
            :disabled="true" flat class="w-100" @click:input="onclickselectfolder">
            {{ curValue }}
          </v-text-field>

          <v-row justify="center" v-if="type == `colorpicker`">
            <v-col cols="6">
              <a class="text-h6 d-block">Dark 主题</a>
              <v-color-picker :swatches="darkswatches" show-swatches
                @update:model-value="ondarkmaincolorupdated">
              </v-color-picker>
            </v-col>
            <v-col cols="6">
              <a class="text-h6 d-block">Light 主题</a>
              <v-color-picker :swatches="lightswatches" show-swatches
                @update:model-value="onlightmaincolorupdated">
              </v-color-picker>
            </v-col>
          </v-row>

        </v-card-actions>
      </v-card>
    </v-hover>
  </v-container>
</template>

<script lang="ts" setup>
import { onMounted, reactive, ref } from 'vue'
import { SelectFolder, GetSetting, SetSetting } from '../../wailsjs/go/imes/Api'
import { useBaseStore } from '../stores/index'
const store = useBaseStore()
const props = defineProps<{
  settingkey: string
  title: string
  desc?: string
  type: string
  items?: string[]
}>()
const curValue = ref()

// color refer vuetify official
// https://vuetifyjs.com/en/styles/colors/#material-colors
const darkswatches = reactive([
  ['#455A64', '#263238'], // blue-grey
  ['#5D4037', '#3E2723'], // brown
  ['#388E3C', '#1B5E20'], // green
  ['#00695C', '#004D40'], // teal
])
const lightswatches = reactive([
  ['#B0BEC5', '#ECEFF1'], // blue-grey
  ['#BCAAA4', '#EFEBE9'], // brown
  ['#A5D6A7', '#E8F5E9'], // green
  ['#80CBC4', '#E0F2F1'], // teal
])


onMounted(() => {
  // 从 config file 中加载默认值
  if (props.type == 'select') {
    GetSetting(props.settingkey).then(
      (v: string) => {
        curValue.value = (v != '') ? v : (props.items ? props.items[0] : '')
      }
    )
  } else if (['input', 'selectfolder'].includes(props.type)) {
    GetSetting(props.settingkey).then(
      (v: string) => {
        curValue.value = v
      }
    )
  }
})

// 用户修改后存储到 config file 中
const onselectupdate = (val: string) => {
  console.log(val)
  SetSetting(props.settingkey, val)
}
const onclickselectfolder = () => {
  SelectFolder('').then(
    (v) => {
      curValue.value = (v != '') ? v : curValue.value
      SetSetting(props.settingkey, curValue.value)
    }
  )
}
const ondarkmaincolorupdated = (c: string) => {
  SetSetting(props.settingkey + '-dark', c)
  store.darkmaincolor = c
}
const onlightmaincolorupdated = (c: string) => {
  SetSetting(props.settingkey + '-light', c)
  store.lightmaincolor = c
}
// 其他辅助函数
const clearMessage = () => {
  curValue.value = ''
}
</script>