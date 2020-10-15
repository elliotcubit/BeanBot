package state

import (
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

var fetchAllServersStatement string = `SELECT * FROM servers`
var createServerStatement string = `INSERT INTO
servers(serverid,
 channelid,
 mostrecentid,
 mostrecentauthorid,
 mostrecentnumber,
 highestnumberachieved
 )
VALUES ('%s', '%s', '%s', '%s', %d, %d)`
var updateServerChannelStatement string = `UPDATE servers
SET
channelid='%s',
mostrecentid='',
mostrecentauthorid=''
WHERE serverid='%s'`
var updateServerStatement string = `UPDATE servers
SET
mostrecentid='%s',
mostrecentauthorid='%s',
mostrecentnumber=%d,
highestnumberachieved=GREATEST(highestnumberachieved, %d)
WHERE channelid='%s'`
var sumAllAmountsInServer string = `SELECT SUM(amount) FROM beans WHERE serverid='%s'`
var createRowStatement string = `INSERT INTO
beans(serverID, userID, amount) VALUES ('%s', '%s', %d)`
var getRowStatement string = `SELECT amount FROM beans
WHERE serverID='%s' AND userID='%s'`
var updateRowStatement string = `UPDATE beans
SET amount=%d WHERE serverID='%s' AND userID='%s'`
var getLeaderboardStatement = `SELECT userID, amount
FROM beans WHERE serverID='%s' ORDER BY amount %s LIMIT %d`
var transferBeansStatement string = `UPDATE beans
SET amount = CASE userID
WHEN '%s' THEN %d
WHEN '%s' THEN %d
ELSE amount END
WHERE serverID='%s' AND userID IN('%s', '%s')`

type BeanData struct {
	User   string
	Amount int
}

type MessageData struct {
	// Necessary to send the message on parsing it
	Session   *discordgo.Session
	GuildID   string
	ChannelID string
	AuthorID  string
	ID        string
	Number    int
	Timestamp time.Time // Required and must be in this format for heap to work
}

type ChannelData struct {
	// The server this channel is registered to
	ServerID string
	// This doesn't really need to be here, but we might as well
	ChannelID string
	// The most recent ID we parsed in the server
	// Needed for bot to recover after a reboot
	// Empty string if we have never seen a message before
	MostRecentID string
	// The author of the last message - prevent duplicates
	MostRecentAuthorID string
	// The most recent number we parsed. -1 if the number isn't set
	MostRecentNumber      int
	HighestNumberAchieved int
}

func GetServerSum(serverID string) int {
	var amount int
	rows, err := database.Query(fmt.Sprintf(
		sumAllAmountsInServer,
		serverID,
	),
	)
	if err != nil {
		return 0
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&amount)
		if err != nil {
			return 0
		}
	}
	return amount
}

func GetAllServers() map[string]*ChannelData {
	result := make(map[string]*ChannelData, 0)

	rows, err := database.Query(fetchAllServersStatement)
	if err != nil {
		panic("Could not load initial server data from database")
	}
	defer rows.Close()
	for rows.Next() {
		var serverid string
		var channelid string
		var mostrecentid string
		var mostrecentauthorid string
		var mostrecentnumber int
		var highestnumberachieved int
		err := rows.Scan(
			&serverid,
			&channelid,
			&mostrecentid,
			&mostrecentauthorid,
			&mostrecentnumber,
			&highestnumberachieved,
		)
		if err != nil {
			panic("Could not load initial server data from database")
		}
		result[channelid] = &ChannelData{
			ServerID:              serverid,
			ChannelID:             channelid,
			MostRecentID:          mostrecentid,
			MostRecentAuthorID:    mostrecentauthorid,
			MostRecentNumber:      mostrecentnumber,
			HighestNumberAchieved: highestnumberachieved,
		}
	}
	return result
}

func AddServerChannelToList(serverID, channelID string) error {
	_, err := database.Exec(
		fmt.Sprintf(
			createServerStatement,
			serverID,
			channelID,
			"",
			"",
			-1,
			-1,
		),
	)
	if err != nil {
		_, err := database.Exec(fmt.Sprintf(
			updateServerChannelStatement,
			channelID,
			serverID,
		),
		)
		return err
	}
	return nil
}

