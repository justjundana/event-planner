package event

import (
	"database/sql"
	"log"
	"time"

	_models "github.com/justjundana/event-planner/models"
)

type EventRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (r *EventRepository) CheckEventAvailable(id int) (_models.Event, error) {
	var event _models.Event

	row := r.db.QueryRow(`
	SELECT
		events.id, events.user_id, events.category_id, events.image, events.title, events.description, events.location, events.date, events.quota
	FROM
		events
	RIGHT JOIN 
		participants ON participants.event_id = events.id
	WHERE
		CURRENT_TIMESTAMP < events.date AND events.id = ?
	GROUP BY 
		events.id
	HAVING 
		COUNT(participants.event_id) < events.quota`, id)

	err := row.Scan(&event.ID, &event.UserID, &event.CategoryId, &event.Image, &event.Title, &event.Description, &event.Location, &event.Date, &event.Quota)
	if err != nil {
		return event, err
	}

	return event, nil
}

func (r *EventRepository) GetEvents() ([]_models.Event, error) {
	var events []_models.Event
	// this condition will run when events joinable
	rows, err := r.db.Query(`
		SELECT
			events.id, events.user_id, events.category_id, events.image, events.title, events.description, events.location, events.date, events.quota
		FROM
			events
		LEFT JOIN 
			participants ON participants.event_id = events.id
		WHERE
			CURRENT_TIMESTAMP < events.date
		GROUP BY 
			events.id
		HAVING 
			COUNT(participants.event_id) < events.quota
		ORDER BY 
			events.date ASC`)
	if err != nil {
		log.Fatalf("Error")
	}

	defer rows.Close()

	for rows.Next() {
		var event _models.Event

		err = rows.Scan(&event.ID, &event.UserID, &event.CategoryId, &event.Image, &event.Title, &event.Description, &event.Location, &event.Date, &event.Quota)
		if err != nil {
			log.Fatalf("Error")
		}

		events = append(events, event)
	}
	return events, nil
}

