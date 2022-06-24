import { defineStore } from 'pinia'
import { useUserStore } from './user'
import { main, imes, testset } from '../../wailsjs/go/models'
import * as api from '../../wailsjs/go/imes/Api'

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
  sysInfo: main.SysInfo,
  defaultRoute: string, // 默认导航的页面
  appTheme: string,   // 颜色主题
  appBarHeight: number,
  appStatusBar: IAppStatusBar, // 状态栏信息显示
  appConfWorkModel: string,     // 工作模式：'1'-本地、'2'-网络
  appConfDefaultLang: string,   // 默认语言
  toolbarheight: number,
  logHeight: number,
  mainWindowHeight: number,
  userStatus: UserStatus,
  appStatus: AppStatus,
  testProductions: imes.TestProduction[], // 所有产品
  activedProductionId: number,            // 选中产品
  testStages: imes.TestStage[],     // 所有工序
  activedTestStageId: number,       // 选中工序
  testStation: imes.TestStation,  // 工位(only one)
  testEntities: imes.TestEntity[],  // 所有被测实体
  activedTestEntityIp: string,      // 选中实体
  testGroups: testset.TestGroup[]
  testitemsLogs: testset.TestItemLog[],
  addEntity: boolean,
  TEorTI: boolean,
}

export const useBaseStore = defineStore('imesBaseStore', {
  state: (): TGlobalState => {
    return {
      sysInfo: { buildtype: '', platform: '', arch: '' },
      defaultRoute: 'test',
      appTheme: 'dark',
      appBarHeight: 30,
      appStatusBar: {},
      appConfWorkModel: '1',
      appConfDefaultLang: '1',
      toolbarheight: 38,
      logHeight: 0,
      mainWindowHeight: 0,
      userStatus: UserStatus.login,
      appStatus: AppStatus.init,
      testProductions: [],
      activedProductionId: 0,
      testStages: [],
      activedTestStageId: 0,
      testStation: { id: 0, title: '', desc: '', enabledTestStageIds: [], activedTestStageIds: [] },
      testEntities: [],
      activedTestEntityIp: '',
      testGroups: [],
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
    testGroupById: (state) => {
      return (id: number): testset.TestGroup | undefined => {
        return state.testGroups.find((tg) => tg.id == id)
      }
    },
    testStageByProductionId: (state) => {
      return (id: number): imes.TestStage[] => {
        const tss: imes.TestStage[] = []
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
  },
  actions: {
    async initConfig() {
      api.InitTestProductions()
      api.InitTestStage()
      api.InitTestStation()
      api.InitTestEntity()
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
    async addTestEntity(te: imes.TestEntity) {
      let _new = true
      this.testEntities.forEach((_te, idx) => {
        if (_te.ip.toString() == te.ip.toString()) {
          this.testEntities[idx] = te
          console.log('update testentity:', te.ip)
          _new = false
        }
      })
      if (_new) {
        this.testEntities.push(te)
        console.log('create testentity:', te.ip)
      }
    },
    async syncTestEntity() {
      // sync: 加载 & 去重 & 去脏 & 写回
      const _ips: string[] = []
      this.testEntities.forEach((te) => {
        _ips.push(te.ip.toString())
      })
      api.LoadTestEntity().then((tes) => {
        if (tes) {
          tes.forEach((te) => {
            if (_ips) {
              if (_ips.indexOf(te.ip.toString()) < 0) {
                this.testEntities.push(te)
              }
            } else {
              this.testEntities.push(te)
            }
          })
        }
      })
    },
    async syncTestSet() {
      // sync: 加载 & 去重 & 去脏 & 写回
      api.LoadPythonTestGroup(false).then((tgs) => {
        if (tgs) {
          console.log(tgs)
          this.testGroups = tgs
        }
      })
    },
  }
})