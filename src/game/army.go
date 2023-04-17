package game

import (
	"api/auth"
	"api/database"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"net/http"

	"github.com/gorilla/mux"
)

const TIME_TO_TRAIN int = 5
const PERCENTAGE_LOOTED float64 = 0.99
const MARCH_TIME int = 180

type Train struct {
	CityName   string `database:"city_name" json:"cityName"`
	TroopCount int    `database:"troop_count" json:"troopCount"`
}

type Training struct {
	CityName  string `database:"city_name" json:"cityName"`
	ArmySize  int    `database:"army_size" json:"armySize"`
	StartTime string `database:"start_time" json:"startTime"`
	EndTime   string `database:"end_time" json:"endTime"`
}

type March struct {
	MarchId       string `database:"march_id" json:"marchId"`
	FromCity      string `database:"from_city" json:"fromCity"`
	ToCity        string `database:"to_city" json:"toCity"`
	FromCityName  string `database:"from_city_name" json:"fromCityName"`
	ToCityName    string `database:"to_city_name" json:"toCityName"`
	FromCityOwner string `database:"from_city_owner" json:"fromCityOwner"`
	ToCityOwner   string `database:"to_city_owner" json:"toCityOwner"`
	ArmySize      int    `database:"army_size" json:"armySize"`
	IsReturn      bool   `database:"returning" json:"returning"`
	IsIncoming    bool   `database:"incoming" json:"incoming"`
	IsAttack      bool   `database:"attack" json:"attack"`
	StartTime     string `database:"start_time" json:"startTime"`
	EndTime       string `database:"end_time" json:"endTime"`
}

type Battle struct {
	FromCityName     string  `database:"from_city_name" json:"fromCityName"`
	ToCityName       string  `database:"to_city_name" json:"toCityName"`
	FromCityOwner    string  `database:"from_city_owner" json:"fromCityOwner"`
	ToCityOwner      string  `database:"to_city_owner" json:"toCityOwner"`
	AttackerArmySize int     `database:"attacker_army_size" json:"attackerArmySize"`
	DefenderArmySize int     `database:"defender_army_size" json:"defenderArmySize"`
	AmountLooted     float64 `database:"amount_looted" json:"amountLooted"`
	BattleTime       string  `database:"battle_time" json:"battleTime"`
	AttackVictory    bool    `database:"attack_victory" json:"attackVictory"`
	Incoming         bool    `database:"incoming" json:"incoming"`
}

func HandleArmyRoutes(r *mux.Router) {
	r.HandleFunc("/armies/train", armyTrain).Methods("POST")
	r.HandleFunc("/armies/move", armyMove).Methods("POST")

	r.HandleFunc("/armies/marches", getMarches).Methods("GET")
	r.HandleFunc("/armies/training/global", getGlobalTraining).Methods("GET")
	r.HandleFunc("/armies/training", getTraining).Methods("GET")
	r.HandleFunc("/armies/battles", getBattleLogs).Methods("GET")

	go func() {
		for {
			handleMarches()
		}
	}()
}

