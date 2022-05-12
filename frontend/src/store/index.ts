import { InjectionKey, UnwrapRef } from 'vue'
import { defineStore } from 'pinia'

export enum UserStatus {
  login,
  logout
}
export enum AppStatus {
  init,
  ready,
  testing
}
export interface TestItem {
  id: string,
  title: string,
  desc?: string,
  func?: string
}
export type TGlobalState = {
  apptheme: string,
  username: string,
  userstatus: UserStatus,
  appstatus: AppStatus,
  testitems: TestItem[]
}

export const useBaseStore = defineStore('imesBaseStore', {
  state: (): TGlobalState => {
    return {
      apptheme: 'dark',
      username: 'admin',
      userstatus: UserStatus.login,
      appstatus: AppStatus.init,
      testitems: [],
    }
  },
  getters: {
    firstTestItem: (state) => {
      return state.testitems[0]
    },
    lastTestItem: (state) => {
      return state.testitems[-1]
    }
  },
  actions: {
    addTestItem(ti: TestItem) {
      this.testitems.push(ti)
    }
  }
})