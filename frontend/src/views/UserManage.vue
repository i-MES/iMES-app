<template>
  <v-card id="login">
    <v-card-title>Login</v-card-title>
    <v-card-text>
      <v-container>
        <v-form ref="form" v-model="valid" lazy-validation>
          <v-text-field v-model="name" :counter="10" :rules="nameRules" label="Name"
            required></v-text-field>

          <v-text-field v-model="email" :rules="emailRules" label="E-mail">
          </v-text-field>

          <v-text-field v-model="password" :rules="passwordRules" label="PassWord"
            required>
          </v-text-field>

          <v-checkbox v-model="checkbox" label="Remeber user and password?">
          </v-checkbox>

          <v-btn :disabled="!valid" color="primary" class="mr-4" @click="validate">
            Login
          </v-btn>
        </v-form>
      </v-container>
    </v-card-text>
  </v-card>

</template>

<script lang="ts" setup>
import { ref } from 'vue'

const checkbox = ref(false)
const valid = ref(true)
const name = ref('')
const nameRules = [
  (v: string) => !!v || 'Name is required',
  (v: string) => (v && v.length <= 10) || 'Name must be less than 10 characters',
]
const email = ref('')
const emailRules = [
  (v: string) => !!v || 'E-mail is required',
  (v: string) => /.+@.+\..+/.test(v) || 'E-mail must be valid',
]
const password = ref('')
const passwordRules = [
  (v: string) => !!v || 'Password is required',
  (v: string) => /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[$@$!%*?&_])[A-Za-z\d$@$!%*?&_]{12,20}$/.test(v) || 'Password must be valid',
]

const validate = () => {
  console.log('login')
}


</script>

<style scoped>
#login {
  width: 500px;
  margin: auto;
  margin-top: 100px;
  padding: 20px;
}
</style>