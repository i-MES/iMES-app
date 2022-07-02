import { defineStore, storeToRefs } from 'pinia'
import { useUserStore } from './user'
import { main, imes, target } from '../../wailsjs/go/models'
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
enum TestStatus {
  ready,
  running,
  passed,
  failed
}
interface ITestStatus {
  testGroupTitle: string,
  testClassId: number,
  testItemId: number,
  status: TestStatus
}
interface ITestStatuses {
  [entityId: number]: ITestStatus[]
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
  testEntities: target.TestEntity[],  // 所有被测实体
  activedTestEntityId: string,      // 选中实体
  testGroups: target.TestGroup[],
  testitemsLogs: target.TestItemLog[],
  addEntity: boolean,
  TEsNotTE: boolean,
  testStatuses: ITestStatuses,
  canSortTestClass: boolean
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
      activedTestEntityId: '',
      testGroups: [],
      testitemsLogs: [],
      addEntity: false,
      TEsNotTE: true,
      testStatuses: {},
      canSortTestClass: false
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
      return (id: number): target.TestGroup | undefined => {
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
    async addTestEntity(te: target.TestEntity) {
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
    async LoadTestGroup() {
      // sync: 加载 & 去重 & 去脏 & 写回
      api.LoadTestGroup(false).then((tgs) => {
        if (tgs) {
          // console.log('load testgroup:', tgs)
          this.testGroups = tgs
        }
      })
    },
    async SaveTestGroup() {
      // sync: 加载 & 去重 & 去脏 & 写回
      api.SaveTestGroup(this.testGroups)
    },
    async newTestGroup(preid: string) {
      this.testGroups.forEach((tg, idx) => {
        if (tg.id == preid) {
          api.UUID().then(
            (_id) => {
              this.testGroups.splice(idx + 1, 0, {
                id: _id,
                title: '',
                desc: '',
                testclasses: []
              })
              api.SaveTestGroup(this.testGroups)
              return
            }
          )
        }
      })
    },
    async delTestGroup(id: string) {
      this.testGroups.forEach((tg, idx) => {
        if (tg.id == id) {
          console.log('-=-=', this.testGroups[idx].testclasses.length)
          if (this.testGroups[idx].testclasses.length == 0) {
            this.testGroups.splice(idx, 1)
            this.appStatusBar.Tips = ''
          } else {
            this.appStatusBar.Tips = '只有空组才允许删除'
            setTimeout(() => {
              delete this.appStatusBar.Tips
            }, 5000)
          }
          api.SaveTestGroup(this.testGroups)
          return
        }
      })
    }
  }
})