func armyTrain(response http.ResponseWriter, request *http.Request) {
	status := false
	defer func() {
		json.NewEncoder(response).Encode(Status{Status: status})
	}()

	if request.Header["Token"] == nil {
		log.Println("Header Token is nil")
		return
	}

	claims, err := auth.ParseJWT(request.Header["Token"][0])

	if err != nil {
		log.Println("error parsing the token")
		return
	}

	var train Train
	err = json.NewDecoder(request.Body).Decode(&train)

	if err != nil {
		log.Println("Error Decoding Body")
		return
	}

	var query string

	if len(train.CityName) > 0 {
		query = fmt.Sprintf(
			`
			SELECT COUNT(*)
			FROM Buildings
			WHERE city_id = (SELECT city_id FROM Cities WHERE city_name = '%s')
				AND building_type = 'Barracks'
			`,
			train.CityName)
	} else {
		query = fmt.Sprintf(
			`
			SELECT COUNT(*)
			FROM Buildings
			WHERE city_id = (SELECT city_id FROM Cities WHERE city_owner='%s' AND town=0)
				AND building_type = 'Barracks'
			`,
			claims["playerId"])
	}

	var barrackCount int

	database.QueryValue(query, &barrackCount)

	if barrackCount == 0 {
		log.Println("You don't have barracks!")
		return
	}

	if len(train.CityName) > 0 {
		query = fmt.Sprintf(
			`
			INSERT INTO Training VALUES(
				(
					SELECT city_id
					FROM Cities
					WHERE city_name='%s'
				),
				%d,
				NOW(),
				TIMESTAMPADD(SECOND, %d	, NOW())
			)
			`,
			train.CityName, train.TroopCount, int(math.Floor(float64(train.TroopCount*TIME_TO_TRAIN)/float64(barrackCount))))
	} else {
		query = fmt.Sprintf(
			`
			INSERT INTO Training VALUES(
				(
					SELECT city_id
					FROM Cities
					WHERE city_owner='%s' AND town=0
				),
				%d,
				NOW(),
				TIMESTAMPADD(SECOND, %d, NOW())
			)
			`,
			claims["playerId"], train.TroopCount, int(math.Floor(float64(train.TroopCount*TIME_TO_TRAIN)/float64(barrackCount))))
	}

	result, err := database.Execute(query)

	if err != nil {
		log.Println("Error querying database")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		log.Println("No rows affected")
		return
	}

	status = true
}

