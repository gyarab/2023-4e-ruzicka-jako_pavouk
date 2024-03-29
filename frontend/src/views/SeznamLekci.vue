<script setup lang="ts">
import axios from "axios"
import BlokLekce from "../components/BlokLekce.vue";
import { onMounted, onUnmounted, ref } from "vue"
import { Oznacene, checkTeapot, getToken, pridatOznameni, napovedaKNavigaci } from "../utils";
import { useHead } from 'unhead'

useHead({
    title: "Lekce",
    meta: [
        {
            name: "description",
            content: "Jako Pavouk nabízí spoustu lekcí s řadou typů cviční. Od jednotlivých písmen, přes slova až k celým větám. Naučte se psát všemi deseti!",
        },
    ],
    link: [
        {
            rel: "canonical",
            href: "https://jakopavouk.cz/lekce"
        }
    ]
})

const lekce = ref([[]])
const dokoncene = ref([] as number[])
const o = new Oznacene()
const prvniNedokoncena = ref(1)

onMounted(() => {
    const header = getToken() ? { headers: { Authorization: `Bearer ${getToken()}` } } : {}
    axios.get("/lekce", header)
        .then(response => {
            lekce.value = response.data.lekce
            dokoncene.value = response.data.dokoncene
            o.setMax(lekce.value.join(',').split(',').length) // pocet lekci

            for (let i = 0; i < lekce.value.length; i++) {
                for (let j = 0; j < lekce.value[i].length; j++) {
                    if (dokoncene.value.includes(lekce.value[i][j]['id'])) prvniNedokoncena.value += 1
                    else return
                }
            }
        }).catch(e => {
            if (!checkTeapot(e)) {
                pridatOznameni()
                console.log(e)
            }
        })
    document.addEventListener('keydown', e1)
    document.addEventListener('mousemove', e2)
})

onUnmounted(() => {
    document.removeEventListener('keydown', e1)
    document.removeEventListener('mousemove', e2)
})

function e1(e: KeyboardEvent) {
    if (e.key == 'ArrowUp') {
        e.preventDefault()
        o.mensi()
        let lekce: HTMLElement | null = document.querySelector(`[i="true"]`)
        window.scrollTo({ top: lekce?.offsetTop! - 500 })
    } else if (e.key == 'ArrowDown') {
        e.preventDefault()
        o.vetsi()
        let lekce: HTMLElement | null = document.querySelector(`[i="true"]`)
        window.scrollTo({ top: lekce?.offsetTop! - 200 })
    } else if (e.key == 'Enter') {
        e.preventDefault()
        let lekceE: HTMLElement | null = document.querySelector(`[i="true"]`)
        if (lekceE == null) {
            let lekceE: HTMLElement = document.getElementsByClassName(`lekceBlok`)[prvniNedokoncena.value - 1] as HTMLElement
            if (lekceE == null) lekceE = document.getElementsByClassName(`lekceBlok`)[0] as HTMLElement
            lekceE?.click()
        } else lekceE?.click()
    } else if (e.key == 'Tab') {
        e.preventDefault()
        napovedaKNavigaci()
    }
}

function e2() {
    o.index.value = 0
}
</script>

<template>
    <h1>Lekce</h1>
    <div id="seznam">
        <h2>Střední řada</h2>
        <BlokLekce v-if="lekce[0].length == 0" v-for="_ in 4" pismena="..." :jeDokoncena="false" />
        <!-- jen aby tam něco bylo než se to načte -->
        <BlokLekce v-else v-for="l in lekce[0]" :pismena="l['pismena']" :jeDokoncena="dokoncene.includes(l['id'])"
            :oznacena="o.is(l['id'])" :i="o.is(l['id'])" :class="{ nohover: o.index.value != 0 }" />
        <h2>Horní řada</h2>
        <BlokLekce v-if="lekce[0].length == 0" v-for="_ in 5" pismena="..." :jeDokoncena="false" />
        <BlokLekce v-else v-for="l in lekce[1]" :pismena="l['pismena']" :jeDokoncena="dokoncene.includes(l['id'])"
            :oznacena="o.is(l['id'])" :i="o.is(l['id'])" :class="{ nohover: o.index.value != 0 }" />
        <h2>Dolní řada</h2>
        <BlokLekce v-if="lekce[0].length == 0" v-for="_ in 3" pismena="..." :jeDokoncena="false" />
        <BlokLekce v-else v-for="l in lekce[2]" :pismena="l['pismena']" :jeDokoncena="dokoncene.includes(l['id'])"
            :oznacena="o.is(l['id'])" :i="o.is(l['id'])" :class="{ nohover: o.index.value != 0 }" />
        <h2>Diakritika</h2>
        <BlokLekce v-if="lekce[0].length == 0" v-for="_ in 5" pismena="..." :jeDokoncena="false" />
        <BlokLekce v-else v-for="l in lekce[3]" :pismena="l['pismena']" :jeDokoncena="dokoncene.includes(l['id'])"
            :oznacena="o.is(l['id'])" :i="o.is(l['id'])" :class="{ nohover: o.index.value != 0 }" />
        <h2>Poslední soud</h2>
        <BlokLekce v-if="lekce[0].length == 0" v-for="_ in 2" pismena="..." :jeDokoncena="false" />
        <BlokLekce v-else v-for="l in lekce[4]" :pismena="l['pismena']" :jeDokoncena="dokoncene.includes(l['id'])"
            :oznacena="o.is(l['id'])" :i="o.is(l['id'])" :class="{ nohover: o.index.value != 0 }" />
        <h2>Pro programátory</h2>
        <BlokLekce v-if="lekce[0].length == 0" v-for="_ in 2" pismena="..." :jeDokoncena="false" />
        <BlokLekce v-else v-for="l in lekce[5]" :pismena="l['pismena']" :jeDokoncena="dokoncene.includes(l['id'])"
            :oznacena="o.is(l['id'])" :i="o.is(l['id'])" :class="{ nohover: o.index.value != 0 }" />
    </div>
</template>

<style scoped>
#seznam {
    display: flex;
    flex-direction: column;
    gap: 20px;
    text-align: left;
}

h2 {
    margin-top: 10px;
    margin-left: 5px;
}

.nohover {
    pointer-events: none;
    transition-duration: 0.1s;
}

@media screen and (max-width: 1100px) {
    #seznam {
        width: 70vw;
        align-items: center;
    }

    h2 {
        align-self: start;
    }
}
</style>