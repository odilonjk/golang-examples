package main

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func createBooking(db *sql.DB, start, end string) (uuid.UUID, error) {
	// iniciando transacao
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("Erro ao iniciar transacao no banco: ", err.Error())
	}
	defer func() {
		if err != nil {
			// em caso de retorno com erro, eh realizado o rollback
			// eh aqui que realiza o rollback para o caso de transacoes
			// tentando executar simultaneamente
			_ = tx.Rollback()
		} else {
			// tudo deu certo entao eh realizado o commit
			err = tx.Commit()
			if err != nil {
				log.Println("Erro ao commit transacao!!! ", err.Error())
			}
		}
	}()

	// definindo nivel de isolamento da transacao
	// serializable = garante que nao havera registros sobrepostos
	// voce pode alterar para outro nivel de isolamento para testar
	_, err = tx.Exec(`set transaction isolation level serializable`)
	if err != nil {
		log.Fatal("Erro ao definir nivel de isolamento da transacao: ", err.Error())
	}

	// valida se ja nao existe uma reserva utilizando a data
	exists := existsOverlappingBooking(tx, start, end)
	if exists {
		return uuid.UUID{}, errors.New("Ja existe uma reserva utilizando estas datas")
	}

	// um sleep para fingir que existe algum processamento a mais ocorrendo
	time.Sleep(time.Millisecond * 200)

	// cria reserva no banco
	id := uuid.New()
	_, err = tx.Exec(`insert into booking(id, start_date, end_date) values ($1, $2, $3)`, id, start, end)
	if err != nil {
		log.Println("Erro ao persistir reserva: ", err.Error())
	}

	return id, err
}

func existsOverlappingBooking(tx *sql.Tx, start, end string) (exists bool) {
	res, err := tx.Query(`
		SELECT EXISTS(
			SELECT 
				1
			FROM 
				booking
			WHERE 
				(start_date < $2
				and end_date >= $2)
			OR 
				(start_date <= $1
				and end_date > $1)
		)`, start, end)
	if err != nil {
		log.Fatal("Erro ao buscar se existe reserva se sobrepondo: ", err.Error())
	}
	defer res.Close()
	res.Next()
	err = res.Scan(&exists)
	if err != nil {
		log.Fatal("Erro ao ler retorno do banco: ", err.Error())
	}
	return
}
