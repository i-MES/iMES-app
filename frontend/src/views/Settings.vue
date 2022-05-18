<template>
  <v-row>
    <v-col cols="12" sm="8" offset-sm="2">
      <v-card class="mt-5">
        <v-list lines="two">
          <v-list-subheader>{{ t('settingpage.pagetitle') }}</v-list-subheader>

          <v-list-item prepend-avatar="@/images/logo.png" title="版权声明" value="foobar">
            <template v-slot:subtitle>
              <span class="font-weight-bold">wkevin</span> &mdash; I'll be in your neighborhood doing errands this
              weekend. Do you want to hang out?
            </template>
          </v-list-item>

          <v-divider></v-divider>

          <v-list-item prepend-avatar="@/images/logo.png">
            <template v-slot:title> 默认主题 </template>
            <template v-slot:subtitle> </template>
          </v-list-item>

          <v-divider></v-divider>

          <v-list-item prepend-avatar="@/images/logo.png" title="iMES 工作模式">
            <v-list-item-media>
              <v-menu>
                <template v-slot:activator="{ props }">
                  <v-btn color="primary" v-bind="props"> Activator slot </v-btn>
                </template>
                <v-list>
                  <v-list-item v-for="(item, index) in items" :key="index" :value="index">
                    <v-list-item-title>{{ item.title }}</v-list-item-title>
                  </v-list-item>
                </v-list>
              </v-menu>
            </v-list-item-media>
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
    </v-col>
  </v-row>
  <v-sheet></v-sheet>
</template>

<script lang="ts" setup>
import { useI18n } from 'vue-i18n'
import { useBaseStore } from '../stores/index'
const { t } = useI18n({ useScope: 'global' })
const store = useBaseStore()

const items = [{ title: '工作模式 1' }, { title: '工作模式 2' }, { title: '工作模式 3' }]
</script>