func armyMove(response http.ResponseWriter, request *http.Request) {
	status := false

	defer func() {
		json.NewEncoder(response).Encode(Status{Status: status})
	}()

	if request.Header["Token"] == nil {
		return
	}

	claims, err := auth.ParseJWT(request.Header["Token"][0])

	if err != nil {
		return
	}

	var march March
	err = json.NewDecoder(request.Body).Decode(&march)

	if err != nil {
		return
	}

	result, err := database.Execute(
		fmt.Sprintf(
			`
			UPDATE Cities
			SET army_size=army_size-%d
			WHERE city_name='%s' AND city_owner='%s'
			`,
			march.ArmySize, march.FromCity, claims["playerId"]))

	if err != nil {
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		return
	}

	result, err = database.Execute(
		fmt.Sprintf(
			`
			INSERT INTO Marches
				(march_id, from_city, to_city, army_size, attack, start_time, end_time)
			VALUES(
				uuid(),
				(SELECT city_id FROM Cities WHERE city_name='%s' AND city_owner='%s'),
				(SELECT city_id FROM Cities WHERE city_name='%s'),
				%d,
				(SELECT city_owner FROM Cities WHERE city_name='%s')!=(SELECT city_owner FROM Cities WHERE city_name='%s'),
				NOW(),
				TIMESTAMPADD(SECOND, %d, NOW()))
			`,
			march.FromCity, claims["playerId"], march.ToCity, march.ArmySize, march.FromCity, march.ToCity, MARCH_TIME))

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
			database.Execute(
				fmt.Sprintf(
					`
					DELETE FROM Marches
					WHERE end_time <= NOW() AND march_id IN %s
					`,
					completedMarchesSQL))
		}
		time.Sleep(time.Millisecond * 250)
	}()

	var marches []March

	database.Query(
		`
		SELECT * FROM Marches WHERE end_time <= NOW()
		`,
		&marches)

	for _, march := range marches {

		if !march.IsAttack {
			// runs if the player is moving between two cities that they own
			result, err := database.Execute(
				fmt.Sprintf(
					`
					UPDATE Cities
					SET army_size=army_size+%d
					WHERE city_id='%s'
					`,
					march.ArmySize, march.ToCity))

			if err != nil {
				return
			}

			rowsAffected, err := result.RowsAffected()
			if err != nil || rowsAffected == 0 {
				return
			}

		} else {
			var willConquer bool
			database.QueryValue(
				fmt.Sprintf(
					`SELECT town=1
					FROM Cities
					WHERE city_id='%s'
					`,
					march.ToCity),
				&willConquer)

			type EnemyPlayer struct {
				Balance  float64 `database:"balance"`
				ArmySize int     `database:"army_size"`
			}

			var enemyPlayer []EnemyPlayer
			database.Query(
				fmt.Sprintf(
					`
					SELECT balance, army_size
					FROM Cities JOIN Accounts ON player_id=city_owner
					WHERE city_id='%s'
					`,
					march.ToCity),
				&enemyPlayer)

			change := 0.0

			if willConquer {
				// runs if the player is attacking a town to conquer it
				if enemyPlayer[0].ArmySize < march.ArmySize {
					// if the player has won the battle, change ownership, set new army size to be remainder of attacking army
					result, err := database.Execute(
						fmt.Sprintf(
							`
							UPDATE Cities
							SET
								city_owner=
									(
										SELECT city_owner
										FROM (SELECT * FROM Cities) AS TempCities
										WHERE city_id='%s'
									),
								army_size=%d-army_size
							WHERE city_id='%s'
							`,
							march.FromCity, march.ArmySize, march.ToCity))

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
							`
							UPDATE Cities
							SET army_size=army_size-%d
							WHERE city_id='%s'
							`,
							march.ArmySize, march.ToCity))

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
							`
							UPDATE Cities SET army_size=0
							WHERE city_id='%s'
							`,
							march.ToCity))

					if err != nil {
						return
					}

					// remove gold taken from owner of city
					result, err = database.Execute(
						fmt.Sprintf(
							`
							UPDATE Accounts
							SET balance=balance-%v
							WHERE player_id=
								(
									SELECT city_owner FROM Cities WHERE city_id='%s'
								)
							`,
							change, march.ToCity))

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
							`
							UPDATE Accounts
							SET balance=balance+%v
							WHERE player_id=
								(
									SELECT city_owner FROM Cities WHERE city_id='%s'
								)
							`,
							change, march.FromCity))

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
							`
							INSERT INTO Marches
								(march_id, from_city, to_city, army_size, attack, start_time, end_time)
							VALUES (
								uuid(),
								'%s',
								'%s',
								%d,
								0,
								NOW(),
								TIMESTAMPADD(SECOND, %d, NOW())
							)
							`,
							march.ToCity, march.FromCity, march.ArmySize-enemyPlayer[0].ArmySize, MARCH_TIME))

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
							`
							UPDATE Cities
							SET army_size=%d
							WHERE city_id='%s'
							`,
							enemyPlayer[0].ArmySize-march.ArmySize, march.ToCity))

					if err != nil {
						return
					}

					rowsAffected, err := result.RowsAffected()
					if err != nil || rowsAffected == 0 {
						return
					}
				}

			}
			result, err := database.Execute(
				fmt.Sprintf(
					`
					INSERT INTO Battles
					VALUES(
						uuid(),
						'%s',
						'%s',
						'%s',
						'%d',
						'%d',
						%v,
						%v
					)
					`,
					march.FromCity, march.ToCity, march.EndTime, march.ArmySize, enemyPlayer[0].ArmySize, enemyPlayer[0].ArmySize < march.ArmySize, change))

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

