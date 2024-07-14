-- =================
-- Schema 1.0
-- Add base user accounts for fixtures. The only one that can be logged into
-- is the test-service account. Other users are for show only.
-- =================
DELETE FROM kitchens WHERE TRUE;
DELETE FROM accounts WHERE TRUE;

INSERT INTO accounts (account_id, user_id, email, first_name, last_name)
VALUES
  ('acc_2jEwcS7Rla6E5ik5ELa8uoULKOW', 'auth0|665e3646139d9f6300bad5e9', 'test-service@kitchens-app.com', 'Sam', 'Smith'),
  ('acc_2jEx1hZPbnNEoZRmkWqP2BDBB87', 'auth0|000000000000000000000001', 'test-mack@kitchens-app.com', 'Mack', 'Campbell'),
  ('acc_2jEx1fZrWeWKxxAciNcc1ng3fq5', 'auth0|000000000000000000000002', 'test-bill@kitchens-app.com', 'Bill', 'Williams');

INSERT INTO kitchens (kitchen_id, account_id, kitchen_name, bio, handle)
VALUES
  ('ktc_2jEx1e1esA5292rBisRGuJwXc14', 'acc_2jEwcS7Rla6E5ik5ELa8uoULKOW', "Sam's Kitchen", NULL, 'sammycooks'),
  ('ktc_2jEx1eCS13KMS8udlPoK12e5KPW', 'acc_2jEx1hZPbnNEoZRmkWqP2BDBB87', 'The Campbell Kitchen', "The Campbell's ladle out delicious delights with soup-erb flavor", 'Campbell_Soup'),
  ('ktc_2jEx1j3CVPIIAaOwGIORKqHfK89', 'acc_2jEx1fZrWeWKxxAciNcc1ng3fq5', 'Bill in the Kitchen', NULL, 'bbq_bill');