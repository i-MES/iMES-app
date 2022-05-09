<template>
  <div class="about">
    <div class="content">
      <div class="comeon">
        <img :src="logo" alt="Logo" />
      </div>
      <ul class="info">
        <li class="info-item">
          <div class="name">{{ t("aboutpage.project-repository") }}</div>
          <OpenLink class="link" href="https://github.com/i-MES/iMES-app">https://github.com/i-MES/iMES-app</OpenLink>
        </li>
      </ul>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { useI18n } from "vue-i18n";
import logo from "@/assets/images/logo.svg";

const { t } = useI18n({ useScope: "global" });
</script>

<style lang="scss">
.about {
  .content {
    position: relative;
    margin: 36px 20px;
    .comeon {
      position: absolute;
      left: 26px;
      top: 38px;
      max-width: 66%;
      img {
        width: 220px;
        height: 180px;
      }
    }
    .info {
      margin: 0 0 0 33%;
      font-size: 24px;
      text-align: left;
      .info-item {
        margin-bottom: 10px;
        .name {
          line-height: 40px;
          font-size: 28px;
          color: #6d6363;
        }
        .link {
          line-height: 30px;
          font-size: 20px;
          color: #5f6c86;
        }
      }
    }
  }

  .thank {
    height: 68px;
    line-height: 68px;
    margin: 36px auto;
    text-align: center;
    font-size: 40px;
  }
}
</style>
