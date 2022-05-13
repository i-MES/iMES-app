import { defineStore, storeToRefs } from 'pinia'
import { useUserStore } from './user'
import { imes } from "../../wailsjs/go/models"
import * as imesMid from "../../wailsjs/go/imes/Middleware";

export enum UserStatus {
  login,
  logout
}
export enum AppStatus {
  init,
  ready,
  testing
}

export type TGlobalState = {
  counter: number,
  anothercounter: number,
  appname: string,
  apptheme: string,
  userstatus: UserStatus,
  appstatus: AppStatus,
  testitems: imes.TestItem[]
}

export const useBaseStore = defineStore('imesBaseStore', {
  state: (): TGlobalState => {
    return {
      counter: 0,
      anothercounter: 0,
      appname: 'iMES',
      apptheme: 'dark',
      userstatus: UserStatus.login,
      appstatus: AppStatus.init,
      testitems: [],
    }
  },
  getters: {
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
      imesMid.AddCounter().then(
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
    async loadTestItem(path: string = '~/.imes/config.yml') {
      if (path === null) throw new Error("Need TI's file path")
      const ids: number[] = []
      this.testitems.forEach((ti) => {
        ids.push(ti.id)
      })
      return imesMid.LoadTestitems(path)
        .then((tis) => {
          tis.forEach((ti) => {
            if (ids) {
              if (ids.indexOf(ti.id) < 0) {
                this.testitems.push(ti)
              }
            } else {
              this.testitems.push(ti)
            }
          })
        })
    }
  }
})