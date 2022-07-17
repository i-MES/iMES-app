<template>
  <v-container>
    <v-row>
      <v-col cols="12" sm="8" offset-sm="2">
        <v-card class="mt-5">
          <v-list lines="two" class="ma-1">
            <v-list-subheader>APP 配置</v-list-subheader>
            <v-list-item prepend-avatar=" " title="iMES 工作模式">
              <template v-slot:subtitle>
                <li>影响产品、工序、合同等数据加载方式。</li>
                <li>影响测试业务包的导入方式。</li>
                <li>影响 Log 的存储方式。</li>
              </template>
              <v-radio-group v-model="store.appConfWorkModel">
                <v-radio key="1" label="本地模式" value="1"> </v-radio>
                <v-radio key="2" label="SaaS模式" value="2" disabled> </v-radio>
              </v-radio-group>
            </v-list-item>

            <v-list-item prepend-avatar=" " title="iMES 数据存储模式">
              <v-menu>
                <template v-slot:activator="{ props }">
                  <v-btn color="primary" v-bind="props"> Activator slot </v-btn>
                </template>
                <v-list>
                  <v-list-item v-for="(item, index) in items" :key="index"
                    :value="index">
                    <v-list-item-title>{{ item.title }}</v-list-item-title>
                  </v-list-item>
                </v-list>
              </v-menu>
            </v-list-item>

            <v-list-item prepend-avatar=" " title="Test Config 加载模式">
              <v-btn @click="selectConfigFolder">选择配置文件目录 </v-btn>
            </v-list-item>
          </v-list>
        </v-card>

        <v-card class="mt-5">
          <v-list lines="two">
            <v-list-subheader>UI 配置</v-list-subheader>

            <v-list-item prepend-avatar=" ">
              <template v-slot:title> 默认主题 </template>
              <template v-slot:subtitle> </template>
              <v-radio-group v-model="store.appTheme">
                <v-radio key="1" label="Dark" value="1"> </v-radio>
                <v-radio key="2" label="Light" value="2"> </v-radio>
              </v-radio-group>
            </v-list-item>

            <v-divider></v-divider>

            <v-list-item prepend-avatar=" ">
              <template v-slot:title> 默认语言 </template>
              <v-radio-group v-model="store.appConfDefaultLang">
                <v-radio key="1" label="English" value="1"> </v-radio>
                <v-radio key="2" label="中文" value="2"> </v-radio>
              </v-radio-group>
            </v-list-item>

          </v-list>
        </v-card>

        <v-card class="mt-5">
          <v-list lines="three" active-strategy="multiple">
            <v-list-subheader>用户配置</v-list-subheader>

            <v-list-item value="notifications">
              <template v-slot:default="{ isActive }">
                <v-list-item-avatar start>
                  <v-checkbox :model-value="isActive" hide-details></v-checkbox>
                </v-list-item-avatar>

                <v-list-item-header>
                  <v-list-item-title>是否启用权限管理</v-list-item-title>
                  <v-list-item-subtitle>登录、测试、配置等操作是否使用权限管理。</v-list-item-subtitle>
                </v-list-item-header>
              </template>
            </v-list-item>

            <v-list-item value="sound">
              <template v-slot:default="{ isActive }">
                <v-list-item-avatar start>
                  <v-checkbox :model-value="isActive" hide-details></v-checkbox>
                </v-list-item-avatar>

                <v-list-item-header>
                  <v-list-item-title>是否启用超级用户</v-list-item-title>
                  <v-list-item-subtitle>超级用户可以删除其他用户、执行危险测试项。</v-list-item-subtitle>
                </v-list-item-header>
              </template>
            </v-list-item>
          </v-list>
        </v-card>

        <v-card class="mt-5">
          <v-list lines="three" active-strategy="multiple">
            <v-list-subheader>数据配置</v-list-subheader>
            <v-list-item prepend-avatar=" ">
              <v-row justify="center">
                <v-col cols="4">
                  <v-btn @click="onclickInitConfigData">初始化 Config Data</v-btn>
                </v-col>
                <v-col cols="8">
                  <v-expand-transition>
                    <v-card v-show="expand" width="80%" class="mx-auto secondary"
                      text="Init config data done">
                    </v-card>
                  </v-expand-transition>
                </v-col>
              </v-row>
            </v-list-item>
          </v-list>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { useBaseStore } from '../stores/index'
import { SelectFolder, CreateTargetExample } from '../../wailsjs/go/imes/Api'
const expand = ref(false)
const store = useBaseStore()

const items = [{ title: '模式 1' }, { title: '模式 2' }, { title: '模式 3' }]

const selectConfigFolder = () => {
  SelectFolder('')
}

const onclickInitConfigData = () => {
  CreateTargetExample()
  expand.value = true
}
</script>
