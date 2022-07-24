<template>
  <v-container>
    <v-row>
      <v-col cols="12" sm="8" offset-sm="2">

        <one-setting settingkey="workmode" title="iMES 工作模式" type="select"
          :items="['本地模式', 'SaaS模式']">
          <template v-slot:desc>
            <li>影响产品、工序、合同等数据加载方式。</li>
            <li>影响测试业务包的导入方式。</li>
            <li>影响 Log 的存储方式。</li>
          </template>
        </one-setting>

        <one-setting settingkey="datacachepath" title="iMES 数据存储路径"
          desc="config、setting 数据存储文件的路径" type="selectfolder" />

        <one-setting settingkey="maincolor" title="主色调" desc="一些主色块的颜色"
          type="colorpicker"></one-setting>
      </v-col>
    </v-row>
  </v-container>
</template>


<script lang="ts" setup>
import OneSetting from '../components/OneSetting.vue'
</script>
