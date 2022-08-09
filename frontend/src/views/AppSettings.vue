<template>
  <v-container>
    <v-row>
      <v-col cols="12" sm="8" offset-sm="2">

        <drop-select settingkey="workmode" title="iMES 工作模式" type="select"
          :items="['本地模式', 'SaaS模式']" :updatestore="true">
          <template v-slot:desc>
            <li>影响产品、工序、合同等数据加载方式。</li>
            <li>影响测试业务包的导入方式。</li>
            <li>影响 Log 的存储方式。</li>
          </template>
        </drop-select>

        <folder-select settingkey="usercachepath" title="用户数据存储路径"
          desc="用户 Cache 数据文件存储路径" :default="defaultusercachepath" />

        <folder-select settingkey="pythonvenvpath" title="Python 虚拟环境"
          desc="Python虚拟环境路径，选中 python 文件" />

        <drop-select settingkey="groupparse" title="TestGroup 解析组策略" type="select"
          :items="['组合成一组', '每个文件一组']" :updatestore="true" />

        <color-picker settingkey="maincolor" title="主色调" desc="影响页面主色块、底色等颜色"
          :updatestore="true" />
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import ColorPicker from '../components/settings/ColorPicker.vue'
import DropSelect from '../components/settings/DropSelect.vue'
import FolderSelect from '../components/settings/FolderSelect.vue'
import { GetUserCacheDefaultPath } from '../../wailsjs/go/imes/Api'

const defaultusercachepath = ref('')
GetUserCacheDefaultPath().then((v: string) => { defaultusercachepath.value = v })
</script>
