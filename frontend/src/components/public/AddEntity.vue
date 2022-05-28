<template>
  <v-dialog v-model="addEntity" persistent max-width="600px">
    <template v-slot:activator="{ isActive, props }">
      <v-card
        :color="store.appTheme == 'dark' ? 'blue-grey-darken-1' : 'blue-grey-lighten-5'"
        @click="addEntity = true">
        <v-card-text>
          <v-row>
            <v-col class="text-h5">添加 Entity {{ isActive }}
            </v-col>
          </v-row>
        </v-card-text>
        <v-card-actions>
          <v-icon class="ma-auto" icon="mdi-overscan" size="85"></v-icon>
        </v-card-actions>
      </v-card>
    </template>
    <v-card>
      <v-card-title>
        <span class="text-h5">添加 Entity</span>
      </v-card-title>
      <v-card-text>
        <v-container>
          <v-row>
            <v-col cols="12" sm="6" md="6">
              <v-text-field label="产品" readonly
                :placeholder="store.testProductionById(store.activedProductionId)?.title"
                persistent-placeholder>
              </v-text-field>
            </v-col>
            <v-col cols="12" sm="6" md="6">
              <v-text-field label="工序" readonly
                :placeholder="store.testStageById(store.activedTestStageId)?.title"
                persistent-placeholder>
              </v-text-field>
            </v-col>

            <v-col cols="12" sm="6" md="6">
              <v-text-field v-model="ip1" label="IP Address [Start]*" required
                hint="IP 段扫描起始地址" placeholder="127.0.0.1" persistent-placeholder>
              </v-text-field>
            </v-col>
            <v-col cols="12" sm="6" md="6">
              <v-text-field v-model="ip2" label="IP Address End" hint="IP 段扫描终止地址">
              </v-text-field>
            </v-col>
            <v-col cols="12">
              <v-text-field v-model="code" label="扫条码" required hint="英文"
                placeholder="entity-realme-x40"></v-text-field>
            </v-col>
            <v-col cols="12">
              <v-text-field v-model="tags" label="Tags" hint="支持多个标签">
              </v-text-field>
            </v-col>
          </v-row>
        </v-container>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="blue darken-1" text @click="addEntity = false">
          Close
        </v-btn>
        <v-btn color="blue darken-1" text @click="progressDialog = !progressDialog">
          Save
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>

  <v-dialog v-model="progressDialog" hide-overlay persistent width="300">
    <template>
    </template>
    <v-card color="primary" dark>
      <v-card-text>
        Please stand by
        <v-progress-linear indeterminate color="white" class="mb-0">
        </v-progress-linear>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue'
import { imes } from '../../../wailsjs/go/models';
import { useBaseStore } from '../../stores/index'
import { MsgDialog } from '../../../wailsjs/go/imes/Api'
const store = useBaseStore()
const addEntity = ref(false)
const progressDialog = ref(false)
const ip1 = ref('')
const ip2 = ref('')
const code = ref('')
const tags = ref('')

watch(
  () => progressDialog.value,
  (nv) => {
    if (nv) {
      setTimeout(() => {
        progressDialog.value = false
      }, 500)
      console.log("ip1:", ip1)
      console.log("ip2:", ip2)
      if (ip1.value == "") {
        MsgDialog("IP address 非法")
        addEntity.value = true
      } else {
        var _cf = false
        store.testEntities.forEach((te) => {
          if (te.ip == ip1.value) {
            MsgDialog("IP address 重复")
            _cf = true
          }
        })
        if (!_cf) {
          var te: imes.TestEntity = {
            id: store.testEntities.length + 1,
            ip: ip1,
            code: code.value,
            tags: tags.value,
          }
          console.log(te)
          store.testEntities.push(te)
          addEntity.value = false
        }
      }
    }
  }
)
</script>