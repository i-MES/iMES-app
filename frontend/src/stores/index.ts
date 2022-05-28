import { defineStore, storeToRefs } from 'pinia'
import { useUserStore } from './user'
import { imes } from "../../wailsjs/go/models"
import * as api from "../../wailsjs/go/imes/Api";

export enum UserStatus {
  login,
  logout
}
export enum AppStatus {
  init,
  ready,
  testing
}
export interface IAppStatusBar {
  [key: string]: number | string
}

export type TGlobalState = {
  defaultRoute: string, // 默认导航的页面
  appTheme: string,   // 颜色主题
  appBarHeight: number,
  appStatusBar: IAppStatusBar, // 状态栏信息显示
  appConfWorkModel: string,     // 工作模式：'1'-本地、'2'-网络
  appConfDefaultLang: string,   // 默认语言
  toolbarheight: number,
  logHeight: number,
  availableHeight: number,
  userStatus: UserStatus,
  appStatus: AppStatus,
  testProductions: imes.TestProduction[], // 所有产品
  activedProductionId: number,            // 选中产品
  testStages: imes.TestStage[],     // 所有工序
  activedTestStageId: number,       // 选中工序
  testStation: imes.TestStation,  // 工位(only one)
  testEntities: imes.TestEntity[],  // 所有被测实体
  activedTestEntity: number,        // 选中实体
  testitems: imes.TestItem[]
  testitemsLogs: imes.TestItemLog[],
  addEntity: boolean,
  TEorTI: boolean,
}

export const useBaseStore = defineStore('imesBaseStore', {
  state: (): TGlobalState => {
    return {
      defaultRoute: 'test',
      appTheme: 'dark',
      appBarHeight: 30,
      appStatusBar: {},
      appConfWorkModel: '1',
      appConfDefaultLang: '1',
      toolbarheight: 38,
      logHeight: 0,
      availableHeight: 0,
      userStatus: UserStatus.login,
      appStatus: AppStatus.init,
      testProductions: [],
      activedProductionId: 0,
      testStages: [],
      activedTestStageId: 0,
      testStation: { id: 0, title: '', desc: '', enabledTestStageIds: [], activedTestStageIds: [] },
      testEntities: [],
      activedTestEntity: 0,
      testitems: [],
      testitemsLogs: [],
      addEntity: false,
      TEorTI: true
    }
  },
  getters: {
    testProductionById: (state) => {
      return (id: number): imes.TestProduction | undefined => {
        return state.testProductions.find((tp) => tp.id == id)
      }
    },
    testStageById: (state) => {
      return (id: number): imes.TestStage | undefined => {
        return state.testStages.find((tp) => tp.id == id)
      }
    },
    testStageByProductionId: (state) => {
      return (id: number): imes.TestStage[] => {
        var tss: imes.TestStage[] = []
        state.testStages.forEach((ts, _) => {
          if (ts.pid == id) {
            tss.push(ts)
          }
        })
        return tss
      }
    },
    userInfo: (state) => {
      const user = useUserStore()
      return {
        ...user
      }
    },
    firstTestItem: (state) => {
      return state.testitems[0]
    },
    lastTestItem: (state) => {
      return state.testitems[-1]
    }
  },
  actions: {
    async initConfig() {
      api.InitTestProductions()
      api.InitTestStage()
      api.InitTestStation()
      api.InitTestEntity()
      api.InitTestItems()
    },
    async syncTestProductions() {
      // sync: 加载 & 去重 & 去脏 & 写回
      api.LoadTestProductions().then(
        (_tps) => {
          if (_tps) {
            this.testProductions = _tps
            console.log(this.testProductions)
          }
        }
      )
      // api.SaveTestProductions([])
    },
    async syncTestStages() {
      // sync: 加载 & 去重 & 去脏 & 写回
      const _ids: number[] = []
      this.testStages.forEach((ts) => {
        _ids.push(ts.id)
      })
      api.LoadTestStages().then(
        (tss) => {
          console.log(tss)
          tss.forEach((ts) => {
            if (_ids) {
              if (_ids.indexOf(ts.id) < 0) {
                this.testStages.push(ts)
              }
            } else {
              this.testStages.push(ts)
            }
          })
        }
      )
    },
    async syncTestStation() {
      // sync: 加载 & 去重 & 去脏 & 写回
      api.LoadTestStation().then((ts) => {
        if (ts) {
          this.testStation = ts
        }
      })
    },
    async syncTestEntity() {
      // sync: 加载 & 去重 & 去脏 & 写回
      api.LoadTestEntity().then((tes) => {
        if (tes) {
          tes.forEach((te) => {
            this.testEntities.push(te)
          })
        }
      })
    },
    async syncTestItem() {
      // sync: 加载 & 去重 & 去脏 & 写回
      api.LoadTestItems().then((tis) => {
        tis.forEach((ti) => {
          this.testitems.push(ti)
        })
      })
    },
  }
})