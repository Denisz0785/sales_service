package schema

import "github.com/jmoiron/sqlx"

const seeds = `
INSERT INTO products (id,name,cost,quantity,date_created,date_updated) VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd390a21','Lego City',3000,56,'2024-05-05T12:12:12Z','2024-05-06T14:15:12Z'),
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11','Lego Chima',2000,50,'2024-05-05T12:12:12Z','2024-05-06T14:15:12Z')
ON CONFLICT DO NOTHING;

INSERT INTO sales (sale_id,product_id,quantity,paid,date_created) VALUES
('b0eebc99-9c0b-4ef8-bb6d-6bb9bd390a41','a0eebc99-9c0b-4ef8-bb6d-6bb9bd390a21',1,3000,'2024-05-05T12:12:12Z'),
('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a51','a0eebc99-9c0b-4ef8-bb6d-6bb9bd390a21',2,6000,'2024-05-05T12:12:12Z'),
('b0eebc99-9c0b-4ef8-bb6d-6bb9bd380a61','a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',1,2000,'2024-05-05T12:12:12Z')
ON CONFLICT DO NOTHING;

INSERT INTO users (user_id,name,email,password_hash,roles,date_created,date_updated) VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd390a03','Admin','admin1@mail.ru','$2y$10$DwO6jVhE6cM7Y1eJg8eX7ePQ6k5wXwzN0y3p4m6o1p4Dx4x3n2h0','{ADMIN,USER}','2024-05-05T12:12:12Z','2024-05-05T12:12:12Z'),
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd390a04','User','user1@mail.ru','$2y$10$DwO6jVhE6cM7Y1eJg8eX7ePQ6k5wXwzN0y3p4m6o1p4Dx4x3n2h0','{USER}','2024-05-05T12:12:12Z','2024-05-05T12:12:12Z')
`

func Seed(db *sqlx.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	if _, err := tx.Exec(seeds); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
	}
	return tx.Commit()
}
