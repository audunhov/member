<script setup lang="ts">
import { useRoute, useRouter } from 'vue-router'


const router = useRouter()

async function onSubmit(e) {
  const fd = new FormData(e.target)
  const res = await fetch("/api/v1/login", {
    method: "POST",
    body: JSON.stringify(Object.fromEntries(fd.entries()))
  })
  if (res.ok) {
    router.push("/dashboard")
  }
}
async function register() {
  const form = document.getElementById("form")
  const fd = new FormData(form)
  const res = await fetch("/api/v1/register", {
    method: "POST",
    body: JSON.stringify(Object.fromEntries(fd.entries()))
  })
  const data = await res.text()
}

</script>

<template>
  <h1>Logg inn</h1>
  <form @submit.prevent="onSubmit" id="form">
    <label for="email">E-post</label>
    <input name="email" id="email" type="email">

    <label for="password">Passord</label>
    <input name="password" id="password" type="password">

    <button type="submit">Logg inn</button>
    <button @click.prevent="register" type="button">Registrer</button>
  </form>
</template>
