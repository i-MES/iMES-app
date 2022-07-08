<template>
  <div>
    <v-dialog v-model="addEntity" persistent max-width="600px">
      <template v-slot:activator>
        <v-btn class="ma-0 pa-0" min-width="30" min-height="30" stacked
          @click="addEntity = true">
          <v-icon>mdi-creation</v-icon>
        </v-btn>
      </template>
      <v-card>
        <v-card-title>
          <span class="text-h5">添加 Entity</span>
        </v-card-title>
        <v-card-text>
          <v-container>
            <v-row align="center">
              <v-col cols="12" sm="6" md="6">
                <v-text-field label="产品" readonly hide-details="auto"
                  :placeholder="store.testProductionById(store.activedProductionId)?.title"
                  persistent-placeholder>
                </v-text-field>
              </v-col>
              <v-col cols="12" sm="6" md="6">
                <v-text-field label="工序" readonly hide-details="auto"
                  :placeholder="store.testStageById(store.activedTestStageId)?.title"
                  persistent-placeholder>
                </v-text-field>
              </v-col>

              <!-- <v-col cols="12" sm="6" md="6">
              <v-text-field v-model="ip1" label="IP Address [Start]*" required
                hint="IP 段扫描起始地址" placeholder="127.0.0.1" persistent-placeholder>
              </v-text-field>
            </v-col>
            <v-col cols="12" sm="6" md="6">
              <v-text-field v-model="ip2" label="IP Address End" hint="IP 段扫描终止地址">
              </v-text-field>
            </v-col> -->

              <v-col cols="2">
                <v-text-field v-model="ip1" placeholder="127" persistent-placeholder
                  :rules="iprules" label="IP">
                </v-text-field>
              </v-col>
              <v-col cols="2">
                <v-text-field v-model="ip2" placeholder="0" persistent-placeholder
                  :rules="iprules">
                </v-text-field>
              </v-col>
              <v-col cols="2">
                <v-text-field v-model="ip3" placeholder="0" persistent-placeholder
                  :rules="iprules">
                </v-text-field>
              </v-col>
              <v-col cols="2">
                <v-text-field v-model="ip41" placeholder="1" persistent-placeholder
                  :rules="iprules">
                </v-text-field>
              </v-col>
              <v-col cols="1">
                <a>~</a>
              </v-col>
              <v-col cols="2">
                <v-text-field v-model="ip42">
                </v-text-field>
              </v-col>

              <v-col cols="12">
                <v-text-field v-model="code" label="扫条码" required hint="英文"
                  :rules="rules" placeholder="entity-realme-x40"></v-text-field>
              </v-col>

              <v-col cols="12">
                <v-text-field v-model="tags" label="Tags" hint="支持多个标签"> </v-text-field>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="blue darken-1" text @click="addEntity = false"> Close </v-btn>
          <v-btn color="blue darken-1" text @click="progressDialog = !progressDialog">
            Save </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-dialog v-model="progressDialog" hide-overlay persistent width="300">
      <template> </template>
      <v-card color="primary" dark>
        <v-card-text>
          Please stand by
          <v-progress-linear indeterminate color="white" class="mb-0">
          </v-progress-linear>
        </v-card-text>
      </v-card>
    </v-dialog>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch, reactive } from 'vue'
import { useBaseStore } from '../../stores/index'
import { MsgDialog } from '../../../wailsjs/go/imes/Api'
const store = useBaseStore()
const addEntity = ref(false)
const progressDialog = ref(false)
const ip1 = ref('127')
const ip2 = ref('0')
const ip3 = ref('0')
const ip41 = ref('1')
const ip42 = ref('1')
const code = ref('')
const tags = ref('')

const rules = reactive([
  value => !!value || 'Required.',
  value => (value && value.length >= 12) || 'Min 12 characters',
])

const iprules = reactive([
  // value => ((value as number) < 1 || (value as number) > 255) || '1-255',
  value => (value as number) >= 0 || '>=0',
  value => (value as number) < 255 || '<255'
])

watch(
  () => addEntity.value,
  (nv) => {
    if (nv) {
      document.onkeydown = function (e) {
        let key = window.event.keyCode
        console.log(key)
        if (key == 27) {
          addEntity.value = false
        }
      }
    } else {
      document.onkeydown = null
    }
  }
)
watch(
  () => progressDialog.value,
  (nv) => {
    if (nv) {
      setTimeout(() => {
        progressDialog.value = false
      }, 500)

      var _ip1 = Number(ip1.value)
      var _ip2 = Number(ip2.value)
      var _ip3 = Number(ip3.value)
      var _ip41 = Number(ip41.value)
      var _ip42 = Number(ip42.value)
      console.log(_ip1, _ip2, _ip3, _ip41, _ip42)
      if (ip1.value == '') {
        MsgDialog('IP address 非法')
        addEntity.value = true
      } else {
        for (let index = _ip41; index <= _ip42; index++) {
          store.addTestEntity({
            ip: [_ip1, _ip2, _ip3, index],
            code: code.value,
            tags: [tags.value],
          })
          addEntity.value = false
        }
      }
    }
  }
)
</script>
