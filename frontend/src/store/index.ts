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

interface IGlobalState {
  username: string
  userstatus: UserStatus
  appstatus: AppStatus
  testitems: string[]
}
export const useBaseStore = defineStore('imesBaseStore', {
  state: () => {
    return {
      username: 'admin',
      userstatus: UserStatus.login,
      appstatus: AppStatus.init,
      testitems: ['ti1', 'ti2', 'ti3']
    }
  },
  getters: {
    firstTestItem(state) {
      return state.testitems[0]
    },
    lastTestItem(state) {
      return state.testitems[-1]
    }
  },
  actions: {
    addTestItem(ti: string) {
      this.testitems.push(ti)
    }
  }
})