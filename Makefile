include .env

# ---------------- Database Start ---------------------------------------
jet-gen:
	jet -dsn=${DB_DSN} -path=./.jetgen

# ---------------- Database End ---------------------------------------