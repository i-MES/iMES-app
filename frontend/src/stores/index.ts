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
  counter: number,
  defaultRoute: string, // 默认导航的页面
  appName: string,
  appTheme: string,
  appBarHeight: number,
  appStatusBar: IAppStatusBar,
  testPageViewModel: boolean,
  userStatus: UserStatus,
  appStatus: AppStatus,
  testProductions: imes.TestProduction[],
  teststeps: imes.TestStep[],   // 测试工序
  atciveTestProduction: number, // 当前选中产品
  activeTestStepId: number,     // 当前测试工序（的 id）
  testitems: imes.TestItem[]
  testitemsLogs: imes.TestItemLog[],
  toolbarheight: number,
  tiPageAvilableHeight: number
}

export const useBaseStore = defineStore('imesBaseStore', {
  state: (): TGlobalState => {
    return {
      counter: 0,
      defaultRoute: 'test',
      appName: 'iMES',
      appTheme: 'dark',
      appBarHeight: 30,
      appStatusBar: {},
      testPageViewModel: false,
      testProductions: [],
      userStatus: UserStatus.login,
      appStatus: AppStatus.init,
      teststeps: [],
      atciveTestProduction: 0,
      activeTestStepId: 0,
      testitems: [],
      testitemsLogs: [],
      toolbarheight: 38,
      tiPageAvilableHeight: 0
    }
  },
  getters: {
    testProductionById: (state) => {
      return (id: number): imes.TestProduction | undefined => {
        return state.testProductions.find((tp) => tp.id == id)
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
    addCounter() {
      api.AddCounter().then(
        (ctr) => {
          console.log(ctr)
          this.counter = ctr
        })
    },
    addTestItem(id: string) {
      this.loadTestItem(id).then(
        (ti) => {
          this.testitems.push()
        }
      )
    },
    async loadSteps() {
      const _ids: number[] = []
      this.teststeps.forEach((ts) => {
        _ids.push(ts.id)
      })
      api.LoadTestSteps().then(
        (tss) => {
          console.log(tss)
          tss.forEach((ts) => {
            if (_ids) {
              if (_ids.indexOf(ts.id) < 0) {
                this.teststeps.push(ts)
              }
            } else {
              this.teststeps.push(ts)
            }
          })
        }
      )
    },
    async loadTestItem(path: string = '~/.imes/config.yml') {
      if (path === null) throw new Error("Need TI's file path")
      const ids: number[] = []
      this.testitems.forEach((ti) => {
        ids.push(ti.id)
      })
      api.LoadTestItems(path).then(
        (tis) => {
          tis.forEach((ti) => {
            if (ids) {
              if (ids.indexOf(ti.id) < 0) {
                this.testitems.push(ti)
              }
            } else {
              this.testitems.push(ti)
            }
          })
          console.log('LoadTestItems return length:', tis.length)
          console.log('current testitems length:', this.testitems.length)
        })
    },
    async loadTestProductions() {
      api.OpenConfigFile().then(
        (f) => {
          api.LoadJsonConfigData(f).then((b) => {
            if (b) {
              api.GetJsonProductions().then((_tps) => {
                if (_tps) {
                  this.testProductions = _tps
                  console.log(this.testProductions)
                }
              })
            }
          })
        },
        (err) => {
          console.log(err)
        }
      )
    }
  }
})