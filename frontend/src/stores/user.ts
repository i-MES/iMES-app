import { defineStore } from 'pinia'

type UserState = {
  firstName: string
  lastName: string
  userId: number | null
}

export const useUserStore = defineStore('auth/user', {
  state: (): UserState => ({
    firstName: '',
    lastName: '',
    userId: null
  }),
  getters: {
    fullName: (state) => `${state.firstName} ${state.lastName}`,
  },
  actions: {

  }
})
