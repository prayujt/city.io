package game

import (
	"api/database"
	"encoding/json"
	"fmt"
	"time"

	"net/http"

	"github.com/gorilla/mux"
)

const TIME_TO_TRAIN int = 5
const PERCENTAGE_LOOTED float64 = 0.99

type Train struct {
	SessionId  string `database:"session_id" json:"sessionId"`
	TroopCount int    `database:"troop_count" json:"troopCount"`
}

type March struct {
	MarchId      string `database:"march_id" json:"marchId"`
	SessionId    string `database:"session_id" json:"sessionId"`
	FromCity     string `database:"from_city" json:"fromCity"`
	ToCity       string `database:"to_city" json:"toCity"`
	ArmySize     int    `database:"army_size" json:"armySize"`
	TimeToTarget int    `database:"time_to_target" json:"timeToTarget"`
	IsAttack     bool   `database:"attack" json:"attack"`
}

func HandleArmyRoutes(r *mux.Router) {
	r.HandleFunc("/armies/train", armyTrain).Methods("POST")
	r.HandleFunc("/armies/move", armyMove).Methods("POST")

	go func() {
		for {
			handleMarches()
		}
	}()
}

func armyTrain(response http.ResponseWriter, request *http.Request) {
	var train Train
	err := json.NewDecoder(request.Body).Decode(&train)

	status := false

	defer func() {
		json.NewEncoder(response).Encode(Status{Status: status})
	}()

	if err != nil {
		return
	}

	result, err := database.Execute(
		fmt.Sprintf(
			"INSERT INTO Training VALUES((SELECT city_id FROM Cities JOIN Sessions ON player_id=city_owner WHERE session_id='%s'), %d, %d)", train.SessionId, train.TroopCount, train.TroopCount*TIME_TO_TRAIN))

	if err != nil {
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return
	}

	status = true
}

func armyMove(response http.ResponseWriter, request *http.Request) {
	var march March
	err := json.NewDecoder(request.Body).Decode(&march)

	status := false

	defer func() {
		json.NewEncoder(response).Encode(Status{Status: status})
	}()

	if err != nil {
		return
	}

	result, err := database.Execute(
		fmt.Sprintf(
			"UPDATE Cities SET army_size=army_size-%d WHERE city_name='%s' AND city_owner=(SELECT player_id FROM Sessions WHERE session_id='%s')", march.ArmySize, march.FromCity, march.SessionId))

	if err != nil {
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return
	}

	result, err = database.Execute(
		fmt.Sprintf(
			"INSERT INTO Marches(march_id, from_city, to_city, army_size, time_to_target, attack) VALUES(uuid(), (SELECT city_id FROM Cities JOIN Sessions ON player_id=city_owner WHERE city_name='%s' AND session_id='%s'), (SELECT city_id FROM Cities WHERE city_name='%s'), %d, 30, (SELECT city_owner FROM Cities WHERE city_name='%s')!=(SELECT city_owner FROM Cities WHERE city_name='%s'))", march.FromCity, march.SessionId, march.ToCity, march.ArmySize, march.FromCity, march.ToCity))

	if err != nil {
		return
	}

	rowsAffected, err = result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return
	}

	status = true
}

func handleMarches() {
	defer func() {
		database.Execute("DELETE FROM Marches WHERE time_to_target=0")
		time.Sleep(time.Millisecond * 250)
	}()

	var marches []March

	database.Query("SELECT * FROM Marches WHERE time_to_target=0", &marches)

	for _, march := range marches {

		if !march.IsAttack {
			result, err := database.Execute(
				fmt.Sprintf(
					"UPDATE Cities SET army_size=army_size+%d WHERE city_id='%s'", march.ArmySize, march.ToCity))

			if err != nil {
				return
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil || rowsAffected == 0 {
				return
			}

		} else {
			var willConquer bool
			database.QueryValue(fmt.Sprintf("SELECT town=1 FROM Cities WHERE city_id='%s'", march.ToCity), &willConquer)

			type EnemyPlayer struct {
				Balance  float64 `database:"balance"`
				ArmySize int     `database:"army_size"`
			}

			var enemyPlayer []EnemyPlayer
			database.Query(fmt.Sprintf("SELECT balance, army_size FROM Cities JOIN Accounts ON player_id=city_owner WHERE city_id='%s'", march.ToCity), &enemyPlayer)

			if willConquer {
				if enemyPlayer[0].ArmySize < march.ArmySize {
					result, err := database.Execute(
						fmt.Sprintf(
							"UPDATE Cities SET city_owner=(SELECT city_owner FROM (SELECT * FROM Cities) AS TempCities WHERE city_id='%s'), army_size=%d-army_size WHERE city_id='%s'", march.FromCity, march.ArmySize, march.ToCity))

					if err != nil {
						return
					}

					rowsAffected, err := result.RowsAffected()
					if err != nil || rowsAffected == 0 {
						return
					}
				} else {
					result, err := database.Execute(
						fmt.Sprintf(
							"UPDATE Cities SET army_size=army_size-%d WHERE city_id='%s'", march.ArmySize, march.ToCity))

					if err != nil {
						return
					}

					rowsAffected, err := result.RowsAffected()
					if err != nil || rowsAffected == 0 {
						return
					}
				}

			} else {
				if enemyPlayer[0].ArmySize < march.ArmySize {
					change := enemyPlayer[0].Balance * PERCENTAGE_LOOTED

					result, err := database.Execute(
						fmt.Sprintf(
							"UPDATE Cities SET army_size=0 WHERE city_id='%s'", march.ToCity))

					if err != nil {
						return
					}

					result, err = database.Execute(
						fmt.Sprintf(
							"UPDATE Accounts SET balance=balance-%v WHERE player_id=(SELECT city_owner FROM Cities WHERE city_id='%s')", change, march.ToCity))

					if err != nil {
						return
					}

					rowsAffected, err := result.RowsAffected()
					if err != nil || rowsAffected == 0 {
						return
					}

					result, err = database.Execute(
						fmt.Sprintf(
							"UPDATE Accounts SET balance=balance+%v WHERE player_id=(SELECT city_owner FROM Cities WHERE city_id='%s')", change, march.FromCity))

					if err != nil {
						return
					}

					rowsAffected, err = result.RowsAffected()
					if err != nil || rowsAffected == 0 {
						return
					}
				} else {
					result, err := database.Execute(
						fmt.Sprintf(
							"UPDATE Cities SET army_size=%d WHERE city_id='%s'", enemyPlayer[0].ArmySize-march.ArmySize, march.ToCity))

					if err != nil {
						return
					}

					rowsAffected, err := result.RowsAffected()
					if err != nil || rowsAffected == 0 {
						return
					}
				}
			}
		}

	}
}
