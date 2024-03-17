import { ref } from "vue";

export let prihlasen = ref(false);
export const tokenJmeno = "pavouk_token"
export const cislaProcvicJmeno = "pavouk_procvic_"
export const nastaveniJmeno = "pavouk_nastaveni_psani"
export const levelyRychlosti = [80, 130] // spatny |80| ok |130| libovy
export const levelyPresnosti = [90, 97.5] // spatny |90| ok |97.5| libovy
export const maxPismenNaRadek = 629 / 17 // sirka ramecku / sirka pismene