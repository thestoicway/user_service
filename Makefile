genmocks:
	@echo "Generating mocks..."
	@mockgen -source=internal/usr/database/database.go -destination=internal/usr/mocks/database.go -package=mocks
	@mockgen -source=internal/usr/service/service.go -destination=internal/usr/mocks/service.go -package=mocks
	@mockgen -source=internal/usr/session_database/session_database.go -destination=internal/usr/mocks/session_database.go -package=mocks
	@mockgen -source=internal/usr/jsonwebtoken/jwt_manager.go -destination=internal/usr/mocks/jwt_manager.go -package=mocks