// Updates the most recent message in the server data
func UpdateMessageInServer(m *MessageData) error {
	_, err := database.Exec(
		fmt.Sprintf(
			updateServerStatement,
			m.ID,
			m.AuthorID,
			m.Number,
			m.Number,
			m.ChannelID,
		),
	)
	return err
}

func GetUserBalance(server, user string) (int, error) {
	var amount int
	rows, err := database.Query(fmt.Sprintf(getRowStatement, server, user))
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&amount)
		if err != nil {
			return 0, err
		}
	}
	// If we didn't get a result, amount will be 0 so we're gravy
	return amount, nil
}

// Internal call that doesn't pprint results -- TODO make the pprint a wrapper around this one that convert them
func UglyBeanLeaderboard(serverID string, direction bool, n int) ([]*BeanData, error) {
	var user string
	var amount int
	result := make([]*BeanData, 0)
	dir := "DESC"
	if direction {
		dir = "ASC"
	}
	rows, err := database.Query(fmt.Sprintf(getLeaderboardStatement, serverID, dir, n))
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&user, &amount)
		if err != nil {
			return result, err
		}
		result = append(result, &BeanData{User: user, Amount: amount})
	}
	return result, nil
}

func GetBeanLeaderboard(s *discordgo.Session, serverID string, direction bool, n int) ([]*BeanData, error) {
	var user string
	var amount int
	result := make([]*BeanData, 0)
	dir := "DESC"
	if direction {
		dir = "ASC"
	}
	rows, err := database.Query(fmt.Sprintf(getLeaderboardStatement, serverID, dir, n))
	if err != nil {
		return result, err
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&user, &amount)
		if err != nil {
			return result, err
		}
		// Get user identifier from UUID
		// TODO use nicknames instead if possible
		userStruct, err := s.User(user)
		if err != nil {
			// Show a placeholder if this happens somehow
			log.Println("We couldn't retrieve user data for UUID %s", user)
			result = append(result, &BeanData{User: "[???]", Amount: amount})
			continue
		}
		result = append(result, &BeanData{User: userStruct.String(), Amount: amount})
	}
	return result, nil
}

/*
 Transfer $amount beans from $from to $to in $serverID
*/
func TransferBeans(serverID, from, to string, amount int) error {
	fromBalance, err := GetUserBalance(serverID, from)
	if err != nil {
		return err
	}
	toBalance, err := GetUserBalance(serverID, to)
	if err != nil {
		return err
	}

	if fromBalance < amount {
		return fmt.Errorf("Not enough beans to make this transfer")
	}

	_, err = database.Exec(
		fmt.Sprintf(
			transferBeansStatement,
			from,
			fromBalance-amount,
			to,
			toBalance+amount,
			serverID,
			from,
			to,
		),
	)
	if err != nil {
		return fmt.Errorf("Bean transfer failed.")
	}
	return nil
}

func AddBeans(server, user string, amount int) (int, error) {
	var currentScore int
	var updatedScore int
	rows, err := database.Query(fmt.Sprintf(getRowStatement, server, user))
	if err != nil {
		return 0, err
	}
	defer rows.Close()
	didGetResult := false
	for rows.Next() {
		err := rows.Scan(&currentScore)
		if err != nil {
			return 0, err
		}
		didGetResult = true
	}
	// Create user if there wasn't a row
	if !didGetResult {
		updatedScore = amount
		err := bbCreateUser(server, user, updatedScore)
		if err != nil {
			return 0, err
		}
		// Otherwise, update the row
	} else {
		updatedScore = currentScore + amount
		// Update the row
		_, err := database.Exec(fmt.Sprintf(updateRowStatement, updatedScore, server, user))
		if err != nil {
			return 0, err
		}
	}
	return updatedScore, nil
}

func bbCreateUser(server, user string, amount int) error {
	_, err := database.Exec(fmt.Sprintf(createRowStatement, server, user, amount))
	return err
}
