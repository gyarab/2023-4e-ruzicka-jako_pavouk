package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

type lekcos struct {
	Id      uint   `json:"id"`
	Pismena string `json:"pismena"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Nenašel jsem soubor .env v /backend.")
	}

	fmt.Printf("Připojuješ se na %s (.env)\n", os.Getenv("DB_JMENO"))

	connStr := fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=disable", os.Getenv("DB_UZIV"), os.Getenv("DB_HESLO"), os.Getenv("DB_HOST"), os.Getenv("DB_JMENO"))
	DB, err = sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Databaze se pokazila", err)
	}

	for {
		fmt.Print("Slovnik (s) / Pohadky (p): ")
		var input string
		fmt.Scan(&input)
		if input == "s" {
			PushSlovnik()
			break
		} else if input == "p" {
			PushPohadky()
			break
		}
	}

	fmt.Println("\nHotovo!")
}

/*
Tahle funkce projde vsechny slovicka ve csv souboru a nasazi do databaze do tabulky slovnik spolu s id lekce ve které je pomocí dosavandne naucenych slovicek možné ho napsat
*/
func PushSlovnik() {
	fmt.Println("\nJdem na slovník")

	rowsZ, err1 := DB.Query(`SELECT id, pismena FROM lekce WHERE klavesnice = 'qwertz' OR klavesnice = 'oboje' ORDER BY id ASC;`)
	rowsY, err2 := DB.Query(`SELECT id, pismena FROM lekce WHERE klavesnice = 'qwerty' OR klavesnice = 'oboje' ORDER BY id ASC;`)

	if err1 != nil || err2 != nil {
		fmt.Println("Připojení k databázi se nezdařilo...")
		return
	}
	defer rowsZ.Close()
	defer rowsY.Close()

	lekceZ := []lekcos{}
	for rowsZ.Next() {
		var l lekcos = lekcos{}
		rowsZ.Scan(&l.Id, &l.Pismena)
		lekceZ = append(lekceZ, l)
	}

	lekceY := []lekcos{}
	for rowsY.Next() {
		var l lekcos = lekcos{}
		rowsY.Scan(&l.Id, &l.Pismena)
		lekceY = append(lekceY, l)
	}

	fmt.Println("Lekce z DB načteny")

	f, err := os.Open("./slovnik.txt")
	if err != nil {
		log.Println("spatna cesta k slovniku")
		return
	}
	csvReader := csv.NewReader(f)
	records, _ := csvReader.ReadAll()
	f.Close()

	fmt.Println("Slova načteny")

	st := `INSERT INTO slovnik (slovo, lekceQWERTZ_id, lekceQWERTY_id) VALUES `

	var pismenkaZ string
	var pismenkaY string
	var indexZ int
	var indexY int
	for _, v := range records {
		pismenkaZ = ""
		indexZ = -1
		for _, p := range lekceZ {
			if p.Pismena == "zbylá diakritika" {
				pismenkaZ += "óďťň"
			} else if p.Pismena != "velká písmena (Shift)" && p.Pismena != "závorky" && p.Pismena != "operátory" && p.Pismena != "čísla" && p.Pismena != "interpunkce" {
				pismenkaZ += p.Pismena
			}
			if obsahujeJenOKPismena(v[0], pismenkaZ) {
				indexZ = int(p.Id)
			}
			if indexZ != -1 {
				break
			}
		}
		pismenkaY = ""
		indexY = -1
		for _, p := range lekceY {
			if p.Pismena == "zbylá diakritika" {
				pismenkaY += "óďťň"
			} else if p.Pismena != "velká písmena (shift)" && p.Pismena != "závorky" && p.Pismena != "operátory" && p.Pismena != "čísla" && p.Pismena != "interpunkce" {
				pismenkaY += p.Pismena
			}
			if obsahujeJenOKPismena(v[0], pismenkaY) {
				indexY = int(p.Id)
			}
			if indexY != -1 {
				break
			}
		}

		st += fmt.Sprintf(`('%s', %v, %v), `, v[0], indexZ, indexY)
	}
	st = st[:len(st)-2]
	st += ";"

	_, err = DB.Exec(`
		DROP TABLE IF EXISTS slovnik;
		CREATE TABLE
    		IF NOT EXISTS slovnik (
        		id SERIAL PRIMARY KEY,
        		slovo VARCHAR(50),
        		lekceQWERTZ_id INT,
        		lekceQWERTY_id INT
    	);
	`)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("%v slov jde do DB\n", len(records))

	if _, err := DB.Exec(st); err != nil {
		panic(err)
	}
}

func obsahujeJenOKPismena(slovo string, pismena string) bool {
	ok := true
	for _, x := range slovo {
		if !strings.Contains(pismena, string(x)) {
			ok = false
			break
		}
	}
	return ok
}

func PushPohadky() {
	fmt.Println("\nJdem na pohadky")

	f, err := os.Open("./pohadky.txt")
	if err != nil {
		log.Println("spatna cesta k souboru")
		return
	}
	defer f.Close()

	fmt.Println("Pohádky načteny")

	var st string = `INSERT INTO vety (veta, delka) VALUES `

	scanner := bufio.NewScanner(f)
	var pohadky []string
	for scanner.Scan() {
		pohadky = append(pohadky, scanner.Text())
	}

	sort.Slice(pohadky, func(i, j int) bool {
		return utf8.RuneCountInString(pohadky[i]) < utf8.RuneCountInString(pohadky[j])
	})

	var pocet int
	delky := make(map[int]int)
	for _, poh := range pohadky {
		delka := utf8.RuneCountInString(poh)
		_, ok := delky[delka]
		if ok {
			delky[utf8.RuneCountInString(poh)]++
		} else {
			delky[utf8.RuneCountInString(poh)] = 1
		}
		if delka <= 80 {
			pocet++
			st += fmt.Sprintf(`('%s', %v), `, poh, delka)
		}
	}

	st = st[:len(st)-2]
	st += ";"

	_, err = DB.Exec(`
		DROP TABLE IF EXISTS vety;
		CREATE TABLE
    		IF NOT EXISTS vety (
        		id SERIAL PRIMARY KEY,
        		veta TEXT NOT NULL,
				delka INT NOT NULL
    	);
	`)
	if err != nil {
		log.Panic(err)
	}

	fmt.Printf("%v pohádek jde do DB", pocet)

	if _, err := DB.Exec(st); err != nil {
		panic(err)
	}
}
