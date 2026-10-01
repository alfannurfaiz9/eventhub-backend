package repo

import (
	"context"
	"errors"

	"github.com/alfannurfaiz9/eventhub-backend.git/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommunityRepo struct {
	db *pgxpool.Pool
}

func NewCommunityRepo(db *pgxpool.Pool) *CommunityRepo {
	return &CommunityRepo{
		db: db,
	}
}

func (c *CommunityRepo) GetCommunities(ctx context.Context, categories string) ([]model.CommunityList, error) {
	sql := `
	SELECT communities.name, communities.img_url,communities.description, STRING_AGG(categories.name, ', '), COUNT(user_community.community_id), COUNT(events.id)
	FROM communities
	LEFT JOIN community_category ON community_category.community_id = communities.id
	LEFT JOIN categories ON categories.id = community_category.category_id
	LEFT JOIN user_community ON user_community.community_id = communities.id
	LEFT JOIN events ON events.community_id = communities.id
	WHERE categories.name ILIKE $1
	GROUP BY communities.id, categories.id`
	args := []any{"%" + categories + "%"}

	rows, err := c.db.Query(ctx, sql, args...)

	if err != nil {
		return nil, err
	}

	var communities []model.CommunityList

	for rows.Next() {
		var community model.CommunityList

		if err := rows.Scan(
			&community.Community.Name,
			&community.Community.ImgUrl,
			&community.Community.Description,
			&community.Category.Name,
			&community.TotalMember,
			&community.UpcomingEvent,
		); err != nil {
			return nil, err
		}
		communities = append(communities, community)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return communities, nil
}

func (c *CommunityRepo) GetCommunityDetail(ctx context.Context, id int) (model.CommunityList, error) {
	sql := `
	SELECT communities.name, communities.img_url,communities.description, STRING_AGG(categories.name, ', '), COUNT(user_community.community_id), COUNT(events.id)
	FROM communities
	LEFT JOIN community_category ON community_category.community_id = communities.id
	LEFT JOIN categories ON categories.id = community_category.category_id
	LEFT JOIN user_community ON user_community.community_id = communities.id
	LEFT JOIN events ON events.community_id = communities.id
	WHERE categories.id = $1
	GROUP BY communities.id, categories.id`
	args := []any{id}

	var data model.CommunityList
	if err := c.db.QueryRow(ctx, sql, args...).Scan(
		&data.Community.Name,
		&data.Community.ImgUrl,
		&data.Community.Description,
		&data.Category.Name,
		&data.TotalMember,
		&data.UpcomingEvent,
	); err != nil {
		return model.CommunityList{}, err
	}

	return data, nil
}

func (c *CommunityRepo) GetCommunityEvent(ctx context.Context, id int) ([]model.EventList, error) {
	sql := `
	SELECT events.title, events.img_url, STRING_AGG(categories.name, ', '), events.start_at, locations.name, COUNT(user_event.event_id), events.capacity 
	FROM events 
	LEFT JOIN locations ON locations.id = events.location_id 
	LEFT JOIN communities ON communities.id = events.community_id 
	LEFT JOIN event_category ON event_category.event_id = events.id 
	LEFT JOIN categories ON categories.id = event_category.category_id 
	LEFT JOIN user_event ON user_event.event_id = events.id 
	WHERE events.community_id = $1
	GROUP BY events.id, categories.id, locations.id`
	args := []any{id}

	rows, err := c.db.Query(ctx, sql, args...)

	if err != nil {
		return nil, err
	}

	var events []model.EventList

	for rows.Next() {
		var event model.EventList

		if err := rows.Scan(
			&event.Event.Title,
			&event.Event.ImgUrl,
			&event.Category.Name,
			&event.Event.StartAt,
			&event.Location.Name,
			&event.TotalAttendee,
			&event.Event.Capacity,
		); err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return events, nil
}

func (c *CommunityRepo) GetCommunityMember(ctx context.Context, id int) ([]model.CommunityMember, error) {
	sql := `
	SELECT users.full_name
	FROM user_community
	LEFT JOIN users ON users.id = user_community.user_id
	LEFT JOIN communities ON communities.id = user_community.community_id
	WHERE communities.id = $1`
	args := []any{id}

	rows, err := c.db.Query(ctx, sql, args...)

	if err != nil {
		return nil, err
	}

	var members []model.CommunityMember

	for rows.Next() {
		var member model.CommunityMember

		if err := rows.Scan(
			&member.FullName,
		); err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return members, nil
}

func (c *CommunityRepo) GetPopularCommunity(ctx context.Context) ([]model.CommunityList, error) {
	sql := `
	SELECT communities.name, communities.img_url,communities.description, STRING_AGG(categories.name, ', '), COUNT(user_community.community_id), COUNT(events.id)
	FROM communities
	LEFT JOIN community_category ON community_category.community_id = communities.id
	LEFT JOIN categories ON categories.id = community_category.category_id
	LEFT JOIN user_community ON user_community.community_id = communities.id
	LEFT JOIN events ON events.community_id = communities.id
	GROUP BY communities.id
	HAVING  COUNT(user_community.community_id) > 0
	ORDER BY COUNT(user_community.community_id) DESC
	LIMIT 5`

	rows, err := c.db.Query(ctx, sql)

	if err != nil {
		return nil, err
	}

	var communities []model.CommunityList
	for rows.Next() {
		var community model.CommunityList

		if err := rows.Scan(
			&community.Community.Name,
			&community.Community.ImgUrl,
			&community.Community.Description,
			&community.Category.Name,
			&community.TotalMember,
			&community.UpcomingEvent,
		); err != nil {
			return nil, err
		}

		communities = append(communities, community)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return communities, nil
}

func (c *CommunityRepo) JoinCommunity(ctx context.Context, userId int, body model.UserCommunity) error {
	sql := `
	INSERT INTO user_community(user_id, community_id)
	VALUES($1, $2)`
	args := []any{userId, body.CommunityId}

	cmd, err := c.db.Exec(ctx, sql, args...)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("no row affected")
	}

	return nil
}

func (c *CommunityRepo) LeaveCommunity(ctx context.Context, userId int, body model.UserCommunity) error {
	sql := "DELETE FROM user_community WHERE user_id = $1 AND community_id = $2"
	args := []any{userId, body.CommunityId}

	cmd, err := c.db.Exec(ctx, sql, args...)

	if err != nil {
		return err
	}

	if cmd.RowsAffected() == 0 {
		return errors.New("no row affected")

	}

	return nil
}
