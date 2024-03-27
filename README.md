<div align="center">
  <a href="https://jakopavouk.cz" target="_blank">
    <img src='https://jakopavouk.cz/OGnahled.png' width='420' style='border-radius: 10px;'>
  </a>
</div>

# [Jako Pavouk](https://jakopavouk.cz) | Psaní všemi deseti 🕷️
### Webová aplikace na výuku psaní všemi deseti. <br> Chceš se naučit psát rychle Jako Pavouk?

Umět psát všemi deseti je krásná dovednost. Jak se ji ale naučit?
1. založ si **účet** na Pavoukovi
2. dokonči **všechny lekce** (aby sis osvojil/a prstoklad a věděl/a kde jaká klávesa je)
3. piš všemi deseti **_všechno a všude_**, i když zatím píšeš jako šnek (🐌 -> 🕷️)
4. **doporuč** stránku známým

## Co jsem použil?
- **Frontend** je napsaný ve [Vue.js](https://vuejs.org/) s [typescriptem](https://www.typescriptlang.org/) + pure CSS
- **Backend** používá programovací jazyk [Go](https://go.dev/) a framework [Fiber](https://gofiber.io/)
- **Databázi** jsem zvolil [PostgreSQL](https://www.postgresql.org/) hlavně kvůli popularitě a výkonu

## Spuštění
Pro spuštění je nutné mít nainstalovaný nějaký JS runtime ([Node.js](https://nodejs.org/), [Bun](https://bun.sh/) ...), jazyk [Go](https://go.dev/) a [PostgreSQL](https://www.postgresql.org/) databázi. Instalace těchto programů je nejlépe popsána v oficiální dokumentaci na konkrétním webu.

### Návod pro `Arch linux`:

```sh
#----- INSTALACE -----#

# Go
sudo pacman -S go
go version  # go version go1.22 nebo vyšší

# Node.js
sudo pacman -S nodejs
node -v  # v21 nebo vyšší

# PostgreSQL
sudo pacman -S postgresql
postgres --version # postgres (PostgreSQL) 16 nebo vyšší


#----- KONFIGURACE -----#

sudo su postgres
initdb -D /var/lib/postgres/data
exit # odhlásit z postgres účtu

sudo systemctl start postgresql
# teď bychom měli mít databázi 'postgres' vlastněnou uživatelem 'postgres' bez hesla
# tu můžeme použít pro naší aplikaci
psql --user=postgres postgres < initial.sql # načteme testovací data

# soubor .env.example přejmenovat na .env (nic se snad nemusí upravovat)

git clone https://github.com/gyarab/2023-4e-ruzicka-jako_pavouk

cd 2023-4e-ruzicka-jako_pavouk/frontend
npm i


#----- SPUŠTĚNÍ -----#
# nejlépe přes tmux spustit zároveň
go run .
npm run dev

# jelikož ověřování emailu asi fungovat nebude, je v databázi testovací uživatel:
# jméno: test
# heslo: testtest 

```