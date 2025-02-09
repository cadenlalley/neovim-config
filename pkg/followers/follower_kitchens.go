package followers

import "context"

func ListFollowersByKitchenID(ctx context.Context, store Store, kitchenID string) ([]Follower, error) {
	rows, err := store.QueryxContext(ctx, `
		SELECT k.kitchen_id, k.kitchen_name, k.handle, k.avatar
		FROM kitchen_followers f LEFT JOIN kitchens k ON f.kitchen_id = k.kitchen_id
		WHERE f.followed_kitchen_id = ?
		ORDER BY f.created_at;
	`, kitchenID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followers []Follower
	for rows.Next() {
		var follower Follower
		if err := rows.StructScan(&follower); err != nil {
			return nil, err
		}
		followers = append(followers, follower)
	}

	return followers, nil
}

func ListFollowingByKitchenID(ctx context.Context, store Store, kitchenID string) ([]Follower, error) {
	rows, err := store.QueryxContext(ctx, `
		SELECT k.kitchen_id, k.kitchen_name, k.handle, k.avatar
		FROM kitchen_followers f LEFT JOIN kitchens k ON f.followed_kitchen_id = k.kitchen_id
		WHERE f.kitchen_id = ?
		ORDER BY f.created_at;
	`, kitchenID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var followed []Follower
	for rows.Next() {
		var follower Follower
		if err := rows.StructScan(&follower); err != nil {
			return nil, err
		}
		followed = append(followed, follower)
	}

	return followed, nil
}