func (r *EventRepository) Pagination(limit, offset *int) ([]_models.Event, error) {
	var events []_models.Event
	// this condition will run on pagination
	rows, err := r.db.Query(`
		SELECT 
			id, user_id, image, title,category_id, description, location, date, quota 
		FROM 
			events 
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?`, limit, offset)
	if err != nil {
		log.Fatalf("Error")
	}

	defer rows.Close()

	for rows.Next() {
		var event _models.Event

		err = rows.Scan(&event.ID, &event.UserID, &event.Image, &event.Title, &event.CategoryId, &event.Description, &event.Location, &event.Date, &event.Quota)
		if err != nil {
			log.Fatalf("Error")
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepository) GetEvent(id int) (_models.Event, error) {
	var event _models.Event

	row := r.db.QueryRow(`
	SELECT
		events.id, events.user_id, events.image, events.title, events.category_id, events.description, events.location, events.date, events.quota,
		users.id, users.avatar, users.name, users.email, users.address, users.occupation, users.phone
	FROM
		events
	JOIN 
		users ON users.id = events.user_id
	WHERE
		events.id = ?`, id)

	err := row.Scan(
		&event.ID, &event.UserID, &event.Image, &event.Title, &event.CategoryId, &event.Description, &event.Location, &event.Date, &event.Quota,
		&event.User.ID, &event.User.Avatar, &event.User.Name, &event.User.Email, &event.User.Address, &event.User.Occupation, &event.User.Phone,
	)

	if err != nil {
		return event, err
	}

	return event, nil
}

func (r *EventRepository) SearchEvents(keyword string) ([]_models.Event, error) {
	var events []_models.Event
	rows, err := r.db.Query(`
	SELECT 
		id, user_id, image, title,category_id, description, location, date, quota 
	FROM
		events 
	WHERE 
		CURRENT_TIMESTAMP < date AND (title LIKE ? OR location LIKE ?)
	ORDER BY
		date DESC
		`, "%"+keyword+"%", "%"+keyword+"%")
	if err != nil {
		log.Fatalf("Error")
	}
	defer rows.Close()

	for rows.Next() {
		var event _models.Event

		err := rows.Scan(&event.ID, &event.UserID, &event.Image, &event.Title, &event.CategoryId, &event.Description, &event.Location, &event.Date, &event.Quota)
		if err != nil {
			log.Fatalf("Error")
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepository) GetEventMostAttendant() ([]_models.Event, error) {
	var events []_models.Event
	rows, err := r.db.Query(`
	SELECT
		events.id, events.user_id, events.category_id, events.image, events.title, events.description, events.location, events.date, events.quota,
		COUNT(participants.event_id) AS NumberOfParticipant
	FROM
		events
	JOIN 
		participants ON participants.event_id = events.id
	WHERE
		CURRENT_TIMESTAMP < events.date
	GROUP BY 
		participants.event_id
	ORDER BY 
		NumberOfParticipant DESC
	LIMIT 4`)
	if err != nil {
		log.Fatalf("Error")
	}

	defer rows.Close()

	for rows.Next() {
		var event _models.Event

		err = rows.Scan(&event.ID, &event.UserID, &event.CategoryId, &event.Image, &event.Title, &event.Description, &event.Location, &event.Date, &event.Quota, &event.UserID)
		if err != nil {
			log.Fatalf("Error")
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepository) GetOwnEvent(userID int) ([]_models.Event, error) {
	var events []_models.Event
	rows, err := r.db.Query(`SELECT id, user_id, image, title,category_id, description, location, date, quota FROM events WHERE user_id = ?`, userID)
	if err != nil {
		log.Fatalf("Error")
	}

	defer rows.Close()

	for rows.Next() {
		var event _models.Event

		err := rows.Scan(&event.ID, &event.UserID, &event.Image, &event.Title, &event.CategoryId, &event.Description, &event.Location, &event.Date, &event.Quota)
		if err != nil {
			log.Fatalf("Error")
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepository) GetParticipateEvent(userID int) ([]_models.Event, error) {
	var events []_models.Event
	rows, err := r.db.Query(`
	SELECT 
		events.id, events.image, events.title, events.category_id, events.description, events.location, events.date, events.quota 
	FROM 
		events 
	JOIN 
		participants ON participants.event_id = events.id
	WHERE 
		participants.status = TRUE AND participants.user_id = ?`, userID)
	if err != nil {
		log.Fatalf("Error")
	}

	defer rows.Close()

	for rows.Next() {
		var event _models.Event

		err := rows.Scan(&event.ID, &event.Image, &event.Title, &event.CategoryId, &event.Description, &event.Location, &event.Date, &event.Quota)
		if err != nil {
			log.Fatalf("Error")
		}

		events = append(events, event)
	}

	return events, nil
}

func (r *EventRepository) CreateEvent(event _models.Event) error {
	_, err := r.db.Exec("INSERT INTO events(user_id, image, title,category_id, description, location, date, quota, created_at, updated_at) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", event.UserID, event.Image, event.Title, event.CategoryId, event.Description, event.Location, event.Date, event.Quota, time.Now(), time.Now())
	return err
}

func (r *EventRepository) UpdateEvent(event _models.Event) error {
	query := `UPDATE events SET image = ?, title = ?, category_id = ?, description = ?, location = ?, date = ?, quota = ?, updated_at = ? WHERE id = ?`

	statement, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(event.Image, event.Title, event.CategoryId, event.Description, event.Location, event.Date, event.Quota, time.Now(), event.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *EventRepository) DeleteEvent(event _models.Event) error {
	query := `DELETE FROM events WHERE id = ?`

	statement, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	defer statement.Close()

	_, err = statement.Exec(event.ID)
	if err != nil {
		return err
	}

	return nil
}
