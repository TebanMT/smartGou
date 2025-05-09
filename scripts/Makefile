ROOT_DIR := $(shell powershell -Command "(Get-Item ..).FullName")

build-lambda-function:
	powershell -Command "New-Item -ItemType Directory -Path '$(ROOT_DIR)/bin/$(FUNC_NAME)' -Force" && \
	cd $(ROOT_DIR)/src/modules/$(MODULE_NAME)/interfaces/lambdas/$(FUNC_NAME) && \
	powershell -Command "$$env:GOOS='linux'; $$env:GOARCH='amd64'; $$env:CGO_ENABLED='0'; go build  -tags lambda.norpc -o '$(ROOT_DIR)/bin/$(FUNC_NAME)/bootstrap' handler.go" && \
	C:\Users\Mendi\go\bin\build-lambda-zip.exe -o $(ROOT_DIR)/bin/$(FUNC_NAME)/function.zip $(ROOT_DIR)/bin/$(FUNC_NAME)/bootstrap

build-all-lambdas:
	make build-lambda-function MODULE_NAME=security FUNC_NAME=request_singup_by_phone
	make build-lambda-function MODULE_NAME=security FUNC_NAME=verify_otp_from_phone
	make build-lambda-function MODULE_NAME=security FUNC_NAME=sign_up_by_email
	make build-lambda-function MODULE_NAME=security FUNC_NAME=confirm_otp_by_email
	make build-lambda-function MODULE_NAME=security FUNC_NAME=login_with_email
	make build-lambda-function MODULE_NAME=users FUNC_NAME=compleate_onbording
	make build-lambda-function MODULE_NAME=users FUNC_NAME=get_user_profile
	make build-lambda-function MODULE_NAME=security FUNC_NAME=refresh_token
	make build-lambda-function MODULE_NAME=security FUNC_NAME=logout
	make build-lambda-function MODULE_NAME=security FUNC_NAME=request_recovery_password
	make build-lambda-function MODULE_NAME=security FUNC_NAME=reset_password
	make build-lambda-function MODULE_NAME=reference FUNC_NAME=get_meta_categories
	make build-lambda-function MODULE_NAME=reference FUNC_NAME=get_categories

build-custom-auth:
	powershell -Command "New-Item -ItemType Directory -Path '$(ROOT_DIR)/bin/custom_auth/define_auth_challenge' -Force" && \
	cd $(ROOT_DIR)/src/modules/security/infrastructure/cognito/define_auth_challenge && \
	powershell -Command "$$env:GOOS='linux'; $$env:GOARCH='amd64'; $$env:CGO_ENABLED='0'; go build  -tags lambda.norpc -o '$(ROOT_DIR)/bin/custom_auth/define_auth_challenge/bootstrap' handler.go" && \
	C:\Users\Mendi\go\bin\build-lambda-zip.exe -o $(ROOT_DIR)/bin/custom_auth/define_auth_challenge/function.zip $(ROOT_DIR)/bin/custom_auth/define_auth_challenge/bootstrap

	powershell -Command "New-Item -ItemType Directory -Path '$(ROOT_DIR)/bin/custom_auth/create_auth_challenge' -Force" && \
	cd $(ROOT_DIR)/src/modules/security/infrastructure/cognito/create_auth_challenge && \
	powershell -Command "$$env:GOOS='linux'; $$env:GOARCH='amd64'; $$env:CGO_ENABLED='0'; go build  -tags lambda.norpc -o '$(ROOT_DIR)/bin/custom_auth/create_auth_challenge/bootstrap' handler.go" && \
	C:\Users\Mendi\go\bin\build-lambda-zip.exe -o $(ROOT_DIR)/bin/custom_auth/create_auth_challenge/function.zip $(ROOT_DIR)/bin/custom_auth/create_auth_challenge/bootstrap

	powershell -Command "New-Item -ItemType Directory -Path '$(ROOT_DIR)/bin/custom_auth/verify_auth_challenge_response' -Force" && \
	cd $(ROOT_DIR)/src/modules/security/infrastructure/cognito/verify_auth_challenge_response && \
	powershell -Command "$$env:GOOS='linux'; $$env:GOARCH='amd64'; $$env:CGO_ENABLED='0'; go build  -tags lambda.norpc -o '$(ROOT_DIR)/bin/custom_auth/verify_auth_challenge_response/bootstrap' handler.go" && \
	C:\Users\Mendi\go\bin\build-lambda-zip.exe -o $(ROOT_DIR)/bin/custom_auth/verify_auth_challenge_response/function.zip $(ROOT_DIR)/bin/custom_auth/verify_auth_challenge_response/bootstrap

	powershell -Command "New-Item -ItemType Directory -Path '$(ROOT_DIR)/bin/custom_auth/custom_email_messages' -Force" && \
	cd $(ROOT_DIR)/src/modules/security/infrastructure/cognito/custom_email_messages && \
	powershell -Command "$$env:GOOS='linux'; $$env:GOARCH='amd64'; $$env:CGO_ENABLED='0'; go build  -tags lambda.norpc -o '$(ROOT_DIR)/bin/custom_auth/custom_email_messages/bootstrap' handler.go" && \
	C:\Users\Mendi\go\bin\build-lambda-zip.exe -o $(ROOT_DIR)/bin/custom_auth/custom_email_messages/function.zip $(ROOT_DIR)/bin/custom_auth/custom_email_messages/bootstrap

deploy:
	cd $(ROOT_DIR)/infraestructure/cdk && npm run deploy
