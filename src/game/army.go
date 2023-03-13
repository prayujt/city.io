package game

import (
	"api/database"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"net/http"

	"github.com/gorilla/mux"
)

const TIME_TO_TRAIN int = 5
const PERCENTAGE_LOOTED float64 = 0.99
const MARCH_TIME int = 30

type Train struct {
	SessionId  string `database:"session_id" json:"sessionId"`
	TroopCount int    `database:"troop_count" json:"troopCount"`
}

type March struct {
	MarchId   string `database:"march_id" json:"marchId"`
	SessionId string `database:"session_id" json:"sessionId"`
	FromCity  string `database:"from_city" json:"fromCity"`
	ToCity    string `database:"to_city" json:"toCity"`
	ArmySize  int    `database:"army_size" json:"armySize"`
	IsAttack  bool   `database:"attack" json:"attack"`
	StartTime string `database:"start_time" json:"startTime"`
	EndTime   string `database:"end_time" json:"endTime"`
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
			"INSERT INTO Training VALUES((SELECT city_id FROM Cities JOIN Sessions ON player_id=city_owner WHERE session_id='%s'), %d, NOW(), TIMESTAMPADD(SECOND, %d, NOW()))", train.SessionId, train.TroopCount, train.TroopCount*TIME_TO_TRAIN))

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
			"INSERT INTO Marches(march_id, from_city, to_city, army_size, attack, start_time, end_time) VALUES(uuid(), (SELECT city_id FROM Cities JOIN Sessions ON player_id=city_owner WHERE city_name='%s' AND session_id='%s'), (SELECT city_id FROM Cities WHERE city_name='%s'), %d, (SELECT city_owner FROM Cities WHERE city_name='%s')!=(SELECT city_owner FROM Cities WHERE city_name='%s'), NOW(), TIMESTAMPADD(SECOND, %d, NOW()))", march.FromCity, march.SessionId, march.ToCity, march.ArmySize, march.FromCity, march.ToCity, MARCH_TIME))

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
	var completedMarches []string = make([]string, 0)
	defer func() {
		if len(completedMarches) > 0 {
			completedMarchesSQL := "('" + strings.Join(completedMarches, `', '`) + `')`
			database.Execute(fmt.Sprintf("DELETE FROM Marches WHERE end_time <= NOW() AND march_id IN %s", completedMarchesSQL))
		}
		time.Sleep(time.Millisecond * 250)
	}()

	var marches []March

	database.Query("SELECT * FROM Marches WHERE end_time <= NOW()", &marches)

	for _, march := range marches {

		if !march.IsAttack {
			// runs if the player is moving between two cities that they own
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

			change := 0.0

			if willConquer {
				// runs if the player is attacking a town to conquer it
				if enemyPlayer[0].ArmySize < march.ArmySize {
					// if the player has won the battle, change ownership, set new army size to be remainder of attacking army
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
				// runs if the player is attacking a city to loot
				if enemyPlayer[0].ArmySize < march.ArmySize {
					// if the player has won the battle
					change = enemyPlayer[0].Balance * PERCENTAGE_LOOTED

					// remove city garrison army
					result, err := database.Execute(
						fmt.Sprintf(
							"UPDATE Cities SET army_size=0 WHERE city_id='%s'", march.ToCity))

					if err != nil {
						return
					}

					// remove gold taken from owner of city
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

					// give gold stolen to attacking player
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

					// add march back for remaining troops
					result, err = database.Execute(
						fmt.Sprintf(
							"INSERT INTO Marches(march_id, from_city, to_city, army_size, attack, start_time, end_time) VALUES (uuid(), '%s', '%s', %d, 0, NOW(), TIMESTAMPADD(SECOND, %d, NOW()))", march.ToCity, march.FromCity, march.ArmySize-enemyPlayer[0].ArmySize, MARCH_TIME))

					if err != nil {
						return
					}

					rowsAffected, err = result.RowsAffected()
					if err != nil || rowsAffected == 0 {
						return
					}

				} else {
					// set the new garrison size to be the remaining garrison
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
			result, err := database.Execute(fmt.Sprintf(
				"INSERT INTO Battles VALUES(uuid(), '%s', '%s', '%s', '%d', '%d', %v, %v)", march.FromCity, march.ToCity, march.EndTime, march.ArmySize, enemyPlayer[0].ArmySize, enemyPlayer[0].ArmySize < march.ArmySize, change))

			if err != nil {
				return
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil || rowsAffected == 0 {
				return
			}
		}
		completedMarches = append(completedMarches, march.MarchId)
	}
}
