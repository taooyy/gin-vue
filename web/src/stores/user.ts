import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useUserStore = defineStore('user', () => {
  const token = ref(sessionStorage.getItem('token') || '')
  const role = ref(sessionStorage.getItem('role') || '')
  // orgType: 1=平台, 2=学校, 3=供应商...
  const orgType = ref(Number(sessionStorage.getItem('orgType')) || 0)

  function setLoginState(newToken: string, newRole: string, newType: number) {
    token.value = newToken
    role.value = newRole
    orgType.value = newType
    sessionStorage.setItem('token', newToken)
    sessionStorage.setItem('role', newRole)
    sessionStorage.setItem('orgType', String(newType))
  }

  function logout() {
    token.value = ''
    role.value = ''
    orgType.value = 0
    sessionStorage.clear()
  }

  return { token, role, orgType, setLoginState, logout }
})