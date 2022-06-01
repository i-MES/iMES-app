<template>
  <n-grid :cols="5" item-responsive>
    <n-grid-item span="0 400:1"></n-grid-item>
    <n-grid-item span="5 400:3">
      <n-card :hoverable="true" content-style="padding: 0; text-align:center">
        <template #cover>
          <img class="ma-auto" alt="Vue logo" src="@/assets/images/logo.svg"
            :style="`height:${logoheight}`" />
        </template>
        {{ t('all.imes') }}
        <template #action>
          <n-collapse>
            <n-collapse-item>
              <div>intelligent Manufacturing Execution System</div>
            </n-collapse-item>
          </n-collapse>
        </template>
      </n-card>
    </n-grid-item>
  </n-grid>
</template>

<script lang="ts" setup>
import { NCard, NCollapse, NCollapseItem, NGrid, NGridItem } from 'naive-ui'
import { useI18n } from 'vue-i18n'
const { t } = useI18n({ useScope: 'global' })

withDefaults(
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