func getMarches(response http.ResponseWriter, request *http.Request) {
	var result []March
	defer func() {
		json.NewEncoder(response).Encode(result)
	}()

	if request.Header["Token"] == nil {
		return
	}

	claims, err := auth.ParseJWT(request.Header["Token"][0])

	if err != nil {
		return
	}

	database.Query(
		fmt.Sprintf(
			`
			SELECT 
			(SELECT city_name FROM Cities WHERE city_id=from_city) AS from_city_name,
			(SELECT username FROM Cities JOIN Accounts ON city_owner=player_id WHERE city_id=from_city) AS from_city_owner, 
			(SELECT city_name FROM Cities WHERE city_id=to_city) AS to_city_name,
			(SELECT username FROM Cities JOIN Accounts ON city_owner=player_id WHERE city_id=to_city) AS to_city_owner,
			(SELECT to_city IN (SELECT city_id FROM Cities WHERE city_owner='%s')) as incoming,
			(SELECT 
				(SELECT city_owner FROM Cities WHERE city_id=from_city)!=
				(SELECT city_owner FROM Cities WHERE city_id=to_city)
				AND attack=0
			) AS returning,
			army_size, start_time, end_time, attack 
			FROM Marches
			WHERE 
			from_city IN (SELECT city_id FROM Cities WHERE city_owner='%s')
			OR
			to_city IN (SELECT city_id FROM Cities WHERE city_owner='%s')
			`,
			claims["playerId"], claims["playerId"], claims["playerId"]),
		&result)
}

func getGlobalTraining(response http.ResponseWriter, request *http.Request) {
	var result []Training
	defer func() {
		json.NewEncoder(response).Encode(result)
	}()

	if request.Header["Token"] == nil {
		return
	}

	claims, err := auth.ParseJWT(request.Header["Token"][0])

	if err != nil {
		return
	}

	database.Query(
		fmt.Sprintf(
			`
			SELECT city_name, Training.army_size, start_time, end_time
			FROM Training JOIN Cities ON Training.city_id=Cities.city_id
			WHERE city_owner='%s'
			`,
			claims["playerId"]),
		&result)
}

func getTraining(response http.ResponseWriter, request *http.Request) {
	var result []Training

	defer func() {
		if len(result) == 0 {
			json.NewEncoder(response).Encode(Training{})
		} else {
			json.NewEncoder(response).Encode(result[0])
		}
	}()

	if request.Header["Token"] == nil {
		return
	}

	claims, err := auth.ParseJWT(request.Header["Token"][0])

	if err != nil {
		return
	}

	var query string

	cityName := request.URL.Query()["cityName"]

	if len(cityName) > 0 {
		query = fmt.Sprintf(
			`
			SELECT city_name, Training.army_size, start_time, end_time
			FROM Training JOIN Cities ON Training.city_id=Cities.city_id
			WHERE city_name='%s'
			`,
			cityName[0])
	} else {
		query = fmt.Sprintf(
			`
			SELECT Training.army_size, start_time, end_time
			FROM Training JOIN Cities ON Training.city_id=Cities.city_id
			WHERE city_owner='%s' AND town=0
			`,
			claims["playerId"])
	}

	database.Query(query, &result)
}

func getBattleLogs(response http.ResponseWriter, request *http.Request) {
	var result []Battle
	defer func() {
		json.NewEncoder(response).Encode(result)
	}()

	if request.Header["Token"] == nil {
		return
	}

	claims, err := auth.ParseJWT(request.Header["Token"][0])

	if err != nil {
		return
	}

	database.Query(
		fmt.Sprintf(
			`
			SELECT
			(SELECT city_name FROM Cities WHERE city_id=from_city) AS from_city_name,
			(SELECT username FROM Cities JOIN Accounts ON city_owner=player_id WHERE city_id=from_city) AS from_city_owner,
			(SELECT city_name FROM Cities WHERE city_id=to_city) AS to_city_name,
			(SELECT username FROM Cities JOIN Accounts ON city_owner=player_id WHERE city_id=to_city) AS to_city_owner,
			(SELECT to_city IN (SELECT city_id FROM Cities WHERE city_owner='%s')) as incoming,
			attacker_army_size, defender_army_size, battle_time, amount_looted, attack_victory
			FROM Battles
			WHERE
			from_city IN (SELECT city_id FROM Cities WHERE city_owner='%s')
			OR
			to_city IN (SELECT city_id FROM Cities WHERE city_owner='%s')
			AND battle_time > TIMESTAMPADD(DAY, -14)
			ORDER BY battle_time DESC
			`,
			claims["playerId"], claims["playerId"], claims["playerId"]),
		&result)
}
